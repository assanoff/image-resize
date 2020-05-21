#docker build -t resizer .
#docker run -p 3220:8080 --name resizer --rm -d -t resizer

# Initial stage: download modules
FROM golang:1.13 as modules
ADD go.mod go.sum /m/
RUN cd /m && go mod download

# Intermediate stage: Build the binary
FROM golang:1.13 AS builder
COPY --from=modules /go/pkg /go/pkg
RUN mkdir /app
ADD . /app
WORKDIR /app
# We want to build our application's binary executable
RUN CGO_ENABLED=0 GOOS=linux make build

# Final stage: Run the binary
FROM alpine:latest AS production
COPY --from=builder /app/bin/ .
CMD ["./imageresize"]