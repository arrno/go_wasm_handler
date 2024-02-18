  gcloud functions deploy rebug-go-wasm \
    --gen2 \
    --runtime=go121 \
    --region=us-east1 \
    --source=. \
    --entry-point HandleWasm \
    --trigger-http \
    --allow-unauthenticated