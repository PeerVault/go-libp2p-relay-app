package main

import (
	"context"
	"fmt"
	"flag"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"

	ma "github.com/multiformats/go-multiaddr"
)

// create peer addr info
func p2pAddrInfo(addrStr string) (*peer.AddrInfo, error) {
	addr, err := ma.NewMultiaddr(addrStr)
	if err != nil {
		panic(err)
	}
	return peer.AddrInfoFromP2pAddr(addr)
}

func main() {
  relayHost := flag.String("relay", "", "Relay Host URL")
  flag.Parse()

	// Zero out the listen addresses for the host, so it can only communicate
	// via p2p-circuit for our example
	node, err := libp2p.New(
	  context.Background(),
	  libp2p.ListenAddrs(),
	  libp2p.EnableRelay(),
  )
	if err != nil {
		panic(err)
	}

	// Creates relay peer.AddrInfo
  relayAddrInfo, err := p2pAddrInfo(*relayHost)
  if err != nil {
    panic(err)
  }

	if err := node.Connect(context.Background(), *relayAddrInfo); err != nil {
		panic(err)
	}

	// Now, to test things, let's set up a protocol handler on node
	node.SetStreamHandler("/cats", func(s network.Stream) {
		fmt.Println("Meow! It worked!")
		s.Close()
	})

	fmt.Println("Bob Relay ID:", node.ID().Pretty());

	for _, addr := range node.Addrs() {
    fmt.Printf("Bob Relay Addr: %s\n", addr.String())
  }

	select {}
}
