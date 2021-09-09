# jsonhttp

A http server for testing.

# Usage

```
% jsonhttp -h
Usage of dist/jsonhttp:
  -c string
        config file (default "server.json")
  -debug
        enable debug log
  -p int
        port number
```

The format of `config file` is `message Server` in `pb/origin.proto`.

## Hello

```
{
  "handlers": [
    {
      "path": "/hello",
      "methodType": "GET",
      "action": {
        "return": {
          "status": 200,
          "templates": [
            {
              "type": "BODY",
              "value": {
                "m": {
                  "values": {
                    "message": {
                      "s": "Hello!"
                    }
                  }
                }
              }
            }
          ]
        }
      }
    }
  ]
}
```

## Echo with random delay

```
{
  "handlers": [
    {
      "path": "/echo",
      "methodType": "GET",
      "action": {
        "return": {
          "status": 200,
          "delay": {
            "util": {
              "random": {
                "dice": {
                  "min": 200,
                  "max": 1000
                }
              }
            }
          }
        }
      }
    }
  ]
}
```

## Proxy

```
{
  "handlers": [
    {
      "path": "/gw",
      "methodType": "GET",
      "action": {
        "gateway": {
          "path": {
            "s": "http://127.0.0.1:10000/hello"
          },
          "methodType": "GET",
          "timeout": {
            "n": 300
          }
        }
      }
    }
  ]
}
```

## Dynamic proxy

```
{
  "handlers": [
    {
      "path": "/dgw",
      "methodType": "POST",
      "action": {
        "gateway": {
          "path": {
            "add": {
              "values": [
                {
                  "s": "http://127.0.0.1:"
                },
                {
                  "body": {
                    "keys": ["port"]
                  }
                },
                {
                  "s": "/hello"
                }
              ]
            }
          },
          "methodType": "GET"
        }
      }
    }
  ]
}
```

# Build

```
make build
```
