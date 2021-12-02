FROM golang:1.17.3-alpine AS GO_BUILD
WORKDIR /api
COPY . /api
RUN go build -o /go/bin/api

FROM alpine:3.10
ENV MONGO_URI=mongodb://172.19.0.2 \
    DB=TestDB \
    PORT=8082
WORKDIR /api
COPY --from=GO_BUILD /go/bin/api/ ./
EXPOSE 8082
CMD ./api
