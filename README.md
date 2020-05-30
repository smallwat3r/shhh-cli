<h3 align="center">shhh-cli</h3>
<p align="center">Go CLI client for Shhh</p>

---

**shhh-cli** is a Command Line Interface tool interacting with the 
[Shhh](https://github.com/smallwat3r) API.  
This allows you to create and read secrets directly from the 
terminal.

![shhh-cli](https://i.imgur.com/zGF2015.gif)  

## How to install it?  

If you are a Go user:
```sh
go get -u github.com/smallwat3r/shhh-cli \
    && mv $GOPATH/bin/shhh-cli $GOPATH/bin/shhh
```

Or, shhh-cli has no runtime dependencies, you can download a binary for 
your platform [here](https://github.com/smallwat3r/shhh-cli/releases).

## Using your own instance of Shhh?  

shhh-cli interacts by default with the official Shhh server.  

If you've set-up your own Shhh server, and want to create secrets 
from this server by default, you will need to set-up an `SHHH_SERVER`
environment variable.

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
  -s string  Optional, Shhh target server (ex: https://shhh-encrypt.herokuapp.com).

Usage of read:
  -h         Show help message.
  -l string  URL link to access secret.
  -p string  Passphrase to decrypt secret.

Examples:
  shhh create -m 'this is a secret msg.' -p P!dhuie0e3bdiu -d 2
  shhh read -l https://shhh-encrypt.herokuapp.com/r/jKD8Uy0A9_51c8asqAYL -p P!dhuie0e3bdiu
```
