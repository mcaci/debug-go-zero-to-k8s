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

docker-push-debug:
	docker build -f deploy/docker/debug.Dockerfile -t mcaci/gif-app:debug .
	docker push mcaci/gif-app:debug

docker-debug:
	docker ps -f name=gif-app
	docker images mcaci/gif-app
	docker run -it --rm --cap-add=SYS_PTRACE --pid="container:gif-app" mcaci/gif-app:debug 

docker-up:
	docker run -it --detach --rm -p 8080:8080 --name gif-app mcaci/gif-app:latest

docker-down:
	docker stop gif-app

# Kubernetes

cluster-up:
	sh ./deploy/kind/setup.sh
	helm install metallb metallb/metallb
	kubectl wait deployment.apps/metallb-controller --for=condition=ready --timeout=5m

cluster-down:
	kind delete cluster --name go-debug-cluster

app-install:
	helm install gif-app ./deploy/gif-app

app-uninstall:
	helm uninstall gif-app

app-upgrade:
	helm upgrade gif-app ./deploy/gif-app

curl-hello-k8s:
	curl -v -JLO ${IP}:8080/goCol?text=Hello,+Gophers

curl-jeans-k8s:
	curl -v -JLO ${IP}:8080/byBlink?text=Gopher+Jeans

curl-complex-k8s:
	curl -v -JLO ${IP}:8080/free?text=Container+Days -d '{"delay":150, "figlet":"speed", "gifType":"alt","bgHex":"0x030303", "fgHex":"0x78C475"}'

app-debug-k8s:
	kubectl debug ${POD} -it --image mcaci/gif-app:debug --profile=general --target=gif-app

app-debug-k8s-copy:
	kubectl debug ${POD} -it --image mcaci/gif-app:debug --profile=general --share-processes --copy-to=gif-app-debug
	