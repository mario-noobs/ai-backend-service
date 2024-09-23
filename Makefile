.PHONY: run

# Default environment variables
SERVER_PORT ?= 8080
SERVER_HOST ?= localhost
LOG_LEVEL ?= info

run:
    @echo "Starting server with configuration:"
    @echo "HOST: $(SERVER_HOST)"
    @echo "PORT: $(SERVER_PORT)"
    @echo "LOG LEVEL: $(LOG_LEVEL)"
    SERVER_PORT=$(SERVER_PORT) SERVER_HOST=$(SERVER_HOST) LOG_LEVEL=$(LOG_LEVEL) go run main.go
