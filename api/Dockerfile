FROM golang:1.19-alpine

WORKDIR /app
ADD . ./
# COPY go.mod ./
# COPY go.sum ./
RUN go mod download
# COPY api.go ./
# COPY cmd/ ./cmd/
# COPY pkg ./pkg/
RUN go build -o /main

# Remember to expose port with -p flag
EXPOSE 8080

CMD ["/main"]
