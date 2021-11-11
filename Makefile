build:
	go build -o ./apisix-opa-plugin ./

docker:
	docker build -t apisix-opa-plugin .

build-linux:
	CGO_ENABLED=0 GOOS=linux go build -o ./apisix-opa-plugin ./