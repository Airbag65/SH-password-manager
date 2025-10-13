.PHONY: confirm 
confirm:
	@echo "Preparing clean up"
	@echo -n "Are you sure you want to clean up project? [y/N] " && read ans && [ $${ans:-N} = y ]

.PHONY: clean
clean: confirm
	@echo "Cleaning up ..."
	@rm -rf target

.PHONY: build
build: 
	@echo "Setting up target"
	@mkdir -p target/bin
	# @touch auth/auth.json
	# @echo '{"auth_token": ""}' > auth/auth.json
	@echo "Building project"
	@go mod tidy
	@go build -o ./target/bin/main

.PHONY: run
run: clean build
	@echo "Running project\n\n\n"
	@./target/bin/main


