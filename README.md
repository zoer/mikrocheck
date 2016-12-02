# MikroTik Firmware checker
[![Build Status](https://travis-ci.org/zoer/mikrocheck.svg)](https://travis-ci.org/zoer/mikrocheck)

## Install

```
$ go get github.com/zoer/mikrocheck
```
## Usage

```
$ mikrocheck -h
Usage of ./mikrocheck:
  -addr string
        SMTP server address (default "smtp.yandex.ru:25")
  -from string
        Mail 'From:' header
  -password string
        SMTP Auth password
  -store string
        File to store current version (default "$HOME/.mikrocheck")
  -to string
        Mail 'To:' header
  -username string
        SMTP Auth username (if isn't set it takes the 'from' flag value.)


$ mikrocheck -to myemail1@foo.com,myemai2@foo.com -from foo@baz.com -password mypass
```

### Crone
```
0 12 * * * mikrocheck -to myemail1@foo.com,myemai2@foo.com -from foo@baz.com -password mypass >> /var/log/mikrocheck.log 2>&1
```
