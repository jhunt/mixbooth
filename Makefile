IMAGE ?= filefrog/mixbooth
TAG   ?= latest

default:
	plackup bin/app.psgi

docker:
	docker build -t $(IMAGE):$(TAG) .

push: docker
	docker push $(IMAGE):$(TAG)
