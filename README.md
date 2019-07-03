# BorGoZ


## Intro

__Bor__ g + __Go__ + __Z__ abbix

This is the service written on Go for Borg Backup repo monitoring.

[What is Borg?](https://borgbackup.readthedocs.io/en/stable/)

## Install and run.

Download the binary from releases.

Export variables and start the service:
```bash
$ declare -x BORGOZ_DEFAULT_REPO_KEY="XuolaiziVeiRa8DuPhoopeish5XexingaiveiTaaWaeg8elai1"
$ declare -x BORGOZ_LOG_LEVEL="DEBUG"
$ declare -x BORGOZ_REPOS_DIR="/home/borguser/repos"
$ declare -x BORGOZ_HOST="0.0.0.0"
$ declare -x BORGOZ_PORT="8080"
$ declare -x BORGOZ_BORG_BIN="borg-linux64"

$ borgoz-linux64

```
## Build from source.

```bash
$ git clone https://github.com/vkhodor/borgoz
$ cd ./borgoz
$ make build

```

## Environment variables.

  - __BORGOZ_HOST__
  - __BORGOZ_PORT__
  - __BORGOZ_REPOS_DIR__ - directory where is all Borg Backup repos
  - __BORGOZ_DEFAULT_REPO_KEY__ - default repo key (it uses when :key is not set)
  - __BORGOZ_LOG_LEVEL__ - set LogLevel
            - "DEBUG"
            - "INFO" (DEFAULT)
            - "WARN"
            - "ERROR"

## Checks

### Check the last backup not older then :time in seconds.

__/check/lastBackupTime/:repo/:time/:key__ - with specific repo key
__/check/lastBackupTime/:repo/:time__ - with default repo key

  - __:repo__ - borg repo name
  - __:time__ - count of seconds
  - __:key__  - borg repo key (uses BORGOZ_DEFAULT_REPO_KEY if not set)  

__Example:__

```bash
$ curl http://your.host.local:8080/check/lastBackupTime/some.host.borg.repo-dir/86400    # check last backup not older then 24h (in sec.)
```

The same with specific repokey.

```bash
$ curl http://your.host.local:8080/check/lastBackupTime/some.host.borg.repo-dir/86400/XuolaiziVeiRa8DuPhoopeish5XexingaiveiTaaWaeg8elai1
```

Check returns 200 if OK or 40x if not OK.
