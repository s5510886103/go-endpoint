# Dockerfile
FROM golang:1.18 as build

RUN mkdir /app
WORKDIR /app
COPY ./go.mod .
COPY ./go.sum .
RUN go mod download
COPY . .
RUN go build -o /build/app .

FROM golang:1.18 as run
COPY --from=build /build/app /app
ENTRYPOINT [ "/app" ]
