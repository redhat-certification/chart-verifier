FROM golang:1.15 AS build

WORKDIR /tmp/src

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN ./hack/build.sh

FROM registry.access.redhat.com/ubi8/ubi-minimal

COPY --from=build /tmp/src/out/chart-verifier /app/chart-verifier

ENTRYPOINT ["/app/chart-verifier"]
