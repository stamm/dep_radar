# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [3.0.0] - 2018-05-07
### Added
- Create single page application to select organization on github

### Changed
- Huge refactoring: get rid of src package, move interfaces from one place

## [2.0.0] - 2018-02-27
### Added
- More tests
### Changed
- Huge refactoring: add context to all function working with network
### Fixed
- Linter's warnings
### Removed
- Unused code


## [1.1.4] - 2018-02-23
### Changed
- Fix release docker image


## [1.1.3] - 2018-02-23
### Changed
- Refactor travis.yaml, Makefile to create release


## [1.1.2] - 2018-02-23
### Changed
- Use go 1.10.0 for docker image
- Use alpine 3.7 for docker image
### Added
- Workaround for release in Makefile
- Use upx for minify binary in docker image

## [1.1.1] - 2018-02-23
### Fixed
- Typo


## [1.1.0] - 2018-02-16
### Added
- Binary to analyze github organization
- Docker image


## [1.0.0] - 2018-02-12
### Added
- Support github and private bitbucket
- Support dep and glide


## [0.1.0] - 2017-09-06
### Added
- First release
