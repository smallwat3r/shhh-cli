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
$ shhh --help
Create or read secrets from a Shhh server.

Find more information at https://github.com/smallwat3r/shhh-cli/blob/master/README.md

Usage:
  shhh [mode] [options]

Options:
  -h, --help   Show help message and exit.

Modes:
  create       Creates a secret message.
  read         Read a secret message.

Usage of create:
  -h, --help                 Show create help message and exit.
  -m, --message    <string>  Secret message to encrypt.
  -p, --passphrase <string>  Passphrase to encrypt secret.
  -d, --days       <int>     (optional) Number of days to keep the secret alive (default: 3).
  -s, --server     <string>  (optional) Shhh target server (ex: https://shhh-encrypt.herokuapp.com).
  example: shhh create -m 'a secret msg' -p P!dhuie0e3bdiu -d 2

Usage of read:
  -h, --help                 Show read help message and exit.
  -l, --link       <string>  URL link to access secret.
  -p, --passphrase <string>  Passphrase to decrypt secret.
  example: shhh read -l https://shhh-encrypt.herokuapp.com/r/jKD8Uy0A9_51c8asqAYL -p P!dhuie0e3bdiu
```
