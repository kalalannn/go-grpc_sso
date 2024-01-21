#! Local
_mkdir_bin:
	mkdir -p bin
build: _mkdir_bin
	go build -o bin/sso_app cmd/server/main.go
run: build
	./bin/sso_app

#! Docker
build_app_image:
	docker build -t go-grpc_sso-app .
run_app_image: build_app_image
	docker run -d -e APP_CONFIG_PATH=/config/docker.yaml -p 5051:5051 --name go-grpc_sso-app-local go-grpc_sso-app

rm_app_image:
	docker rm -f go-grpc_sso-app-local

sh_app_image:
	docker run --rm -it --entrypoint /bin/bash go-grpc_sso-app

#! Protobuf
generate:
	protoc --go_out=. --go-grpc_out=. proto/sso.proto

#! Database
init_db:
	touch data/app.db

rm_db:
	-rm data/app.db

reinit_db: rm_db init_db

clean_users:
	sqlite3 data/app.db "DELETE FROM users;" ".exit"
