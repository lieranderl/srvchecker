# Makefile for deploying the gcloud function

# Variables
RUNTIME = go123
SOURCE = .
ENTRY_POINT = Srvprocess
MEMORY = 128Mi

# Default target
.PHONY: deploy
deploy-gcloud:
	@echo "Deploying $(FUNCTION_NAME) to Google Cloud Functions..."
	gcloud config set project $(PROJECT)
	gcloud functions deploy $(FUNCTION_NAME) \
	    --gen2 \
	    --region=$(REGION) \
	    --runtime=$(RUNTIME) \
	    --source=$(SOURCE) \
	    --entry-point=$(ENTRY_POINT) \
	    --memory=$(MEMORY) \
	    --trigger-http \
	    --allow-unauthenticated
	@echo "Deployment complete."

# Clean target (optional)
.PHONY: clean
clean:
	@echo "Cleaning up..."
	# Add any cleanup commands here



## login to AWS ECR 
# aws configure sso --session-name jfedotov
# aws ecr get-login-password --region  eu-west-1 --profile jfedotov | docker login --username AWS --password-stdin 585768188093.dkr.ecr.eu-west-1.amazonaws.com
# aws ecr describe-repositories --region eu-west-1 --profile jfedotov
# aws ecr create-repository --repository-name $(FUNCTION_NAME) --region eu-west-1 --profile jfedotov

## clean docker image
cleanaws:
	@echo "Cleaning up..."
	docker rmi  $(AWS_ECR)/$(FUNCTION_NAME):$(TAG)

## build docker image
buildaws:
	@echo "Building docker image..."
	docker buildx build --platform linux/amd64 --provenance false --load -t  $(AWS_ECR)/$(FUNCTION_NAME):$(TAG) .


## push docker image to AWS ECR
pushaws:
	@echo "Pushing docker image to AWS ECR..."
	docker push $(AWS_ECR)/$(FUNCTION_NAME):$(TAG)

## deploy to AWS lambda
deployaws:
	@echo "Deploying $(FUNCTION_NAME) to AWS Lambda..."
	aws lambda update-function-code \
		--function-name $(FUNCTION_NAME) \
		--image-uri  $(AWS_ECR)/$(FUNCTION_NAME):$(TAG) \
		--region eu-west-1 \
		--profile jfedotov
	@echo "Deployment complete."