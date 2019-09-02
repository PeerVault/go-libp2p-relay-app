# `go-libp2p-relay-app` Relay ready to use

## What is libp2p

libp2p is a modular system of protocols, specifications and libraries that enable the development of peer-to-peer network applications.

More information at [https://docs.libp2p.io/introduction/what-is-libp2p/](https://docs.libp2p.io/introduction/what-is-libp2p/)

## Relay application

Circuit relay is a transport protocol that routes traffic between two peers over a third-party “relay” peer.

This application is a very simple relay that generates a keypair and saves it as a JSON file.
You can easily re-use and start the relay with the same ID and listening port.

## Usage

```bash
## Build the relay application
$ make build
go build -o bin/relay relay.go

## Start the relay listening TCP 3030 and create a new KeyPair
$ ./bin/relay -p 3030
{
 "Name": "Relay Lich 1.0",
 "Id": "QmPeerId",
 "PrivKey": "XXX",
 "PubKey": "XXX"
}
Relay start with ID: QmPeerId
Relay listen on: [/ip4/127.0.0.1/tcp/3030 /ip4/127.94.0.1/tcp/3030 /ip4/192.168.178.53/tcp/3030]

## Start relay with existing json file configuration
$ ./bin/relay -p 3030 -c relay-peer-key.json
Relay start with ID: QmPeerId
Relay listen on: [/ip4/127.0.0.1/tcp/3030 /ip4/127.94.0.1/tcp/3030 /ip4/192.168.178.53/tcp/3030]
```
