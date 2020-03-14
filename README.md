# Shhh-cli  

shhh-cli is a Command Line Interface tool to interact with [Shhh](https://github.com/smallwat3r) web-application API.  
This tool allows you to create and read secrets directly from the command line / terminal.  

## Use with your own Shhh server  

shhh-cli interacts by default with the official Shhh web-application (shhh-encrypt.com)
If you're hosting your own Shhh server, you can set-up an environment variable `SHHH_SERVER`.
This server will then be used as the default one by shhh-cli.  
Example: `export SHHH_SERVER=https://mycustomserver.com`  

## Usage  

```sh
usage: shhh <command> [<args>]
Commands available:
  create    Creates a secret message.
  read      Read a secret message
Examples:
  shhh create -m 'this is a secret msg.' -p SuperPassphrase123 -d 2
  shhh read -l https://shhh-encrypt.com/api/r/jKD8Uy0A9_51c8asqAYL -p SuperPassphrase123
```

### create  
```sh
Usage of create:
  -d int
        Optional, number of days to keep the secret alive (defaults to 3 days).
  -m string
        Secret message to encrypt.
  -p string
        Passphrase to encrypt secret.
```

### read  
```sh
Usage of read:
  -l string
        URL link to access secret.
  -p string
        Passphrase to decrypt secret.
```
