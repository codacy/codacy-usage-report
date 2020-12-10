# codacy-usage-report

codacy-usage-report is a script that generates a CSV file with the following activity information for each Codacy user:

-   Date of creation
-   Date of the last login
-   Date of the last commit
-   Number of commits
-   Email addresses used on the commits
-   Date of removal, if applicable

The script obtains the information directly from your Codacy databases and can only be used on Codacy Self-hosted instances.

## Requirements

To run codacy-usage-report you must have:

-   Read access to the `account` and `analysis` Codacy databases from the environment where you will run codacy-usage-report
-   Either [Go](https://golang.org/dl/) or [Docker](https://www.docker.com/) installed

## Configuration

Create a configuration file `codacy-usage-report.yml` with the example syntax below.

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
# You may need to lower the batch size from the default value
# if you experience timeouts when running the script:
# batchSize: 10000000
```

You must have this configuration file in your current working directory or specify it with the flag `--configFile` when running codacy-usage-report.

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

    **Note:** Make sure that you have [included the Go bin folder in your `PATH` environment variable](https://golang.org/doc/install#install).

### Running codacy-usage-report using Docker

Alternatively, you can run codacy-usage-report using Docker:

```bash
docker run -v $PWD/codacy-usage-report.yml:/app/codacy-usage-report.yml \
           -v $PWD/codacy-usage-report:/app/codacy-usage-report \
           codacy/codacy-usage-report:latest
```

### Command-line options

```bash
codacy-usage-report [--configFile <configuration file path>]
                    [--outputFolder <output folder path>]
                    [--help]
```

-   `--configFile`

    Path of the `codacy-usage-report.yml` configuration file. The default is `./codacy-usage-report.yml`.

-   `--outputFolder`

    Path of the output folder to store the CSV file. The default is `./codacy-usage-report/`.

-   `--help`

    Print usage information.

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
