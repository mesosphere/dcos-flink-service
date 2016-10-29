IMAGENAME = dcos-flink

build:
	docker build -t $(IMAGENAME) .

push: build
	docker push $(IMAGENAME)
