FROM golang:1.15.3-alpine3.12 as builder

WORKDIR /app

COPY go.mod ./go.mod
COPY go.sum ./go.sum

RUN go mod download

COPY config ./config
COPY models ./models
COPY runner ./runner
COPY store ./store
COPY utils ./utils
COPY main.go ./main.go

RUN go build -o bin/codacy-usage-report main.go

FROM alpine:3.12

COPY --from=builder /app/bin/codacy-usage-report /app/

USER nobody:nobody

WORKDIR /app

ENTRYPOINT [ "/app/codacy-usage-report" ]
