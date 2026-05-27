MODULE := build-distributed-systems

.PHONY: help run submit verify-submit clean

help:
	@echo "Targets:"
	@echo "  run CHALLENGE=<path>                 Run a challenge (reads stdin)"
	@echo "  submit CHALLENGE=<path>              Generate paste-ready _submit.go"
	@echo "  verify-submit CHALLENGE=<path>       Confirm _submit.go compiles standalone"
	@echo "  clean                                Remove generated artifacts"

BUNDLE    := go run golang.org/x/tools/cmd/bundle@latest
GOIMPORTS := go run golang.org/x/tools/cmd/goimports@latest


run:
	@test -n "$(CHALLENGE)" || (echo "usage: make run CHALLENGE=<path>"; exit 1)
	@go run ./challenges/$(CHALLENGE)

submit:
	@test -n "$(CHALLENGE)" || (echo "usage: make submit CHALLENGE=<path>"; exit 1)
	@mkdir -p .tmp
	@$(BUNDLE) -pkg main -prefix '' -dst ./challenges/$(CHALLENGE) ./internal/core > .tmp/core_bundle.go
	@awk '/^package main/ {next} /^import \(/,/^\)/ {next} /^import "/ {next} {print}' \
	  challenges/$(CHALLENGE)/main.go > .tmp/main_body.go
	@cat .tmp/core_bundle.go .tmp/main_body.go > challenges/$(CHALLENGE)/_submit.go
	@$(GOIMPORTS) -w challenges/$(CHALLENGE)/_submit.go
	@echo "✓ challenges/$(CHALLENGE)/_submit.go ready to paste"

verify-submit: submit
	@d=$$(mktemp -d) && cp challenges/$(CHALLENGE)/_submit.go $$d/main.go && \
	  (cd $$d && go mod init submittest >/dev/null 2>&1 && go build .) && \
	  echo "✓ _submit.go compiles standalone" && rm -rf $$d

clean:
	rm -rf .tmp bin
	find challenges -name '_submit.go' -delete
