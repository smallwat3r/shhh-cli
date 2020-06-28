<h3 align="center">shhh-cli</h3>
<p align="center">Go CLI client for Shhh</p>

---

**shhh-cli** is a Command Line Interface tool interacting with the 
[Shhh](https://github.com/smallwat3r) API.  
This allows you to create and read secrets directly from the 
terminal.

![shhh-cli](https://i.imgur.com/qLx3aoV.png)  

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

Shhh-cli has no runtime dependencies, you can download a binary for 
your platform [here](https://github.com/smallwat3r/shhh-cli/releases), 
then rename it `shhh` and place it in your `bin` directory.

## Using a self-hosted Shhh instance?  

shhh-cli interacts by default with the official Shhh server. If 
you've set-up your own Shhh server, and want to create secrets 
from this server by default, you will need to set-up an `SHHH_SERVER`
environment variable.

```sh
# Example (in you bashrc)
export SHHH_SERVER=https://<my-custom-shhh-server>.com
```

Note: this won't impact reading secrets from other Shhh servers, and
you will still be able to create secrets in other servers using the 
`--server` option.

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
  -t, --tries      <int>     (optional) Number of tries to open secret before it gets deleted (default: 5).
  -h, --host       <string>  (optional) Shhh target server (ex: https://shhh-encrypt.herokuapp.com).
  -s, --secure               (optional) If flag set, check passphrase against the haveibeenpwned API.
  example: shhh create --message 'a secret msg' --passphrase PdVUe3bdiu --days 2 --secure

Usage of read:
  -h, --help                 Show read help message and exit.
  -l, --link       <string>  URL link to access secret.
  -p, --passphrase <string>  Passphrase to decrypt secret.
  example: shhh read --link https://shhh-encrypt.herokuapp.com/r/jKD8Uy0A9_51c8asqAYL --passphrase PdVUe3bdiu
```
