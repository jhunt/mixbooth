IMAGE ?= filefrog/mixbooth
TAG   ?= latest

default:
	go build .
	./mixbooth

docker:
	docker build -t $(IMAGE):$(TAG) .

push: docker
	docker push $(IMAGE):$(TAG)
