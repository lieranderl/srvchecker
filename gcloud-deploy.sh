gcloud config set project srvchecker-collab
gcloud functions deploy svrprocess-service \
    --gen2 \
    --region=europe-west1 \
    --runtime=go121 \
    --source=. \
    --entry-point=Srvprocess \
    --memory=128Mi \
    --trigger-http \
    --allow-unauthenticated