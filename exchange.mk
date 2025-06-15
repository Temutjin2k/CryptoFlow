.PHONY: exchange-load-all exchange-run-all exchange-stop-all

exchange-load-all:
	docker load -i exchanges/exchange1_amd64.tar
	docker load -i exchanges/exchange2_amd64.tar
	docker load -i exchanges/exchange3_amd64.tar


exchange-run-all: exchange-load-all
	docker run -d --rm -p 40101:40101 --name exchange1_amd64 exchange1
	docker run -d --rm -p 40102:40102 --name exchange2_amd64 exchange2
	docker run -d --rm -p 40103:40103 --name exchange3_amd64 exchange3


exchange-stop-all:
	docker stop exchange1_amd64
	docker stop exchange2_amd64
	docker stop exchange3_amd64

