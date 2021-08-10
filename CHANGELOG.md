# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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

[Unreleased]: https://github.com/nekr0z/webmention.io-backup/compare/v1.2.0...HEAD
[1.2.0]: https://github.com/nekr0z/webmention.io-backup/releases/tag/v1.2.0
[1.1.0]: https://github.com/nekr0z/webmention.io-backup/releases/tag/v1.1.0
[1.0.2]: https://github.com/nekr0z/webmention.io-backup/releases/tag/v1.0.2
[1.0.1]: https://github.com/nekr0z/webmention.io-backup/releases/tag/v1.0.1
[1.0.0]: https://github.com/nekr0z/webmention.io-backup/releases/tag/v1.0.0
