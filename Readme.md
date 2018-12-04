## Clone Project

```
git clone <this repo>
```

## Run on Local Host

```
cd go-web-redis-docker
run app.go
```

Access the app using curl

```
curl localhost:8080
```

Start Redis instance on localhost

access http://localhost:8080/redis on browser

## Create docker Image

Two phaze docker build

```sh
docker build -t go-app .
```
