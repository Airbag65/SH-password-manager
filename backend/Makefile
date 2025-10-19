.PHONY: build
build:
	@echo "Setting up target"
	@mkdir -p target/pem
	@touch ./target/pem/privateKey.pem
	@touch ./target/pem/publicKey.pem
	@mkdir -p target/log
	@mkdir -p target/bin
	@echo "=== START OF LOG ===" > ./target/log/file.log
	@echo "Building project"
	@go build -o ./target/bin/api

.PHONY: confirm 
confirm:
	@echo -n "Are you sure? [y/N] " && read ans && [ $${ans:-N} = y ]

.PHONY: clean
clean: confirm
	@echo "Cleaning up ..."
	@rm -rf ./target/log
	@rm -rf ./target/bin

.PHONY: run
run: build
	@echo "Running project"
	@./target/bin/api
