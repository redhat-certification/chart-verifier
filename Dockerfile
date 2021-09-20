#
# The golang:1.16 image is a copy of docker.io/library/golang:1.16 hosted in Quay to work around rate limits to
# Dockerhub:
#
# > docker pull golang:1.16
# > docker tag golang:1.16 quay.io/redhat-certification/golang:1.16
# > docker push quay.io/redhat-certification/golang:1.16
#
# To upgrade Go, then a new image should be pushed to Quay and updated below.
#
FROM quay.io/redhat-certification/golang:1.16 as build

WORKDIR /tmp/src

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN ./hack/get-oc.sh

RUN ./hack/get-helm.sh

RUN ./hack/build.sh

FROM registry.access.redhat.com/ubi8/ubi-minimal

COPY --from=build /tmp/src/out/chart-verifier /app/chart-verifier

COPY --from=build /usr/local/bin/* /usr/local/bin/

COPY ./config /app/config

COPY cmd/release /app/releases

ENTRYPOINT ["/app/chart-verifier"]
