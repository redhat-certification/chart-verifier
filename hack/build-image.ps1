$COMMIT_ID = $(git rev-parse --short HEAD)
docker build -t chart-verifier:$COMMIT_ID .
