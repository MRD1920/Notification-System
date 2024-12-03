DB_URL=postgres://postgres:mrd123@localhost:5431/notification_system?sslmode=disable
migrateup:
	migrate -path db/migrations -database "$(DB_URL)" -verbose up
migrateup1:
	migrate -path db/migrations -database "$(DB_URL)" -verbose up 1
migratedown:
	migrate -path db/migrations -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migrations -database "$(DB_URL)" -verbose down 1
sqlc:
	sqlc generate
server:
	go run main.go


.PHONY: migrateup migrateup1 migratedown migratedown1 sqlc server