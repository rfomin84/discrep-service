# discrep-service

## Migrations

### Create migration
```
docker-compose run discrepancy-service migrate create -ext sql -dir database/clickhouse/migrations -seq create_users_table
```

### Run migration
```
docker-compose run discrepancy-service migrate -database "clickhouse://clickhouse:9000?username=homestead&password=secret&database=homestead&x-multi-statement=true" -path database/clickhouse/migrations up
```

### Rollback migration
```
docker-compose run discrepancy-service migrate -database "clickhouse://clickhouse:9000?username=homestead&password=secret&database=homestead&x-multi-statement=true" -path database/clickhouse/migrations down
```

## Запуск cheduler
```
 docker-compose run --rm --service-ports --use-aliases discrepancy-service go run cmd/scheduller/scheduller.go 
```

## Запуск server
```
 docker-compose run --rm  --use-aliases -p 8080:8080 discrepancy-service go run cmd/server/server.go 
```