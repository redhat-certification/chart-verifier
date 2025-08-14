FROM docker.io/golang:1.24 as build

WORKDIR /tmp/src

COPY . .

RUN make bin

FROM registry.access.redhat.com/ubi9/ubi-minimal

COPY --from=build /tmp/src/out/chart-verifier /app/chart-verifier

WORKDIR /app

ENV PATH "$PATH:/app"

ENTRYPOINT ["/app/chart-verifier"]
