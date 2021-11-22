# webmention.io-backup
a little tool to backup webmentions stored in [webmention.io](https://webmention.io/)

[![Build Status](https://github.com/nekr0z/webmention.io-backup/actions/workflows/pre-release.yml/badge.svg)](https://github.com/nekr0z/webmention.io-backup/releases/tag/latest) [![codecov](https://codecov.io/gh/nekr0z/webmention.io-backup/branch/master/graph/badge.svg)](https://codecov.io/gh/nekr0z/webmention.io-backup) [![Go Report Card](https://goreportcard.com/badge/evgenykuznetsov.org/go/webmention.io-backup)](https://goreportcard.com/report/evgenykuznetsov.org/go/webmention.io-backup)

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
use the `/api/mentions.jf2` endpoint instead of `/api/mentions`. The produced JSON will naturally be JF2 in this case.

```
-tlo=false
```
don't create the top-level object in the saved file (i.e. save as an array of webmentions).

```
-p
```
pretty-print (`jq`-style) the saved file.

```
-cd [directory]
```
look in the `directory` for the directory structure that represents the website's structure and try to save webmentions to the individual files (one for each page) in this directory structure; useful for saving webmentions into the source tree of an SSG project.

```
-l [list]
```
list of comma-separated top-level directories to ignore when using `-cd`; with `-cd ./website -l en,fr` both webmentions for `my.site/en/page` and `my.site/fr/page/` will be saved to `./website/page/webmentions.json` instead of separate directories.

```
-ts
```
when using `-cd`, store a timestamp in the root directory and avoid re-fetching webmentions before that timestamp.

## Development
Issues reports and pull requests are always welcome!

## Credits
This software includes the following software or parts thereof:
* [The Go Programming Language](https://golang.org) Copyright © 2009 The Go Authors
