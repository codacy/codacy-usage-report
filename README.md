# codacy-usage-report

## Requirements

To run codacy-usage-report it is required to have GoLang installed.

## Install

You can install the codacy-usage-report using go get command:

    go get -u github.com/codacy/codacy-usage-report

## Usage

To run the codacy usage report, simply run the following command:

    codacy-usage-report

## Configuration file

The configuration file should be placed on the same folder as the executable and should be named `codacy-usage-report.yml`. Example: 

```yaml
accountDB:
  host: localhost
  port: 5432
  database: codacy
  username: username
  password: password
analysisDB:
  host: localhost
  port: 5432
  database: codacy
  username: username
  password: password
# batchSize: 5 - optional
```

## Run from docker

You can also run the codacy usage report using docker:

    docker run -v $PWD/codacy-usage-report.yml:/app/codacy-usage-report.yml -v $PWD/result:/app/result codacy/codacy-usage-report:latest

## What is Codacy?

[Codacy](https://www.codacy.com/) is an Automated Code Review Tool that monitors your technical debt, helps you improve your code quality, teaches best practices to your developers, and helps you save time in Code Reviews.

### Among Codacy’s features:

-   Identify new Static Analysis issues
-   Commit and Pull Request Analysis with GitHub, BitBucket/Stash, GitLab (and also direct git repositories)
-   Auto-comments on Commits and Pull Requests
-   Integrations with Slack, HipChat, Jira, YouTrack
-   Track issues Code Style, Security, Error Proneness, Performance, Unused Code and other categories

Codacy also helps keep track of Code Coverage, Code Duplication, and Code Complexity.

Codacy supports PHP, Python, Ruby, Java, JavaScript, and Scala, among others.

### Free for Open Source

Codacy is free for Open Source projects.
