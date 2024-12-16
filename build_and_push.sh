docker build --platform=linux/amd64 -t scaling-sandbox .
docker tag scaling-sandbox us-east1-docker.pkg.dev/dv01-prj-shared-art-reg-6j49/dv01-shared-art-reg-01/scaling-sandbox:v1
docker push us-east1-docker.pkg.dev/dv01-prj-shared-art-reg-6j49/dv01-shared-art-reg-01/scaling-sandbox:v1