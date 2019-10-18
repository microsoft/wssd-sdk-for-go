# Generating Certs

Here you will learn how to generate certs and what to do with them when you deploy

## Prerequisites

### CFSSL

There are many ways to generate certs. For now I found the community really enjoys cloudflare's tool `cfssl`

These are the steps I used to setup `cfssl`

First make sure your `$GOOS` is set to linux and `$GOPATH` is properly set up

```
$ git clone git@github.com:cloudflare/cfssl.git $GOPATH/src/github.com/cloudflare/cfssl
$ cd $GOPATH/src/github.com/cloudflare/cfssl
$ make
```

The resulting binaries will be in the bin folder:

```
$ tree bin
bin
├── cfssl
├── cfssl-bundle
├── cfssl-certinfo
├── cfssl-newkey
├── cfssl-scan
├── cfssljson
├── mkbundle
└── multirootca

0 directories, 8 files
```

Then copy `cfssl` and `cfssljson` into your `$GOBIN`

```
$ cp cfssl $GOROOT/bin
$ cp cfssljson $GOROOT/bin
```

Now from this directory run 
`$ ./gen_certs.sh`

The resulting cert directory should look like:

```
$ tree certs
certs
├── ca-config.json
├── ca.csr
├── ca-csr.json
├── ca-key.pem
├── ca.pem
├── gen_certs.sh
├── README.md
├── wssdagent-csr.json
├── wssd.csr
├── wssd-key.pem
└── wssd.pem

0 directories, 11 files
```

For more information go to [cfssl installation guide](https://github.com/cloudflare/cfssl#installation).

## Deploying wssd-sdk-for-go and wssdagent (Work in progress)

For now, for deploying with certs you will need to copy over `wssd.pem` for using wssd-sdk-for-go (wssdctl.exe) and both `wssd-key.pem` and `wssd.pem` for wssdagent.exe

For wssd-sdk-for-go set `WSSD_CLIENT_TOKEN` with the full path to `wssd.pem`
(By default it will check your working directory for the .pem)


For wssdagent.exe write the `TLSCertPath` and `TLSKeyPath` variables in the  `.yaml` config file.
(By default it will check your working directory for both `.pem`s)
