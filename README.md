# Shhh-cli  

**shhh-cli** is a Command Line Interface tool to interact with [Shhh](https://github.com/smallwat3r) web-application API.  
This tool allows you to create and read secrets directly from the command line / terminal.  

## Use with your own Shhh server  

shhh-cli interacts by default with the official Shhh server (shhh-encrypt.com).  
If you host Shhh on your own server, you can set-up an env variable `SHHH_SERVER`.  
It will know interact with your server as default.  
```sh
# Example (in you bashrc or zshrc)
export SHHH_SERVER=https://mycustomserver.com
```

## Usage  

```
Usage: shhh-cli <command> [<args>]

Commands available:
  create    Creates a secret message.
  read      Read a secret message

Usage of create:
  -d int
        Optional, number of days to keep the secret alive (defaults to 3 days).
  -m string
        Secret message to encrypt.
  -p string
        Passphrase to encrypt secret.

Usage of read:
  -l string
        Optional, number of days to keep the secret alive (defaults to 3 days).
  -p string
        Passphrase to decrypt secret.

Examples:
    shhh-cli create -m 'this is a secret msg.' -p SuperPassphrase123 -d 2
    shhh-cli read -l https://shhh-encrypt.com/api/r/jKD8Uy0A9_51c8asqAYL -p SuperPassphrase123
```
