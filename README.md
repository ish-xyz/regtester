[LOGO here]

**Regtester** is created to run load tests against one or more private registries.

Regetester will:
1. Run parallel Docker pulls against a registry or multiple registries, requesting one or more docker images
2. Collect metrics and output them as csv


The tool is shipped as CLI and can be extended to support third party systems, such as:
- Send metrics to prometheus (via PushGateway)
- customize requests per registry


## Usage

1. Get `regtester` from the [download page](http://downloadpage).
2. Configure your performance test file.

*perftest_1.yaml:*
```
connection:
  basicAuth:
    username: ""
    password: ""
  CAPath: /tmp/ca-bundle.crt
  extraHeaders:
    key: value
registries:
- registry1
- registry2
images:
- image:1
- image:2
workload:
  pulls: 10000
  maxConcPulls : 10
  maxConcLayers: 10
  checkIntegrity: true
```
