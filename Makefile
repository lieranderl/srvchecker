# Makefile for deploying the gcloud function

# Variables
RUNTIME = go123
SOURCE = .
ENTRY_POINT = Srvprocess
MEMORY = 128Mi

# Default target
.PHONY: deploy
deploy:
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
