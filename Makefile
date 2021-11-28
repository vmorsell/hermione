LINT_VERSION := 2021.1.1
LINT_FLAGS := -checks inherit

.PHONY: test lintinstall lint build

test:
	go test ./...

lintinstall:
	go install honnef.co/go/tools/cmd/staticcheck@$(LINT_VERSION)

lint:
	staticcheck $(LINT_FLAGS) ./...

build:
	go build ./...

corpus:
	wget http://www.csc.kth.se/~jboye/teaching/ir/davisWiki.tar.gz -P corpus
	(cd corpus && tar -xf davisWiki.tar.gz)

ingestdaviswiki: corpus
	go run ./ingest/ingest.go ./corpus/davisWiki
