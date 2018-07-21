all: clean build
build: 
		CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o camel-go
clean: 
		go clean
		rm -f camel-go

docker:
		docker build --tag="lburgazzoli/camel-go" .

dokerrun:
		docker run \
			--rm \
			-ti \
			-v ${PWD}/examples/example-flow/flow-simple.yaml:/home/camel/flow.yaml:Z \
			lburgazzoli/camel-go:latest \
				run