IMAGENAME = dcos-flink:v2.1


build: build-flink, build-container

build-container:
	docker build -t $(IMAGENAME) .

build-flink:
	cd flink; mvn package -DskipTests

push: build
	docker push $(IMAGENAME)
