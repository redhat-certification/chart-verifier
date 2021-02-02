$COMMIT_ID = $(git rev-parse --short HEAD)
docker build -t helmcertifier:$COMMIT_ID .