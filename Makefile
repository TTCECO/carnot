include Makefile.ledger
all: install
install: go.sum
		GO111MODULE=on go install -tags "$(build_tags)" ./cmd/carnot
		GO111MODULE=on go install -tags "$(build_tags)" ./cmd/carnotcli
cd:
		GO111MODULE=on go install ./x/tcchan/contract/carnotcd
go.sum: go.mod
		@echo "--> Ensure dependencies have not been modified"
		GO111MODULE=on go mod verify
