# Dockerfile
FROM golang:1.18 as build
#FROM golang:1.18-alpine as build

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
COPY .env ./
# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOAMD64=v3 go build -a -installsuffix cgo -o /build/app

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
#FROM scratch
COPY --from=build /build/app .
COPY .env .
EXPOSE 8080
CMD [ "./app" ]
