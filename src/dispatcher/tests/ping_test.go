package tests

import (
	"log"
	"testing"

	"github.com/gsergey418alt/akademi/core"
	"github.com/gsergey418alt/akademi/dispatcher"
)

func TestPing(t *testing.T) {
	d := &dispatcher.UDPDispatcher{}
	d.Initialize(core.IPPort(3865))
	nodeID, err := d.Ping(core.Host("127.0.0.1:3865"))
	if err != nil {
		panic(err)
	}
	log.Print("NodeID: ", nodeID.BinStr())
}
