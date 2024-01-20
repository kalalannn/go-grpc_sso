_mkdir_bin:
	mkdir -p bin
build: _mkdir_bin
	go build -o bin/sso_app cmd/server/main.go
run: build
	./bin/sso_app

generate:
	protoc --go_out=. --go-grpc_out=. proto/sso.proto

init_db:
	touch data/app.db

rm_db:
	-rm data/app.db

reinit_db: rm_db init_db

clean_users:
	sqlite3 data/app.db "DELETE FROM users;" ".exit"
