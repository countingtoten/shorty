# Shorty
A URL shortener written in Go

# Startup

Start shorty by using docker-compose or by building a docker image. Shorty will bind to port 3000 by default.

## Docker Compose

```sh
docker-compose up
```

## Creating and running a docker container

```sh
docker build -t shorty:latest .
docker run --name shorty_web_1 -p 3000:3000 -it shorty:latest

# After shutting down the docker container use the following command to start it
docker start -i shorty_web_1
```

# Routes

## POST localhost:3000/new

Create a post request to localhost:3000/new to creates a new short URL associated with the specified user. If the user id does not exist, it is added.

### Request

```json
{
  "user_id": 1,
  "url": "https://example.com"
}
```

### Response

```json
{
  "short_url": "http://localhost:3000/abcdef"
}
```

## GET localhost:3000/abcdef

Accessing the returned short URL redirects to its long URL.
