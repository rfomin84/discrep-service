# discrep-service

## Migrations

### Create migration
```
docker-compose run detailed-feed-statistic-storage migrate create -ext sql -dir database/clickhouse/migrations -seq create_users_table
```

### Run migration
```
docker-compose run detailed-feed-statistic-storage migrate -database "clickhouse://clickhouse:9000?username=homestead&password=secret&database=homestead&x-multi-statement=true" -path database/clickhouse/migrations up
```

### Rollback migration
```
docker-compose run service-list migrate -database "clickhouse://homestead:secret@tcp(service-list-mysql:3306)/homestead" -path database/migrations down
```

## Запуск cheduler
```
 docker-compose run --rm --service-ports --use-aliases discrepancy-service go run cmd/scheduller/scheduller.go 
```