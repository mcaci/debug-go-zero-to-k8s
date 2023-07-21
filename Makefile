build:
	CGO_ENABLED=0 go build -o gif-app

app-install:
	helm install gif-app ./deploy/gif-app

app-uninstall:
	helm uninstall gif-app

cluster-down:
	kind delete cluster --name go-debug-cluster

docker-publish:
	docker build -f deploy/docker/Dockerfile -t localhost:5001/mcaci/gif-app .
	docker push localhost:5001/mcaci/gif-app

cluster-up:
	sh ./deploy/kind/setup.sh
	helm install metallb metallb/metallb