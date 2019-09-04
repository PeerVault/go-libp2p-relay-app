package main

import (
	"context"
	"fmt"
	"flag"
	"os"
	"log"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/peer"

	circuit "github.com/libp2p/go-libp2p-circuit"
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
  recipient := flag.String("recipient", "", "Recipient Peer ID")
  flag.Parse()

  if *relayHost == "" {
		log.Fatal("Please provide relay host with --relay option")
		os.Exit(0)
	}
  if *recipient == "" {
		log.Fatal("Please provide recipient with --recipient option")
		os.Exit(0)
	}

	// Tell the host to discover peer through relay.
	node, err := libp2p.New(context.Background(), libp2p.EnableRelay(circuit.OptDiscovery))
	if err != nil {
		panic(err)
	}

	// Creates relay peer.AddrInfo
  relayAddrInfo, err := p2pAddrInfo(*relayHost)
  if err != nil {
    panic(err)
  }

	// Connect node to relay
	if err := node.Connect(context.Background(), *relayAddrInfo); err != nil {
		fmt.Println("fail connect to relay!")
		panic(err)
	}

  recipientPeerid, err := peer.IDB58Decode(*recipient)
  if err != nil {
    panic(err)
  }

  // SwapToP2pMultiaddrs is a function to make the transition from /ipfs/... multiaddrs to /p2p/... multiaddrs easier The first stage of the rollout is to ship this package to all users so that all users of multiaddr can parse both /ipfs/ and /p2p/ multiaddrs as the same code (P_P2P).
  ma.SwapToP2pMultiaddrs()
	relayaddr, err := ma.NewMultiaddr(*relayHost + "/p2p-circuit/p2p/" + recipientPeerid.Pretty())

  if err != nil {
    panic(err)
  }

	recipientRelayInfo := peer.AddrInfo{
		ID: recipientPeerid,
		Addrs: []ma.Multiaddr{relayaddr},
	}

	// Connect node to recipient using relay
	if err := node.Connect(context.Background(), recipientRelayInfo); err != nil {
		fmt.Println("fail connect to recipient")
		panic(err)
	}

	// Woohoo! we're connected!
	s, err := node.NewStream(context.Background(), recipientPeerid, "/cats")
	if err != nil {
		fmt.Println("huh, this should have worked: ", err)
		return
	}

	s.Read(make([]byte, 1)) // block until the handler closes the stream

	fmt.Println("end")
}
