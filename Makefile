help:
	@cat Makefile | grep '^[a-z]'

deploy:
	rsync -r mysql $(MONEYTEMAA_DEPLOY_TARGET)
	rsync -r bin $(MONEYTEMAA_DEPLOY_TARGET)
	rsync .direnv $(MONEYTEMAA_DEPLOY_TARGET)/bin/

build:
	go build -o bin/api ./src

run:
	go run ./src
