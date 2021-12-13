PORT=8082

build:
	@docker build . -t api

run:
	@docker run -t mongodb -p 27017:27017 -d mongo:latest
	@docker run -it -p PORT:PORT --rm api