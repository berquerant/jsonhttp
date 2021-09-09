proto: pb/origin.proto
	protoc -I=./pb --go_out=. --go_opt=module=github.com/berquerant/jsonhttp ./pb/origin.proto

.PHONY: test
test:
	go test ./...

.PHONY: generate
generate:
	go generate ./...

.PHONY: build
build: clean proto generate
	mkdir -p dist
	go build -o dist/jsonhttp main.go

.PHONY: clean
clean:
	rm -rf dist
