docker login
docker build -f Dockerfile.multistage -t sharroh/cf-temporal-worker:multistage . --no-cache
docker run sharroh/cf-temporal-worker:multistage
