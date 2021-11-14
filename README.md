## Regtester

[LOGO here]

**Regtester** is a tool used to stress test docker registries interfaces, it can:

- Run parallel docker pulls against a registry or multiple registries, requesting one or more docker images
- Customize the requests per registry
- Collect metrics and push them to Prometheus push-gateway.


## Usage

1. Get `regtester` from the [download page](http://downloadpage).
2. Configure your performance test file.

*perftest_1.yaml:*
```
connection:
  basicAuth:
    username: user
    password: pass
  CAPath: /tmp/ca-bundle.crt
  extraHeaders:
    key: value
registries:
- registry1
- registry2
images:
- image:1
- image:2
docker:
  pulls: 10000
  parallel: 10

output:
    prometheus:
      pushGatewayUrl: https://mypushgateway:7000/
    csv:
      path: /tmp/perftest_1.csv
```


## Developer Guide

### How to contribute
### Add custom outputs
