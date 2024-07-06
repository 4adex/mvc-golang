say_hello:
	@echo "Hello World"

say_gull:
	@echo "Gulla"

migrate_up_step:
	migrate -path database/migration/ -database "mysql://root:1234@tcp(localhost:3306)/library_management" -verbose up 1

migrate_down_step:
	migrate -path database/migration/ -database "mysql://root:1234@tcp(localhost:3306)/library_management" -verbose down 1

migrate_up:
	migrate -path database/migration/ -database "mysql://root:1234@tcp(localhost:3306)/library_management" -verbose up

migrate_down:
	migrate -path database/migration/ -database "mysql://root:1234@tcp(localhost:3306)/library_management" -verbose down

