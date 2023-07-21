# Application

build:
	CGO_ENABLED=0 go build -o gif-app

server-up: build
	./gif-app

curl-hello:
	curl -v -JLO localhost:8080/goCol?text=Hello,+We+Are+Devs

curl-jeans:
	curl -v -JLO localhost:8080/byBlink?text=Gopher+Jeans

curl-complex:
	curl -v -JLO localhost:8080/free?text=We+Are+Dev+World+Congress -d '{"delay":150, "figlet":"speed", "gifType":"alt","bgHex":"0x030303", "fgHex":"0x78C475"}'
	
# Docker

docker-push:
	docker build -f deploy/docker/Dockerfile -t mcaci/gif-app .
	docker push mcaci/gif-app

docker-up:
	docker run -it --detach --rm -p 8080:8080 --name gif-app mcaci/gif-app:latest

docker-down:
	docker stop gif-app

# Kubernetes

spin-cluster: cluster-up app-install

cluster-up:
	sh ./deploy/kind/setup.sh
	helm install metallb metallb/metallb

cluster-down:
	kind delete cluster --name go-debug-cluster

app-install:
	helm install gif-app ./deploy/gif-app

app-uninstall:
	helm uninstall gif-app

app-upgrade:
	helm upgrade gif-app ./deploy/gif-app

port-forward:
	kubectl port-forward service/gif-app 8080:8080

restart: app-uninstall build docker-push app-install