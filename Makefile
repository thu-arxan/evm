# Go parameters
GOCMD=go

# go vet:format check, bug check
vet:
	$(GOCMD) vet `go list ./...`

# test:
test:
	@$(GOCMD) test -count=1 evm/util
	@$(GOCMD) test -count=1 evm/core
	@$(GOCMD) test -count=1 evm/abi
	@$(GOCMD) test -count=1 evm/eabi
	@$(GOCMD) test -count=1 evm/tests