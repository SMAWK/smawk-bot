export GOPATH := /var/www/smawk-bot

default:
	@go get gopkg.in/telegram-bot-api.v4

build:
	go build

test:
	go test

deploy:
	@openssl req -newkey rsa:2048 -sha256 -nodes -keyout smawk_key.pem -x509 -days 365 -out smawk_cert.pem -subj "/C=US/ST=South Carolina/L=Lexington/O=My Simple Things/CN=mysimplethings.xyz"
	go build

makecert:
	@openssl req -newkey rsa:2048 -sha256 -nodes -keyout smawk_key.pem -x509 -days 365 -out smawk_cert.pem -subj "/C=US/ST=South Carolina/L=Lexington/O=My Simple Things/CN=mysimplethings.xyz"
