docker build -f Dockerfile.multistage -t cf-temporal-worker:multistage .
docker run cf-temporal-worker:multistage
