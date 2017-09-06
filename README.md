<a href="https://travis-ci.org/stamm/dep_radar"><img src="https://travis-ci.org/stamm/dep_radar.svg?branch=master" alt="Build Status"></img></a>
<a href="https://codeclimate.com/github/stamm/dep_radar"><img src="https://codeclimate.com/github/stamm/dep_radar/badges/gpa.svg" alt="Code Climate"></img></a>
<a href="https://codeclimate.com/github/stamm/dep_radar/coverage"><img src="https://codeclimate.com/github/stamm/dep_radar/badges/coverage.svg" /></a>

## Dep radar
`dep radar` is a prototype to control Go dependencies in microservice world.
`dep radar` is not stable yet. It requires Go 1.9 or newer to compile.

## How it works
You can't just run some binary. You have to write a bit of code.
Your code must implement:
* Get a list of packages of applications what dependencies you want to monitor.
* Init provider detector. It can be a default with support only github, but you can add you own provider.
* A http handler with calling method for generate html table with all apps and dependencies.


## Supported code storage
* Github (env: GITHUB_TOKEN)
* Private Bitbucket (env: BB_GIT_URL, BB_GO_GET_URL, BB_USER, BB_PASSWORD)

## Supported dep tools
* [dep](https://github.com/golang/dep)
* [glide](https://github.com/Masterminds/glide)
