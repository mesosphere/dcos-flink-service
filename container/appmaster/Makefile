IMAGENAME = dcos-flink:1.4.0-1.0

build: build-flink build-container

build-container:
	docker build -t $(IMAGENAME) .

build-flink:
	cd flink;  mvn clean package -DskipTests

push: build
	docker push $(IMAGENAME)
