export GOPATH := /var/www/smawk-bot

default:
	@go get gopkg.in/telegram-bot-api.v4

build:
	@go build

deploy:
	@openssl req -newkey rsa:2048 -sha256 -nodes -keyout smawk_key.pem -x509 -days 365 -out smawk_cert.pem -subj "/C=US/ST=South Carolina/L=Lexington/O=My Simple Things/CN=mysimplethings.xyz"
	@go build
	@mv smawk-bot /var/www/smawk-bot
	@mv smawk_cert.pem /var/www/smawk-bot
	@mv smawk_key.pem /var/www/smawk-bot
	@cd /var/www/smawk-bot
	./smawk-bot

makecert:
	@openssl req -newkey rsa:2048 -sha256 -nodes -keyout smawk_key.pem -x509 -days 365 -out smawk_cert.pem -subj "/C=US/ST=South Carolina/L=Lexington/O=My Simple Things/CN=mysimplethings.xyz"

deploycert:
	@openssl req -newkey rsa:2048 -sha256 -nodes -keyout smawk_key.pem -x509 -days 365 -out smawk_cert.pem -subj "/C=US/ST=South Carolina/L=Lexington/O=My Simple Things/CN=mysimplethings.xyz"
	@mv smawk_cert.pem /var/www/smawk-bot
	@mv smawk_key.pem /var/www/smawk-bot
