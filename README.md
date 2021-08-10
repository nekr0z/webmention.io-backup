# webmention.io-backup
a little tool to backup webmentions stored in [webmention.io](https://webmention.io/)

[![Build Status](https://travis-ci.org/nekr0z/webmention.io-backup.svg?branch=master)](https://travis-ci.org/nekr0z/webmention.io-backup) [![codecov](https://codecov.io/gh/nekr0z/webmention.io-backup/branch/master/graph/badge.svg)](https://codecov.io/gh/nekr0z/webmention.io-backup) [![Go Report Card](https://goreportcard.com/badge/github.com/nekr0z/webmention.io-backup)](https://goreportcard.com/report/github.com/nekr0z/webmention.io-backup)

##### Table of Contents
* [How to use](#how)
  * [Options](#command-line-options)
* [Development](#development)
* [Credits](#credits)

#### Help `webmention.io-backup` get better!
Join the [development](#development) (or just [buy me a coffee](https://www.buymeacoffee.com/nekr0z), that helps, too).

## How
Simpy run in command line with the desired options. For regular backups, set up a cron script or a systemd timer.

### Command line options
```
-t [token]
```
the API token for `webmention.io`.

```
-d [domain]
```
only ask for webmentions received for specific domain (i.e. `example.org`); all the domains associated with the account will be processed otherwise.

```
-f [filename]
```
the name of the file to save webmentions to (and read the already backed-up webmentions from). Defaults to `webmentions.json` in current directory.

```
-jf2
```
use the `api/mentions.jf2` endpoint instead of `api/mentions`. The produced JSON will naturally be JF2 in this case.

## Development
Issues reports and pull requests are always welcome!

## Credits
This software includes the following software or parts thereof:
* [The Go Programming Language](https://golang.org) Copyright © 2009 The Go Authors
