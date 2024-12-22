. ../.env
TAG=v9
echo $FQIN:$TAG
docker build --platform=linux/amd64 -t watching-scaler .
docker tag watching-scaler $FQIN:$TAG
docker tag watching-scaler $FQIN:latest
docker push $FQIN:$TAG
docker push $FQIN:latest