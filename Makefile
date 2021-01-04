help:
	@cat Makefile | grep '^[a-z]'

deploy:
	rsync -r mysql $(MT_DEPLOY_TARGET)
	rsync -r bin $(MT_DEPLOY_TARGET)

build:
	go build -o bin/api ./src

run:
	go run ./src
