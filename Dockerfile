FROM golang:alpine as builder
ADD . /src
WORKDIR /src
RUN go build -o entain-task cmd/server/main.go

FROM alpine
WORKDIR /app

COPY --from=builder /src/entain-task /app/
ADD  /migrations /app/migrations

ENTRYPOINT /app/entain-task --env-file ./my_env.list