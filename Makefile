# Go parameters
GOCMD=go

all: vet test

# go vet:format check, bug check
vet:
	@$(GOCMD) vet `go list ./...`

# test:
test:
	@$(GOCMD) test -count=1 github.com/thu-arxan/evm/util
	@$(GOCMD) test -count=1 github.com/thu-arxan/evm/core
	@$(GOCMD) test -count=1 github.com/thu-arxan/evm/abi
	@$(GOCMD) test -count=1 github.com/thu-arxan/evm/tests

# sol will compile solidity code
sol:
	@-cd tests/sols && solcjs --bin *.sol
	@-cd tests/sols && solcjs --abi *.sol

generate:
	@$(GOCMD) generate