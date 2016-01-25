# Botyard 

[![Circle CI](https://circleci.com/gh/botyard/botyard.svg?style=svg)](https://circleci.com/gh/botyard/botyard)

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
