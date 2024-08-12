startdb:
	 docker run -d --rm --name snippetboxdb -v snippetboxdb:/var/lib/mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=root mysql:latest

stopdb:
	docker stop snippetboxdb

migrate:
	docker run -v ./migrations:/migrations migrate/migrate \
    -path=/migrations -database "mysql://web:web@tcp(docker.for.mac.localhost:3306)/snippetbox" $(direction)

makemigrations:
	docker run  -v ./migrations:/migrations migrate/migrate create -ext=".sql" -dir="./migrations" $(name)

makemigrations-test:
	docker run  -v ./internal/tests/migrations:/migrations migrate/migrate create -ext=".sql" -dir="./migrations" $(name)

migrate-test:
	docker run -v ./internal/tests/migrations:/migrations migrate/migrate \
	-path=/migrations -database "mysql://web_test:web_test@tcp(docker.for.mac.localhost:3306)/snippetbox_test" $(direction) $(flags)

run:
	go run ./cmd/web -config="./config/local.yaml"