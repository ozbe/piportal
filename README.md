# voggle

Toggle openvpn on same machine via HTTP server.

## Build

```
$ go build -o bin/voggle .
```

## Run

```
$ ./bin/voggle
```

## Usage

```
$ curl http://localhost:8081/
```

### Set active

```
$ curl -X POST -i http://localhost:8081/ --data '{"active":true}'
```

### Set inactive

```
$ curl -X POST -i http://localhost:8081/ --data '{"active":false}'
```
