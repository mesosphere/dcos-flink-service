IMAGENAME = dcos-flink:1.10.2-1.0

build: build-flink build-container

build-container:
	docker build -t $(IMAGENAME) .

build-flink:
	cd flink;  mvn clean package -DskipTests -Pfast

push: build
	docker push $(IMAGENAME)
