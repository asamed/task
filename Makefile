all: build run

build:
	@docker build . -t api

run:
	# @docker run --name mongodb -p 27017:27017 -d mongo
	@docker run -it -p 8082:8082 api