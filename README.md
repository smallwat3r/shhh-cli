# Shhh-cli  

## What is it?  

**shhh-cli** is a Command Line Interface tool to interact with the 
[Shhh](https://github.com/smallwat3r) API. This tool allows you to 
create and read secrets directly from the command line / terminal.  

![shhh-cli](https://i.imgur.com/zGF2015.gif)  

## How to install it?  

If you are a Go user:
```sh
go get -u github.com/smallwat3r/shhh-cli \
    && mv $GOPATH/bin/shhh-cli $GOPATH/bin/shhh
```

Or   
shhh-cli has no runtime dependencies, you can download a binary for 
your platform [here](https://github.com/smallwat3r/shhh-cli/releases).

## Using your own instance of Shhh?  

shhh-cli interacts by default with the official Shhh server. If 
you've set-up your own Shhh server, and want to create secrets from
this server by default, you can set-up an environment variable 
`SHHH_SERVER`.

```sh
# Example (in you bashrc)
export SHHH_SERVER=https://<my-custom-shhh-server>.com
```

Note: you will still be able to read secrets from other Shhh servers.

## How to use it?  

```console
$ shhh -h
Create or read secrets from a Shhh server.

Usage:
  shhh <command> [<args>]

Options:
  -h         Show help message.

Modes:
  create     Creates a secret message.
  read       Read a secret message.

Usage of create:
  -h         Show help message.
  -m string  Secret message to encrypt.
  -p string  Passphrase to encrypt secret.
  -d int     Optional, number of days to keep the secret alive. (default 3).

Usage of read:
  -h         Show help message.
  -l string  URL link to access secret.
  -p string  Passphrase to decrypt secret.

Examples:
  shhh create -m 'this is a secret msg.' -p 'P!dhuie0e3bdiu' -d 2
  shhh read -l https://<shhh-server>/api/r/jKD8Uy0A9_51c8asqAYL -p 'P!dhuie0e3bdiu'
```
