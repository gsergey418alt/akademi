package core

import (
	"log"

	mrand "math/rand"
	"time"
)

// AkademiNode constants.
const (
	IDLength   = 20
	BucketSize = 20
)

// List of bootstrap nodes used for first connecting to
// the network.
var BootstrapHosts = [...]Host{
	"akademi_bootstrap:3865",
}

// Akademi uses node and key IDs, whose length is defined
// in bytes by IDLength.
type BaseID [IDLength]byte

// Separate IPPort type because the IP address is
// identified by receiving node.
type IPPort uint16

// Host is used to identify node's IP address and
// port.
type Host string

// DataBytes is a type for values to be stored in akademi
// nodes.
type DataBytes []byte

// RoutingEntry is a structure that stores routing
// information about an akademi node.
type RoutingEntry struct {
	Host   Host
	NodeID BaseID
}

// AkademiNode is a structure containing the core kademlia
// logic.
type AkademiNode struct {
	NodeID        BaseID
	ListenPort    IPPort
	KeyValueStore map[BaseID][]byte

	RoutingTable [IDLength * 8][]RoutingEntry

	Dispatcher Dispatcher
}

// The initialize function assigns a random NodeID to the
// AkademiNode.
func (a *AkademiNode) Initialize(dispatcher Dispatcher, listenPort IPPort, bootstrap bool) error {
	a.ListenPort = listenPort
	a.NodeID = RandomBaseID()
	a.Dispatcher = dispatcher
	err := a.Dispatcher.Initialize(RoutingHeader{ListenPort: a.ListenPort, NodeID: a.NodeID})
	if err != nil {
		return err
	}
	if bootstrap {
		i := mrand.Intn(len(BootstrapHosts))
		var header RoutingHeader
		for header, err = a.Dispatcher.Ping(BootstrapHosts[i]); err != nil; {
			log.Print(err)
			time.Sleep(5 * time.Second)
		}
		log.Print("Connected to bootstrap node \"", BootstrapHosts[i], "\". NodeID: ", header.NodeID.Base64Str())
	}
	return nil
}
