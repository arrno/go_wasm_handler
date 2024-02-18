  gcloud run deploy rebug-go-wasm \
    --region=us-east1 \
    --source=. \
    --service-account storage-admin@rebut-116bb.iam.gserviceaccount.com \
    --allow-unauthenticated