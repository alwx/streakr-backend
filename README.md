# Development

1. Install Docker: https://docs.docker.com/docker-for-mac/install/ 
2. Install [go](https://github.com/golang/go) and [dep](https://github.com/golang/dep): `brew install go dep`
3. `dep ensure` to download dependencies
5. `go get github.com/pilu/fresh` to download [fresh](https://github.com/pilu/fresh)
4. `make run-development-deps` to run Postgres and Redis
6. `fresh` to start (re-)building the app!

# Production

* `docker run --name db -e POSTGRES_PASSWORD=spaces-password -e POSTGRES_DB=spaces-db -d postgres:latest`
* `docker run --name redis -d redis:alpine`
* `docker run --name spaces -p 3000:3000 -e GIN_MODE=release -e SPACES_ENVIRONMENT=production -d alwxx/streakr-go:latest`