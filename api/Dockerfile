FROM golang:1.19-alpine

WORKDIR /app
ADD . ./
RUN go mod download
RUN go build -o /main

EXPOSE 8080
# ENV host=terraform-20221229205211947800000001.clgu6fwjxdcc.me-south-1.rds.amazonaws.com
# ENV redis=redis:6379
# ENV database=testdb
# ENV username=postgres
# ENV password=postgres
# ENV frontendOrigin=http://localhost:3000
# 
# CMD ["/main"]
