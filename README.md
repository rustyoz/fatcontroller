# FatController
[![Build Status](https://travis-ci.org/rustyoz/fatcontroller.svg?branch=master)](https://travis-ci.org/rustyoz/fatcontroller)
## What is it?

A virtual hosting tool written in Golang.

## How do I use it?
Place a toml configuration file in either the working directory or the a .fatcontroller directory in the users home directory.
The default configuration filename is *config.toml*

## Example configuration.
```
port = 80
host = "localhost"

[[virtualhost]]
name = "web"
url = "/"
host = "localhost"
path = "/"
port = 8000

[[virtualhost]]
name = "service"
url = "/service/"
path = "/"
host = "localhost"
port = 9000
```

The top level configuration of *port* and *host* define the address on which FatController will listen on.

### For each virtual host
Define the start of a new virtual host configuration
```
[[virtualhost]]
```
Give it a name
```
name = "service"
```
URL by which it will be accessed.
```
url = "/service/"
```
Path by which the endpoint will except requests
```
path = "/"
```
Hostname and port of the endpoint
```
host = "localhost"
port = 9000
```

## Made using
[Viper](https://github.com/spf13/viper) by Steve Francia [spf13](https://github.com/spf13)
