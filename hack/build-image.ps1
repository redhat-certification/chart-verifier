$COMMIT_ID = $(git rev-parse --short HEAD)
docker build -t quay.io/redhat-certification/chart-verifier:$COMMIT_ID .
