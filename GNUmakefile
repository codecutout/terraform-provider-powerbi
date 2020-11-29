.PHONY: build, testacc, fmt, fmtcheck, docs

default: build

build: fmtcheck
	@go build

testacc: fmtcheck
	@TF_ACC=1 go test -count=1 -v ./...

fmt:
	@gofmt -l -w $(CURDIR)/internal

fmtcheck:
	@test -z $(shell gofmt -l $(CURDIR)/internal | tee /dev/stderr) || { echo "[ERROR] Fix formatting issues with 'make fmt'"; exit 1; }

docs:
	@go run internal/docgen/cmd/main.go