# Go-jenkins-exporter

go-jenkins-exporter is an exporter for Prometheus, which allows you to monitor/alert on multiple jenkins job statuses and properties.

## Installation

### Using source code

```shell
go get -u -v github.com/goodbins/go-jenkins-exporter
cd $GOPATH/src/github.com/goodbins/go-jenkins-exporter
make deps
make install
```

You can create your local Docker image using:

```shell
make image
```

### Using Docker

```shell
docker pull goodbins/go-jenkins-exporter:latest
docker run -it -p 5000:5000 goodbins/go-jenkins-exporter:latest
```

## Usage

Say you have a jenkins instance at http://jenkins-ci:8080

First thing to do is to export some env vars:

```shell
export JENKINS_USERNAME=yourusername
export JENKINS_PASSWORD=yourpassword
```

Note: You can also use a token instead of a password.

Then you can launch the exporter using the following command:

```shell
./go-jenkins-exporter -j jenkins-ci:8080 -r 2s
```

By default, go-jenkins-exporter listens at [localhost:5000](localhost:5000)

Using the public registry Docker image:

```shell
docker run -it \
    -p 5000:5000 \
    -e JENKINS_USERNAME=yourusername \
    -e JENKINS_PASSWORD=yourpassword \
    --restart=unless-stopped \
    goodbins/go-jenkins-exporter:latest -j jenkins-ci:8080 -r 2s
```

For more configuration options you can use:

```shell
./go-jenkins-exporter --help
```

This gives something like:

```console
Usage:
  go-jenkins-exporter [flags]

Flags:
  -h, --help               help for go-jenkins-exporter
  -j, --jenkins string     Jenkins API host:port pair
  -l, --listen string      Exporter host:port pair (default "localhost:5000")
  -m, --metrics string     Path under which to expose metrics (default "/metrics")
  -a, --path string        Jenkins API path (default "/api/json")
  -r, --rate duration      Set metrics update rate in seconds (default 1s)
  -s, --ssl                Enable TLS (default false)
  -t, --timeout duration   Jenkins API timeout in seconds (default 10s)
  -v, --verbose            Enable verbosity
      --version            version for go-jenkins-exporter
```

## Prometheus configuration

You can add the endpoint to your prometheus.yml file:

```yaml
scrape_configs:
  - job_name: 'go-jenkins-exporter'
    scrape_interval: 5s
    static_configs:
      - targets: ['localhost:5000']
```

## Licence
Unless otherwise noted, the go-jenkins-exporter source files are distributed under the MIT license found in the LICENSE file.

## Next steps...

 - Using bndr/gojenkins to interact with Jenkins API
 - Expose the metrics of the slave nodes
 - Create a helm chart to deploy the exporter on k8s
 - write unit tests
 
## Contribute
Go to [contributing.md](CONTRIBUTING.md)
