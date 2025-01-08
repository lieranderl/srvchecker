# Default target
all: help

# Display help
help:
	@echo "Makefile commands:"
	@echo "  make help      - Display this help message"
	@echo "  make deploy    - Deploy to AWS Lambda serverless"

## deploy to AWS lambda
# Default to dev environment if not specified
ENV ?= dev

# Target to build and deploy
deploy:
	sam build
	sam deploy --config-env $(ENV)

local:
	@echo "Running locally..."
	sam local start-api -p 3001