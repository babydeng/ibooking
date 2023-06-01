# database name
DB_NAME ?= restapi

# database type
DB_TYPE ?= postgres

# database username
DB_USER ?= pi

# database password
DB_PWD ?= 123456


.PHONY : postgresup postgresdown createdb


# psql postgres://pi:123456@localhost:5432/restapi -c "\i schema.sql"



postgresup:
	docker run --name postgredb -e POSTGRES_USER=$(DB_USER) -e POSTGRES_PASSWORD=$(DB_PWD) -e POSTGRES_DB=$(DB_NAME) -p 5432:5432 -v ./postgres:/var/lib/postgresql/data/ -d postgres


postgresdown:
	docker stop postgredb  || true && 	docker rm postgredb || true

teardown_recreate: postgresdown postgresup
	sleep 5
	$(MAKE) createdb
