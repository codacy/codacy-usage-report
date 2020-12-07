# codacy-usage-report

codacy-usage-report is a script that generates a CSV file with the following information about the activity of your Codacy users:

-   Last time a user logged in
-   Date of the last commit
-   Number of commits

The script obtains the information directly from your Codacy database and can only be used on Codacy Self-hosted instances.

## Requirements

To run codacy-usage-report you must have either [Go](https://golang.org/dl/) or [Docker](https://www.docker.com/) installed.

## Configuration

Create a configuration file `codacy-usage-report.yml` with the example syntax below. You must have this configuration file in your current working directory when running codacy-usage-report.

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
# batchSize: 5 (optional)
```

## Usage

### Running codacy-usage-report using Go

To run codacy-usage-report directly using Go:

1.  Install codacy-usage-report using Go:

    ```bash
    go get -u github.com/codacy/codacy-usage-report
    ```

2.  Run codacy-usage-report:

    ```bash
    codacy-usage-report
    ```

    **Note:** Make sure that you have [included the Go bin folder in your PATH environment variable](https://golang.org/doc/install#install).

### Running codacy-usage-report using Docker

Alternatively, you can run codacy-usage-report using Docker:

```bash
docker run -v $PWD/codacy-usage-report.yml:/app/codacy-usage-report.yml \
           -v $PWD/result:/app/result \
           codacy/codacy-usage-report:latest
```

## What is Codacy?

[Codacy](https://www.codacy.com/) is an Automated Code Review Tool that monitors your technical debt, helps you improve your code quality, teaches best practices to your developers, and helps you save time in Code Reviews.

### Among Codacy's features:

-   Identify new Static Analysis issues
-   Commit and Pull Request Analysis with GitHub, BitBucket/Stash, GitLab (and also direct git repositories)
-   Auto-comments on Commits and Pull Requests
-   Integrations with Slack, HipChat, Jira, YouTrack
-   Track issues Code Style, Security, Error Proneness, Performance, Unused Code and other categories

Codacy also helps keep track of Code Coverage, Code Duplication, and Code Complexity.

Codacy supports PHP, Python, Ruby, Java, JavaScript, and Scala, among others.

### Free for Open Source

Codacy is free for Open Source projects.
