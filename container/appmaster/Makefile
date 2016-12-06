IMAGENAME = dcos-flink:v1.1


build: build-flink
	docker build -t $(IMAGENAME) .

build-flink:
	cd flink; mvn package -DskipTests

push: build
	docker push $(IMAGENAME)
