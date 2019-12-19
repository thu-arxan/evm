# Go parameters
GOCMD=go

# go vet:format check, bug check
vet:
	$(GOCMD) vet `go list ./...`
