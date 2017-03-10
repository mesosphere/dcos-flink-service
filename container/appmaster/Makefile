IMAGENAME = dcos-flink:1.2.0-1.4
IMAGENAME_2_11 = dcos-flink-2-11:1.2.0-1.4


build: build-flink build-container build-flink-2-11 build-container-2-11

build-container:
	docker build -t $(IMAGENAME) .

build-container-2-11:
	docker build -t $(IMAGENAME_2_11) .

build-flink-2-11:
	cd flink; ./tools/change-scala-version.sh 2.11; mvn clean package -DskipTests

build-flink:
	cd flink; ./tools/change-scala-version.sh 2.10; mvn clean package -DskipTests

push: build
	docker push $(IMAGENAME)
	docker push $(IMAGENAME_2_11)
