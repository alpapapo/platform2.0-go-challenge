# GlobalWebIndex Engineering Challenge

# Project Title

GWI Platform2.0 go challenge (alpapapo version)

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

Get a copy of the project
```
go get github/alpapapo/platform2.0-go-challenge
```
Go to the project folder:
```
cd ~/go/src/github/alpapapo/platform2.0-go-challenge
```
Install prerequisites libraries
```
go mod download
```

### Installing Development

Get the development env running 

Make the .env.development.local from .env.developmet
```
cp .env.developmet .env.developmet.local
```
Edit the following variables in .env.development.local:
Mysql settings
```
DATABASE_HOST=localhost
DATABASE_PORT=3306
DATABASE_USER=gouser
DATABASE_PASSWORD=goPathword
```
Enable or disable orm debug verbosity
```
DATABASE_DEBUG=enabled
```
Enable or disable jwt authentication (enabled|disabled)
```
JWT_AUTH=disabled
```

### Run 
#####(Non Docker)
Install mysql - follow the instructions:

* [installing-mysql-server-on-ubuntu](https://support.rackspace.com/how-to/installing-mysql-server-on-ubuntu/)

Ensure that no other process running on 8000 (sudo lsof -i :8000) and run the goassets application
```
go run main.go
```
The first time the main function runs, the database is created and the migrations are applied (gorm style)

#####(Docker)
Ensure no other application running on 3306 and 8000
```
docker-compose up
```

### Application Playground Steps
#####Register user
```
curl -X POST \
  http://localhost:8000/api/user/register \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
	"email": "alpapapo@gmail.com",
	"username": "alpapapo",
	"password": "%(123456)"
}'
```
#####Login user
```
curl -X POST http://localhost:8000/api/user/login \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
	"email": "alpapapo@gmail.com",
	"password": "%(123456)"
}'
```
INFO: Copy the jwt token from Login Response and replace wherever <jwt> with it in the following curls
or disable JWT_AUTH and remove  "-H 'Authorization: Bearer <jwt>' \" from curls

#####Populate database 

with dummy assets from sample/data.json (This endpoint is only for demonstrating purposes)
```
curl -X POST \
  http://localhost:8000/populate/assets \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjEsImV4cCI6MTU3MTA1MTE0MX0.RYP19dT_RXRheUHdBuHPnnkwqpx9y2PQ0seKRlD-TIg' \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{	
}'
```
#####Get all Assets
```
curl -X GET \
  http://localhost:8000/api/assets \
  -H 'Authorization: Bearer <jwt>' \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache'
```
#####Mark 4 Assets as Favorites
```
curl -X POST \
  http://localhost:8000/api/assets/favorites \
  -H 'Authorization: Bearer <jwt>>' \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '[1,2,4,5]'
```
#####UnMark 1 and 4 From Favorites
```
curl -X DELETE \
  http://localhost:8000/api/assets/favorites \
  -H 'Authorization: Bearer <jwt>' \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '[1,4]'
```
#####User deletes softly 2
```
curl -X DELETE \
  http://localhost:8000/api/assets/2 \
  -H 'Authorization: Bearer <jwt>' \
  -H 'cache-control: no-cache'
``` 
#####Update description of Asset 5
```
curl -X PUT \
  http://localhost:8000/api/assets/5 \
  -H 'Authorization: Bearer <jwt>' \
  -H 'Content-Type: application/json' \
  -H 'cache-control: no-cache' \
  -d '{
	"desc": "other_description"
}'
```
#####Get Asset by ID
```
curl -X GET \
  http://localhost:8000/api/assets/5 \
  -H 'Authorization: Bearer <jwt>' \
  -H 'cache-control: no-cache'
```
### Running the tests
#####(Non Docker version)
edit .env.test.local according to your mysql settings
```
go test -v
```
#####(Docker-Compose version)
Ensure no other application running on 3306
```
docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit
```

## Deployment

TODO: Add additional notes about how to deploy this on a live system

## Built With

* [dgrijalva/jwt-go](https://github.com/dgrijalva/jwt-go) - The JWT implementation used
* [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql) - My SQL Driver
* [gorilla/mux](https://github.com/gorilla/mux) - Gorilla HTTP Router and Dispatcher
* [GORM](https://github.com/jinzhu/gorm) - ORM library
* [joho/godotenv](https://github.com/joho/godotenv) - The dotenv implementation used

