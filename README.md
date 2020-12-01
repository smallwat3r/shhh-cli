<h3 align="center">shhh-cli</h3>
<p align="center">Go CLI client for Shhh</p>

**shhh-cli** is a Command Line Interface tool interacting with 
[Shhh](https://github.com/smallwat3r/shhh).  
This util allows you to create and read secrets using Shhh directly from the 
terminal.  

![shhh-cli](https://i.imgur.com/DAf4qRg.png)  

## How to install it?  

#### If you're a GO user  
```sh
go get -u github.com/smallwat3r/shhh-cli \
    && mv $GOPATH/bin/shhh-cli $GOPATH/bin/shhh
```

#### Using Homebrew  

```sh
brew tap smallwat3r/scripts \
  && brew install shhh
```

#### Manually  

This util has no runtime dependencies, you can download a binary for 
your platform [here](https://github.com/smallwat3r/shhh-cli/releases).

## Using a self-hosted Shhh instance?  

This util interacts by default with the official Shhh server. If 
you've set-up your own Shhh server, and want to create secrets 
from this server by default, you will need to set-up a custom `SHHH_SERVER`
environment variable.

```sh
# Example (in you bashrc)
export SHHH_SERVER=https://<my-custom-shhh-server>.com
```

_Note: this won't impact reading secrets from external Shhh servers, and
you will still be able to create secrets using other Shhh servers using the 
`--host` option._  

## How to use it?  

```console
$ shhh --help
Create or read secrets from a Shhh server.

Find more information at https://github.com/smallwat3r/shhh-cli/blob/master/README.md

Usage:
  shhh [mode] [options]

Options:
  -h, --help     Show help message and exit.
  -v, --version  Show program version and exit.

Modes:
  create         Creates a secret message.
  read           Read a secret message.

Usage of create:
  -h, --help                 Show create help message and exit.
  -m, --message    <string>  Secret message to encrypt.
  -p, --passphrase <string>  Passphrase to encrypt secret.
  -d, --days       <int>     (opt) Nb of days to keep the secret alive (default: 3).
  -t, --tries      <int>     (opt) Max nb of tries to open the secret (default: 5).
  -h, --host       <string>  (opt) Shhh target server (ex: https://<server>.com).
  -s, --secure               (opt) Check passphrase against the haveibeenpwned API.

  example: shhh create -m [secret] -p [passphrase] -d 2 -t 3 -s

Usage of read:
  -h, --help                 Show read help message and exit.
  -l, --link       <string>  URL link to access secret.
  -p, --passphrase <string>  Passphrase to decrypt secret.

  example: shhh read -l [link] -p [passphrase]
```
