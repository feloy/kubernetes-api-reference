FROM golang:1.13-buster AS builder

WORKDIR /app

ADD go.mod go.sum /app/
RUN go mod download

ADD api /app/api
ADD pkg /app/pkg
ADD cmd /app/cmd

# Run unit tests
RUN go test ./...

# Build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /kubernetes-api-reference cmd/main.go

# Verify that resourceslist.txt is up to date
RUN /kubernetes-api-reference resourceslist -f api/v1.18/swagger.json > /tmp/resourceslist.txt
RUN diff /tmp/resourceslist.txt api/v1.18/resourceslist.txt

# final stage
FROM alpine:latest
COPY --from=builder /kubernetes-api-reference ./
RUN chmod +x ./kubernetes-api-reference
ENTRYPOINT ["./kubernetes-api-reference"]
