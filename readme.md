### Criando migration
`migrate create -ext=sql -dir=sql/migrations -seq init`

### Executando migration
`migrate --path=sql/migrations -database "mysql://root:root@tcp(localhost:3306)/courses" -verbose up`