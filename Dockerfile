FROM golang:1.17.3-alpine AS GO_BUILD
WORKDIR /api
COPY . /api
RUN go build -o /go/bin/api

FROM alpine:3.10
WORKDIR /api
COPY --from=GO_BUILD /go/bin/api/ ./
COPY .env ./
EXPOSE 8082
CMD ./api