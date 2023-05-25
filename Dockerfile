FROM docker.io/golang:1.20 as build

WORKDIR /tmp/src

COPY . .

RUN go build -o ./out/chart-verifier main.go

FROM registry.access.redhat.com/ubi8/ubi-minimal

COPY --from=build /tmp/src/out/chart-verifier /app/chart-verifier

WORKDIR /app

ENV PATH "$PATH:/app"

ENTRYPOINT ["/app/chart-verifier"]
