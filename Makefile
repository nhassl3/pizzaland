.PHONY: build buikd-run run clean migrate -DEFAULT-GOAL migrate-test test

BINARY_NAME := pizzaland
BUILD_DIR := build
MAIN_PACKAGE := ./cmd/pizzaland # Adjust if your main package is elsewhere

# Target to build the Go application
build:
	@mkdir -p $(BUILD_DIR)
	@GOOS=$(shell go env GOOS) GOARCH=$(shell go env GOARCH) go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)

# Target to run the application
build-run: build
	@$(BUILD_DIR)/$(BINARY_NAME) --config="./config/local_config.yaml"

run:
	@$(BUILD_DIR)/$(BINARY_NAME) --config="./config/local_config.yaml"

clean:
	@rm -rf $(BUILD_DIR)

migrate:
	@if [ "$(word 2, $(MAKECMDGOALS))" = "down" ]; then \
		echo "Running migrations with down direction"; \
		go run ./cmd/migrator/ --storage-path="./storage/$(BINARY_NAME).db" --migrations-path="./migrations/" --down=true; \
	elif [ "$(word 2, $(MAKECMDGOALS))" = "up"]; then \
		echo "Running migrations with up direction"; \
		go run ./cmd/migrator/ --storage-path="./storage/$(BINARY_NAME).db" --migrations-path="./migrations/" --down=false; \
	else \
	  	echo "Running migrations with up direction"; \
      	go run ./cmd/migrator/ --storage-path="./storage/$(BINARY_NAME).db" --migrations-path="./migrations/" --down=false; \
	fi

migrate-test:
	@if [ "$(word 2, $(MAKECMDGOALS))" = "down" ]; then \
    		echo "Running migrations with down direction"; \
    		go run ./cmd/migrator/ --storage-path="./storage/$(BINARY_NAME).db" --migrations-path="./tests/migrations" --migrations-table=migrations_test --down=true; \
    	elif [ "$(word 2, $(MAKECMDGOALS))" = "up"]; then \
    		echo "Running migrations with up direction"; \
    		go run ./cmd/migrator/ --storage-path="./storage/$(BINARY_NAME).db" --migrations-path="./tests/migrations" --migrations-table=migrations_test --down=false; \
    	else \
    	  	echo "Running migrations with up direction"; \
          	go run ./cmd/migrator/ --storage-path="./storage/$(BINARY_NAME).db" --migrations-path="./tests/migrations" --migrations-table=migrations_test --down=false; \
    	fi

test:
	@go test ./tests
# Игнорируем аргументы как цели
%:
	@:

-DEFAULT-GOAL: run