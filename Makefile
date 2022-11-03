environment:
	# docker-compose build --no-cache
	docker build -t service-build docker/grpc/build

run/protoc:
ifeq ($(OUT),null)
	echo "エラー: 'OUT' を指定する必要があります";
	exit 1;
endif
	docker run --rm \
		--mount type=bind,source=${PWD}/docker/grpc/build,target=/mnt/out \
		-v ${PWD}/service.proto:/mnt/service/service.proto \
		service-build \
		sh -c "\
			mkdir -p /mnt/out/protobuf; \
			mkdir -p /mnt/out/doc; \
			protoc -I=/mnt/service service.proto \
			--go_out=plugins=grpc:/mnt/out/protobuf \
			--go_opt=paths=source_relative \
			--go_opt=Mservice.proto=service/infrastructure/service\
			--doc_out=/mnt/out/doc --doc_opt=markdown,service.md \
		"
	(mkdir $(OUT)/protobuf || true) > /dev/null 2>&1
ifeq ($(shell uname),Linux)
	sudo chown -R $(USER):$(USER) ./docker/grpc/build/
endif
	mv ${PWD}/docker/grpc/build/protobuf/* $(OUT)
	mv ${PWD}/docker/grpc/build/doc/* ${PWD}
	rmdir ${PWD}/docker/grpc/build/protobuf
	rmdir ${PWD}/docker/grpc/build/doc