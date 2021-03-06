PREFIX ?= filefrog
TAG    ?= latest

default:
	go build .
	./mixbooth

docker:
	docker build -t $(PREFIX)/mixbooth:$(TAG) .
	docker build -t $(PREFIX)/icecast2:$(TAG) icecast2/
	docker build -t $(PREFIX)/liquidsoap:$(TAG) liquidsoap/

push: docker
	docker push $(PREFIX)/mixbooth:$(TAG)
	docker push $(PREFIX)/icecast2:$(TAG)
	docker push $(PREFIX)/liquidsoap:$(TAG)
