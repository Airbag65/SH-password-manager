# Passport
Passport is an extremely simple self hosted password manager, that consists of a backend/server(WIP), a TUI(WIP), a CLI(TBD) and an adroid application(TBD).

# Getting started
If it is the first time you are using `passport`, you will need to sign a SSL certificate, for the HTTPS server to run. The
easiest way to do this is by running this command while in the `backend` directory.
```bash
$ openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -sha256 -days 365 -nodes
```
This certificate will however expire after one year so, if you are planing to use passport for a longer period of time, be prepared to renew it using the same command in a years time.

 TODO: write the rest
