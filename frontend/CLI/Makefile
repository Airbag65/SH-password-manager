.PHONY: build
build:
	@echo "Setting up target"
	@mkdir -p target/bin
	@echo "Building project"
	@go build -o ./target/bin/main

.PHONY: confirm 
confirm:
	@echo "Preparing clean up"
	@echo -n "Are you sure you want to clean up project? [y/N] " && read ans && [ $${ans:-N} = y ]

.PHONY: clean
clean: confirm
	@echo "Cleaning up ..."
	@rm -rf target


.PHONY: run
run: clean build
	@echo "Running project"
	@./target/bin/main
	
