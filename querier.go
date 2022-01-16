package main

import (
	"fmt"
	"strings"
)

type Querier interface {
	Intersection(tokens ...string) ([]Posting, error)
	Phrase(phrase string) ([]Posting, error)
}

type querier struct {
	idx Index

	intersectionFn func(a, b []Posting) []Posting
}

func NewQuerier(idx Index) Querier {
	return &querier{
		idx: idx,

		// Use the default functions
		intersectionFn: intersection,
	}
}

// Intersection fetches the postings lists for all given terms and returns the
// document ID's present in all lists.
func (q *querier) Intersection(tokens ...string) ([]Posting, error) {
	if len(tokens) == 0 {
		return nil, fmt.Errorf("no tokens provided")
	}

	postingsLists := make([][]Posting, 0, len(tokens))
	var lowestDocFreqIdx int

	// Fetch all postings lists.
	for _, t := range tokens {
		postings, err := q.idx.Postings(t)
		if err != nil {
			return nil, fmt.Errorf("get postings list: %w", err)
		}
		postingsLists = append(postingsLists, postings)

		// Keep track of the shortest postings list. We will start with
		// that one when doing the intersection.
		if len(postingsLists[lowestDocFreqIdx]) > len(postings) {
			lowestDocFreqIdx = len(postingsLists) - 1
		}
	}

	res := postingsLists[lowestDocFreqIdx]
	for i, l := range postingsLists {
		if i == lowestDocFreqIdx {
			continue
		}

		res = q.intersectionFn(res, l)
	}
	return res, nil
}

// intersection returns the common document ID's from the two given postings lists.
func intersection(a, b []Posting) []Posting {
	if len(a) == 0 {
		return nil
	}
	if len(b) == 0 {
		return nil
	}

	var res []Posting

	aCur := 0
	bCur := 0

	for {
		if aCur > len(a)-1 || bCur > len(b)-1 {
			break
		}

		if a[aCur].DocID == b[bCur].DocID {
			f := a[aCur].Freq
			bf := b[bCur].Freq
			if bf < f {
				f = bf
			}
			res = append(res, Posting{
				DocID: a[aCur].DocID,
				Freq:  f,
			})
			aCur++
			bCur++
			continue
		}

		if a[aCur].DocID > b[bCur].DocID {
			bCur++
		} else {
			aCur++
		}
	}
	return res
}

// Phrase search for an exact phrase and returns all matching documents.
func (q *querier) Phrase(phrase string) ([]Posting, error) {
	tokens := strings.Split(phrase, " ")
	if len(tokens) == 0 {
		return nil, nil
	}

	res, err := q.idx.Postings(tokens[0])
	if err != nil {
		return nil, fmt.Errorf("first postings: %w", err)
	}

	for _, t := range tokens[1:] {
		postingsList, err := q.idx.Postings(t)
		if err != nil {
			return nil, fmt.Errorf("postings: %w", err)
		}
		res = q.phrase(res, postingsList)
	}

	return res, nil
}

func (q *querier) phrase(a, b []Posting) []Posting {
	if len(a) == 0 {
		return nil
	}
	if len(b) == 0 {
		return nil
	}

	var res []Posting

	i := 0
	j := 0

	for i <= len(a)-1 && j <= len(b)-1 {
		if a[i].DocID > b[j].DocID {
			j++
			continue
		}
		if a[i].DocID < b[j].DocID {
			i++
			continue
		}

		var positions []int
		for _, ap := range a[i].Positions {
			for _, bp := range b[j].Positions {
				if ap+1 == bp {
					positions = append(positions, bp)
				}
			}

		}
		if len(positions) != 0 {
			f := a[i].Freq
			bf := b[j].Freq
			if bf < f {
				f = bf
			}
			res = append(res, Posting{
				DocID:     a[i].DocID,
				Freq:      f,
				Positions: positions,
			})
		}
		i++
		j++
	}
	return res
}
