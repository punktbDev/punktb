include .env

.DEFAULT_GOAL := build

save:
	docker save backend:${TAG} > ~/Downloads/backend.tar
build:
	docker-compose build
migrate_init:
	migrate create -ext sql -dir ./migrations -seq init
migrate_up:
	migrate -path ./migrations -database 'postgres://postgres:Ytn100vfys!@127.0.0.1:5432/postgres?sslmode=disable' up
migrate_down:
	migrate -path ./migrations -database 'postgres://postgres:Ytn100vfys!@127.0.0.1:5432/postgres?sslmode=disable' down
cert_gen:
	#FDQN это адрес сервера на котором будет работать бэк, например 127.0.0.1:8080
	openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout key.pem -out cert.pem