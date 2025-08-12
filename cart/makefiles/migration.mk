DB_NAME := cart
DB_PASSWORD := cart1234
DB_USER := user_cart
POSTGRESQL_URL := "postgresql://${DB_USER}:${DB_PASSWORD}@localhost:5432/${DB_NAME}?sslmode=disable"
MIGRATE_VERSION := ${v}

.PHONY: migrate migrateup migratedown migratefix

migrate: ## 📁 Create a new migration file in internal/migrations (name: init_cart)
	migrate create -ext sql -dir internal/migrations -seq init_${DB_NAME}

migrateup: ## ⬆️ Apply all up migrations
	migrate -database ${POSTGRESQL_URL} -path internal/migrations -verbose up ${MIGRATE_VERSION}

migratedown: ## ⬇️ Roll back migrations
	migrate -database ${POSTGRESQL_URL} -path internal/migrations -verbose down ${MIGRATE_VERSION}

migratefix: ## 🛠 Force fix migration version
	migrate -database ${POSTGRESQL_URL} -path internal/migrations -verbose force ${MIGRATE_VERSION}
