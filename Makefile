all: build run

build:
	@docker build . -t api

run:
	@docker run -t mongodb -p 27017:27017 -d mongo:latest
	@docker run -it --name api -p 8082:8082 --rm api
