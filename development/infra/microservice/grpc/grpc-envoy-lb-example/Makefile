.PHONY: proto build push docker-cluster docker-test k8s-cluster k8s-export k8s-test clean

proto:
	protoc -I=./proto --go_out=./server/pb --go_opt=paths=source_relative --go-grpc_out=./server/pb --go-grpc_opt=paths=source_relative proto/service.proto
	protoc -I=./proto --go_out=./client/pb --go_opt=paths=source_relative --go-grpc_out=./client/pb --go-grpc_opt=paths=source_relative proto/service.proto

build:
	docker build -t rche/grpc-greeter -f server/Dockerfile ./server

push: build
	docker push rche/grpc-greeter

docker-cluster:
	docker-compose up --build -d

docker-test:
	cd client; go run ./main.go --addr 127.0.0.1:10000

k8s-cluster:
	kubectl create configmap envoy-config \
		--from-file=envoy/envoy.yaml\
		--from-file=envoy/lds.yaml \
		--from-file=envoy/cds.yaml \
		-o yaml --dry-run=client | kubectl apply -f -
	kubectl apply -f k8s.yaml

k8s-export:
	kubectl port-forward  svc/envoy 10000:10000

k8s-test:
	cd client; go run ./main.go --addr 127.0.0.1:10000

clean:
	docker-compose down
