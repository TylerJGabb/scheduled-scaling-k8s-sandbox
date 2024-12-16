docker build --platform=linux/amd64 -t watching-scaler .
docker tag watching-scaler us-east1-docker.pkg.dev/dv01-prj-shared-art-reg-6j49/dv01-shared-art-reg-01/watching-scaler:v3
docker push us-east1-docker.pkg.dev/dv01-prj-shared-art-reg-6j49/dv01-shared-art-reg-01/watching-scaler:v3