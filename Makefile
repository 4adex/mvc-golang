migrate_up_step:
	@read -p "Enter MySQL username: " DB_USER; \
	read -sp "Enter MySQL password: " DB_PASS; \
	echo; \
	read -p "Enter MySQL port (default 3306): " DB_PORT; \
	DB_PORT=$${DB_PORT:-3306}; \
	migrate -path database/migration/ -database "mysql://$$DB_USER:$$DB_PASS@tcp(localhost:$$DB_PORT)/library_management" -verbose up 1

migrate_down_step:
	@read -p "Enter MySQL username: " DB_USER; \
	read -sp "Enter MySQL password: " DB_PASS; \
	echo; \
	read -p "Enter MySQL port (default 3306): " DB_PORT; \
	DB_PORT=$${DB_PORT:-3306}; \
	migrate -path database/migration/ -database "mysql://$$DB_USER:$$DB_PASS@tcp(localhost:$$DB_PORT)/library_management" -verbose down 1

migrate_up:
	@read -p "Enter MySQL username: " DB_USER; \
	read -sp "Enter MySQL password: " DB_PASS; \
	echo; \
	read -p "Enter MySQL port (default 3306): " DB_PORT; \
	DB_PORT=$${DB_PORT:-3306}; \
	migrate -path database/migration/ -database "mysql://$$DB_USER:$$DB_PASS@tcp(localhost:$$DB_PORT)/library_management" -verbose up

migrate_down:
	@read -p "Enter MySQL username: " DB_USER; \
	read -sp "Enter MySQL password: " DB_PASS; \
	echo; \
	read -p "Enter MySQL port (default 3306): " DB_PORT; \
	DB_PORT=$${DB_PORT:-3306}; \
	migrate -path database/migration/ -database "mysql://$$DB_USER:$$DB_PASS@tcp(localhost:$$DB_PORT)/library_management" -verbose down