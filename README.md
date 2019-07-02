# BorGoZ


## Intro

__Bor__ g + __Go__ + __Z__ abbix

This is the service written on Go for Borg Backup repo monitoring.


## Enviroment variables.

  - __BORGOZ_HOST__
  - __BORGOZ_PORT__
  - __BORGOZ_REPOS_DIR__ - directory where is all repos
  - __BORGOZ_DEFAULT_REPO_KEY__ - default repo key (it uses when :key is not set)
  - __BORGOZ_LOG_LEVEL__ - set LogLevel
            - "DEBUG"
            - "INFO" (DEFAULT)
            - "WARN"
            - "ERROR"


## Check the last backup not older then :time in seconds.
__/check/lastBackupTime/:repo/:time/:key__ - with specific repo key
__/check/lastBackupTime/:repo/:time__ - with default repo key

  - __:repo__ - borg repo name
  - __:time__ - repo should be yanger then time in minutes
  - __:key__  - borg repo key ( uses BORGOZ_DEFAULT_REPO_KEY if not set)  


# TODO:

## TODO: Check repos health (borg check).
__/check/repoHealth/:repo/:key__
  - __:repo__ - borg repo name
  - __:key__ - borg repo key
