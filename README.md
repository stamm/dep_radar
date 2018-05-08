[![Travis](https://img.shields.io/travis/stamm/dep_radar.svg?style=flat-square)](https://travis-ci.org/stamm/dep_radar)
[![Code Climate](https://img.shields.io/codeclimate/github/stamm/dep_radar.svg?style=flat-square)](https://codeclimate.com/github/stamm/dep_radar)
[![Code Climate](https://img.shields.io/codeclimate/coverage/github/stamm/dep_radar.svg?style=flat-square)](https://codeclimate.com/github/stamm/dep_radar/coverage)
[![Go Report Card](https://goreportcard.com/badge/github.com/stamm/dep_radar)](https://goreportcard.com/report/github.com/stamm/dep_radar)

## Dep radar
`dep radar` is a prototype to control Go dependencies in microservice world.
`dep radar` requires Go 1.10 or newer to compile.

You can try demo: [https://dep-radar.zagirov.name](https://dep-radar.zagirov.name)

# Screenshots
![Frontend](https://github.com/stamm/dep_radar/raw/master/docs/3.0.png)


## How it works
You can't just run some binary. You have to write a bit of code.
Your code must implement:
* Get a list of packages of applications what dependencies you want to monitor
* Init provider detector. It can be a default with support only github, but you can add you own provider
* A http handler with calling method for generate html table with all apps and dependencies

Simple example that shows a table with dependencies for entered github organization:

`docker run -p 8081:8081 stamm/dep_radar:3.1.2`

To put your token for github use this command:

`docker run -e "GITHUB_TOKEN=t0ken" -p 8081:8081 stamm/dep_radar:3.1.2`


Or with showing state of dependencies:

* Recommended: restriction for version, for example `>=0.13`
* Mandatory: lib must be in an app
* Exclude: lib must be absent in an app
* NeedVersion: you must use version, not hash


```
cat <<EOT > /tmp/recommended.json
{
	"github.com/pkg/errors": {
		"Recommended": ">=0.8.0",
		"Mandatory":  true
	},
	"github.com/kr/fs": {
		"Exclude": true
	},
	"github.com/pkg/profile": {
		"NeedVersion": true
	}
}
EOT
```

`docker run -v /tmp:/cfg -p 8081:8081 stamm/dep_radar:2.0.0 -recommended_file="/cfg/recommended.json" -github_org="dep-radar"`


You can find more in [examples](examples/).



## Supported code storage
* Github
* Private Bitbucket

## Supported dep tools
* [dep](https://github.com/golang/dep)
* [glide](https://github.com/Masterminds/glide)
