# Default target
all: help

# Display help
help:
	@echo "Makefile commands:"
	@echo "  make help      - Display this help message"
	@echo "  make deploy    - Deploy to AWS Lambda serverless"

build:
	@echo "Building..."
	sam build \
		--parameter-overrides ImageTag=$(TAG)

## deploy to AWS lambda
deploy:
	@echo "Deploying to AWS Lambda..."
	sam build \
		--parameter-overrides ImageTag=$(TAG) \
		&& \
	sam deploy \
		--parameter-overrides ImageTag=$(TAG)

local:
	@echo "Running locally..."
	sam local start-api -p 3001