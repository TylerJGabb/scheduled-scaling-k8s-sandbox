. ../.env
echo $FQIN
docker build --platform=linux/amd64 -t watching-scaler .
docker tag watching-scaler $FQIN:v7
docker push $FQIN:v7