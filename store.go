package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type Store interface {
	Get(id string) ([]byte, error)
	PutFromStream(r io.Reader, id string) error
}

type store struct {
	root string
}

func NewStore(root string) (Store, error) {
	if _, err := os.Stat(root); !os.IsNotExist(err) {
		if err := os.RemoveAll(root); err != nil {
			return nil, fmt.Errorf("remove root: %w", err)
		}
	}

	if err := os.Mkdir(root, 0755); err != nil {
		return nil, fmt.Errorf("mkdir: %w", err)
	}

	return &store{
		root: root,
	}, nil
}

func (s *store) Get(id string) ([]byte, error) {
	file, err := os.Open(fmt.Sprintf("%s/%s", s.root, id))
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	//bytes, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", s.root, id))
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}
	return bytes, nil
}

func (s *store) PutFromStream(r io.Reader, id string) error {
	log.Printf("id: %s\n", id)
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return fmt.Errorf("read all: %w", err)
	}

	file, err := os.Create(fmt.Sprintf("%s/%s", s.root, id))
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}
	defer file.Close()

	file.Write(bytes)
	return nil
}
