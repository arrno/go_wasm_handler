export GCLOUD_PROJECT="rebut-116bb" 
export REPO="wasm"
export REGION="us-east4"
export IMAGE="go-wasm"

export IMAGE_TAG=${REGION}-docker.pkg.dev/$GCLOUD_PROJECT/$REPO/$IMAGE

docker build -t $IMAGE_TAG -f ./Dockerfile --platform linux/x86_64 .
docker push $IMAGE_TAG