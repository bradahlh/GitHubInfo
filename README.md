# REST API to retrieve GitHub information about a certain project
Individual school project where we created an API to retrieve specific information about a GitHub project. API retrieves owner name, project name, name of top committer, number of commits top committer has made and all languages used in the project.

## Installation
To build:
```
go get github.com/bradahlh/githubinfo
go install github.com/bradahlh/githubinfo
```


## Instructions
1. Start server: `$GOPATH/bin/githubinfo`
2. In your web browser, enter the following basepath: http://localhost:8080/projectinfo/v1/ followed by an organization/username, and then a project name (for example http://localhost:8080/projectinfo/v1/kde/okular for KDE's Okular project)
3. The API should return JSON formatted information about the project in your web browser
