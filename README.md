# your_project

Just a bloilerplate for a GoLang service with Fiber and SQLC.

## Features

* Postgres Database Integration With PGX
* Raw automated-generation queries with SQLC
* JWT generation and validation middleware with HS512 alg
* Ready user resource with login, fetch user data, update and soft delete
* Error logger integrated with Fiber and easy to format
* Direct SSL support for production mode

Please help me to support this boilerplate if you can.

## Setup

### Migrate

```shell
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### SQLc

```shell
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

## Build

### From Windows
#### To Windows
```shell
set CGO_ENABLED=0 & set GOOS=windows& go build -o ./build/server.exe
```
#### To Linux
```shell
set CGO_ENABLED=0 & set GOOS=linux& go build -o ./build/server
```
#### To MacOS
```shell
set CGO_ENABLED=0 & set GOOS=darwin& go build -o ./build/server
```

### From Linux
#### To Windows
```shell
CGO_ENABLED=0 GOOS=windows go build -o ./build/server.exe
```
#### To Linux
```shell
CGO_ENABLED=0 GOOS=linux go build -o ./build/server
```
#### To MacOS
```shell
CGO_ENABLED=0 GOOS=darwin go build -o ./build/server
```

For MacOS, it will be probably the same that Linux.

## Migrations

### Create a new migration

```shell
migrate create -ext sql -dir sql/migration -seq file_name
```

### Grant permissions in database

```sql
grant all on database your_database to your_user;
grant usage on schema public TO your_user;
```

### Run the migrations

Up:
```shell
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/your_database?sslmode=disable"
make migrate-up
```

Down:
```shell
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/your_database?sslmode=disable"
make migrate-down
```

## Running

### Install service on Linux

```shell
sudo cp ./boilerplate-go-fiber-sqlc.service /etc/systemd/system/
sudo systemctl enable boilerplate-go-fiber-sqlc.service
sudo systemctl restart boilerplate-go-fiber-sqlc.service
```

### Environment variables

```shell
APP_ENV=dev
APP_HEADER=Your Project Server Header
APP_NAME=Your Project Name
APP_HOST=0.0.0.0
APP_PORT=8080
DATABASE_URL=postgresql://{user}:{password}@{host}:{port}/{dbname}?sslmode=disable
JWT_ISSUER=your_project
JWT_SECRET=your_secret
```

Note: you need to enable sslmode if you are in production mode!
