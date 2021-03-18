## Example

Launch Redis and run `main.go`

```bash
$ docker-compose up -d
```

Check response header

```bash
$ curl -sI -X GET 127.0.0.1:8080/ping | grep -iE 'X-Ratelimit|HTTP/1.1'
HTTP/1.1 200 OK
X-Ratelimit-Remaining: 975
X-Ratelimit-Reset: 3372
```

