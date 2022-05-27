# build server stage
FROM golang:1.17-alpine3.15 as build_server
WORKDIR /server
COPY . .
RUN go build -o /bin/app .

# build runtime stage
FROM alpine:3.15
WORKDIR /server
COPY --from=build_server /bin/app /server/app
EXPOSE 8080
CMD ["/server/app"]
