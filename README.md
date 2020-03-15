# Shhh-cli  

**shhh-cli** is a Command Line Interface tool to interact with [Shhh](https://github.com/smallwat3r) web-application API.  
This tool allows you to create and read secrets directly from the command line / terminal.  

![shhh-cli](https://i.imgur.com/HntOMrf.gif)  

## Tell shhh-cli to talk to your own Shhh server  

shhh-cli interacts by default with the official Shhh server (shhh-encrypt.com).  
If you host Shhh on your own server, you can set-up an env variable `SHHH_SERVER`.  
It will know interact with your server as default.  
```sh
# Example (in you bashrc or zshrc)
export SHHH_SERVER=https://my_own_shhh_app_server.com
```

## Install  

If you are a Go user:
```sh
go get -u github.com/smallwat3r/shhh-cli && \
    mv $GOPATH/bin/shhh-cli $GOPATH/bin/shhh
```

Other: shhh-cli has no runtime dependencies. Download a binary for your 
platform [here](https://github.com/smallwat3r/shhh-cli/releases).

```sh
# Unzip and move to bin
unzip shhh-cli-darwin-amd64-0.1.0.zip
sudo mv shhh /usr/local/bin/
```

## Usage  

```
Usage: shhh <command> [<args>]

Commands available:
  create    Creates a secret message.
  read      Read a secret message.

Usage of create:
  -m string
        Secret message to encrypt.
  -p string
        Passphrase to encrypt secret.
  -d int
        Optional, number of days to keep the secret alive. (default 3)

Usage of read:
  -l string
        URL link to access secret.
  -p string
        Passphrase to decrypt secret.

Examples:
    shhh create -m "this is a secret msg." -p SuperPassphrase123 -d 2
    shhh read -l https://shhh-encrypt.com/api/r/jKD8Uy0A9_51c8asqAYL -p SuperPassphrase123
```
