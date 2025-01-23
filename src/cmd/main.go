package main

import (
	"fmt"
	"log"
	"net/rpc"
	"os"

	"github.com/gsergey418alt/akademi/core"
	"github.com/gsergey418alt/akademi/daemon"
	akademiRPC "github.com/gsergey418alt/akademi/rpc"
)

const (
	// Never expose RPC to the public! For docker.
	defaultRpcListenAddr  = "127.0.0.1:3855"
	defaultNodeListenAddr = "0.0.0.0:3865"
)

// Settings populated by parseArgs()
type cmdOptions struct {
	cmd            string
	target         string
	rpcListenAddr  string
	nodeListenAddr string
	bootstrap      bool
}

// Global instance of cmdOptions
var opts cmdOptions

// Commands with no positional arguments.
var noPosArgs map[string]bool

// The function parseArgs is responsible for command line
// argument parsing.
func parseArgs() {
	noPosArgs = map[string]bool{
		"daemon":        true,
		"routing_table": true,
		"info":          true,
	}

	opts.bootstrap = true
	opts.nodeListenAddr = defaultNodeListenAddr
	opts.rpcListenAddr = defaultRpcListenAddr

	argLen := len(os.Args)
	if argLen < 2 {
		fmt.Print("Not enough arguments, please provide a command.\n")
		os.Exit(1)
	}
	optStart, optStop := 2, argLen
	opts.cmd = os.Args[1]
	if _, ok := noPosArgs[opts.cmd]; !ok {
		opts.target = os.Args[argLen-1]
		optStop--
	}
	for i := optStart; i < optStop; i++ {
		switch os.Args[i] {
		case "--no-bootstrap":
			opts.bootstrap = false
		case "--rpc-addr":
			opts.rpcListenAddr = os.Args[i+1]
			i++
		default:
			fmt.Print("Unknown argument: \"", os.Args[i], "\".\n")
			os.Exit(1)
		}
	}
	return
}

// Wrapper for RPC calls.
func RPCSessionManager(f func(client *rpc.Client) error) {
	client, err := rpc.DialHTTP("tcp", opts.rpcListenAddr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer client.Close()
	err = f(client)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Akademi entrypoint.
func main() {
	parseArgs()
	switch opts.cmd {
	case "daemon":
		log.Fatal(daemon.Daemon(opts.nodeListenAddr, opts.bootstrap, opts.rpcListenAddr))
	case "ping":
		args := akademiRPC.PingArgs{Host: core.Host(opts.target)}
		reply := akademiRPC.PingReply{}
		RPCSessionManager(func(client *rpc.Client) error {
			return client.Call("AkademiNodeRPCServer.Ping", args, &reply)
		})
		fmt.Print("Received reply from ", opts.target, ". NodeID: ", reply.Header.NodeID, ".\n")
	case "lookup":
		id, err := core.B32ToID(opts.target)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		args := akademiRPC.LookupArgs{ID: id}
		reply := akademiRPC.LookupReply{}
		RPCSessionManager(func(client *rpc.Client) error {
			return client.Call("AkademiNodeRPCServer.Lookup", args, &reply)
		})
		fmt.Print("Node located successfully. Address: ", reply.RoutingEntry, ".\n")
	case "routing_table":
		args := akademiRPC.RoutingTableArgs{}
		reply := akademiRPC.RoutingTableReply{}
		RPCSessionManager(func(client *rpc.Client) error {
			return client.Call("AkademiNodeRPCServer.RoutingTable", args, &reply)
		})
		fmt.Print("Node routing table:\n", reply.RoutingTable, "\n")
	case "info":
		args := akademiRPC.NodeInfoArgs{}
		reply := akademiRPC.NodeInfoReply{}
		RPCSessionManager(func(client *rpc.Client) error {
			return client.Call("AkademiNodeRPCServer.NodeInfo", args, &reply)
		})
		fmt.Print("Node information:\n", reply.NodeInfo, "\n")
	case "bootstrap":
		args := akademiRPC.BootstrapArgs{Host: core.Host(opts.target)}
		reply := akademiRPC.BootstrapReply{}
		RPCSessionManager(func(client *rpc.Client) error {
			return client.Call("AkademiNodeRPCServer.Bootstrap", args, &reply)
		})
		fmt.Print("Successfully bootstrapped node with ", opts.target, ".\n")
	default:
		fmt.Print("Command \"", opts.cmd, "\" not found.\n")
		os.Exit(1)
	}
}
