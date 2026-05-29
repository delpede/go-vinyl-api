# go-vinyl-api
An API to keep track of your vinyl collection

## Migration
Run the migration like this
```
migrate \
  -path migrations \
  -database "postgres://vinyl:vinyl@localhost:5432/vinyl?sslmode=disable" \
  up
```
