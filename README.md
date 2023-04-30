# Requirements

Make sure you have installed Docker, Docker Compose

# Usage

## Run server

```
$ docker compose up
```

## Connect to server

```
$ nc -u localhost 2000
```

## Usage

```
Usage: get_result format
All supported formats: avro, protobuf, json, msgpack, yaml
```

## Example
```
$ nc -u localhost 2000
get_result protobuf
avro-476-17280us-1608us
```


