# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Pre-release]
### Added
* language separation option (`-lang`)

### Fixed
* Go version

## [1.4.0] - 2021-11-15
### Changed
* bump Go to 1.16
* upgrade dependencies

### Added
* option to pretty-print the output file(s)

## [1.3.0] - 2021-08-13
### Added
* option to save webmentions for separate pages to separate files in the directory structure that represents the site structure
* option to use timestamp to avoid re-downloading old webmentions in SSG mode

## [1.2.0] - 2021-08-10
### Added
* option to only query webmentions for a specific domain
* option to have the top-level object in the output file

### Fixed
* default behaviour to use the top-level object (broken in 1.1.0)
* all mentions were treated as new when using JF2

## [1.1.0] - 2021-08-09
### Changed
* output file no longer contains the redundant top-level object

### Added
* option to use JF2 endpoint

## [1.0.2] - 2021-08-09
### Fixed
* data loss due to too many assumptions about what's stored for each webmention
* windows builds packaged as .tar.gz

## [1.0.1] - 2020-05-15
### Fixed
* build and packaging process

## [1.0.0] - 2020-05-15
*initial release*

[Pre-release]: https://github.com/nekr0z/webmention.io-backup/releases/tag/latest
[1.4.0]: https://github.com/nekr0z/webmention.io-backup/releases/tag/v1.4.0
[1.3.0]: https://github.com/nekr0z/webmention.io-backup/releases/tag/v1.3.0
[1.2.0]: https://github.com/nekr0z/webmention.io-backup/releases/tag/v1.2.0
[1.1.0]: https://github.com/nekr0z/webmention.io-backup/releases/tag/v1.1.0
[1.0.2]: https://github.com/nekr0z/webmention.io-backup/releases/tag/v1.0.2
[1.0.1]: https://github.com/nekr0z/webmention.io-backup/releases/tag/v1.0.1
[1.0.0]: https://github.com/nekr0z/webmention.io-backup/releases/tag/v1.0.0
