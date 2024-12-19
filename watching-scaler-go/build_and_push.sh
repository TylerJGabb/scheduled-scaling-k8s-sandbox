. ../.env
TAG=v8
echo $FQIN:$TAG
docker build --platform=linux/amd64 -t watching-scaler .
docker tag watching-scaler $FQIN:$TAG
docker push $FQIN:$TAG