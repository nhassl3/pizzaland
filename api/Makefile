.PHONY: gen install gen-validate genall -DEFAULT-GOAL

gen:
	@protoc -I proto proto/pizzaland/pizzaland.proto --go_out=./generated/go --go_opt=paths=source_relative --go-grpc_out=./generated/go/ --go-grpc-opt_out=paths=source_relative

install:
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go install github.com/envoyproxy/protoc-gen-validate@latest

gen-validate:
	@protoc \
	  -I \
	  proto \
	  proto/pizzaland/pizzaland.proto \
	  --go_out=./generated/go \
	  --go_opt=paths=source_relative \
	  --go-grpc_out=./generated/go \
	  --go-grpc_opt=paths=source_relative \
	  --validate_out="lang=go:./generated/go" \
	  --validate_opt=paths=source_relative

genall: install gen-validate

-DEFAULT-GOAL: genall
