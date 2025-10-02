.PHONY: build
build:
	@echo "Setting up target"
	@mkdir -p target/bin
	@mkdir -p target/log
	@echo "=== START OF LOG ===" >> target/log/file.log
	@echo "Building project"
	@go build -o ./target/bin/main

.PHONY: confirm 
confirm:
	@echo -n "Are you sure? [y/N] " && read ans && [ $${ans:-N} = y ]

.PHONY: clean
clean: confirm
	@echo "Cleaning up ..."
	@rm -rf target


