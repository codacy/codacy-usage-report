# codacy-usage-report

## Install

You can install the codacy-usage-report using go get command:

    go get -u github.com/codacy/codacy-usage-report

## Build docker image

Build a docker image for the usage report script using:

    make dockerbuild

## Configuration file

The configuration file should be placed on the same folder as the executable and should be named `codacy-usage-report.yml`. Example: 

```yaml
accountDB:
  host: localhost
  port: 5432
  database: accountDB
  username: username
  password: password
analysisDB:
  host: localhost
  port: 5432
  database: analysisDB
  username: username
  password: password
# batchSize: 5 - optional
```

## Run from docker

    docker run -v $PWD/codacy-usage-report.yml:/app/codacy-usage-report.yml -v $PWD/result:/app/result codacy-usage-report:latest
