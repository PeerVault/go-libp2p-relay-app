package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"flag"
	"log"

	"encoding/json"
	"io/ioutil"
	b64 "encoding/base64"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	peer "github.com/libp2p/go-libp2p-core/peer"
	tcp "github.com/libp2p/go-tcp-transport"
	circuit "github.com/libp2p/go-libp2p-circuit"
)

type PeerIdentityJson struct {
  Name string
  Id string
  PrivKey string
  PubKey string
}

// Create json file with identity information
func writeIdentityJson(privKey crypto.PrivKey) {
  ID, err := peer.IDFromPrivateKey(privKey)
  if err != nil {
    panic(err)
  }
  _ = ID

  pvtBytes, err := privKey.Raw()
  if err != nil {
    panic(err)
  }
  _ = pvtBytes
  pubBytes, err := privKey.GetPublic().Raw()
  if err != nil {
    panic(err)
  }
  _ = pubBytes

  identityJson := &PeerIdentityJson {
    Name: "Relay Lich 1.0",
    Id: ID.Pretty(),
    PrivKey: b64.StdEncoding.EncodeToString(pvtBytes),
    PubKey: b64.StdEncoding.EncodeToString(pubBytes),
  }

	file, err := json.MarshalIndent(identityJson, "", " ")
	if err != nil {
    panic(err)
  }

  fmt.Println(string(file))

	_ = ioutil.WriteFile("relay-peer-key.json", file, 0644)
}

func readIdentityJson(filePath string) (crypto.PrivKey, crypto.PubKey, error) {
  // Open our jsonFile
  jsonFile, err := os.Open(filePath)
  // if we os.Open returns an error then handle it
  if err != nil {
    return nil, nil, err
  }
  // defer the closing of our jsonFile so that we can parse it later on
  defer jsonFile.Close()

  // read our opened xmlFile as a byte array.
  byteValue, _ := ioutil.ReadAll(jsonFile)

  // we initialize our Users array
  var peerIdentity PeerIdentityJson

  // we unmarshal our byteArray which contains our
  // jsonFile's content into 'users' which we defined above
  json.Unmarshal(byteValue, &peerIdentity)

  privKeyByte, err := b64.StdEncoding.DecodeString(peerIdentity.PrivKey)
  if err != nil {
    return nil, nil, err
  }

  pvtKey, err := crypto.UnmarshalRsaPrivateKey(privKeyByte)
  if err != nil {
    return nil, nil, err
  }

  return pvtKey, pvtKey.GetPublic(), nil
}

func startRelay(listenPort int, priv crypto.PrivKey) {
  // The context governs the lifetime of the libp2p node
  ctx, cancel := context.WithCancel(context.Background())
  defer cancel()

  // Transport configuration
	transports := libp2p.ChainOptions(
		libp2p.Transport(tcp.NewTCPTransport),
	)

  // Listen address with port
  listenAddrs := libp2p.ListenAddrStrings(
    fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", listenPort),
  )

	relayNode, err := libp2p.New(
	  ctx,
	  libp2p.Identity(priv),
	  transports,
	  listenAddrs,
	  libp2p.EnableRelay(circuit.OptHop),
  )
	if err != nil {
		panic(err)
	}

  fmt.Printf("Relay start with ID: %s\n",relayNode.ID().Pretty())
  fmt.Printf("Relay listen on: %s\n",relayNode.Addrs())

  // Start relay indefinitely until signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT)

  select {
	case <-stop:
	  fmt.Printf("Stop relay")
		relayNode.Close()
		os.Exit(0)
  }
}

func main() {
	listenPort := flag.Int("p", 0, "TCP listen port")
	configPeerIdentity := flag.String("c", "", "Config Peer identity JSON File")
	flag.Parse()

  if *listenPort == 0 {
		log.Fatal("Please provide a port to bind on with -p")
		os.Exit(0)
	}

  if *configPeerIdentity == "" {
    priv, _, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
    if err != nil {
      panic(err)
    }
    writeIdentityJson(priv)
    startRelay(*listenPort, priv)
  } else {
    priv, _, err := readIdentityJson(*configPeerIdentity)
    if err != nil {
      panic(err)
    }
    startRelay(*listenPort, priv)
  }
}
