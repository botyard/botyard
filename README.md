# Botyard 

[![Circle CI](https://circleci.com/gh/botyard/botyard.svg?style=svg)](https://circleci.com/gh/botyard/botyard) [![Go Report](https://goreportcard.com/badge/github.com/botyard/botyard)](https://goreportcard.com/report/github.com/botyard/botyard) [![codebeat badge](https://codebeat.co/badges/6a35eba5-e617-4d80-8fd4-2920dffcd1f1)](https://codebeat.co/projects/github-com-botyard-botyard) [![Docker](https://img.shields.io/docker/pulls/botyard/botyard.svg)](https://hub.docker.com/r/botyard/botyard/)

## Configuration

```
name: boty
commands:
    - name: "hello"
      command: "hello <world>"
      url: "http://example.com/hello?world=<world>"
      type: http
      method: GET
irc:
    - name: irc1
      server: irc.org 
      channels:
        - "#channel1"
slack:
    - name: slack1
      token: TOKEN
```

I think that the config file explains all about what the `botyard` is. :-) 

## State

"WORK IN PROGRESS"
