package main

import (
	"log"

	"github.com/piperun/hashsync/netcom"
	"github.com/piperun/hashsync/registry"

	"github.com/piperun/hashsync/config"
	"github.com/piperun/hashsync/hashfunc"
)

type OSvars struct {
	separator string
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	var (
		filename     string
		conf_content config.Content
	)
	connection := netcom.Connection{}

	registry.HandleFlags()

	conf_content.Setup(registry.GetConfigPath())
	conf_content.LoadTree()

	filename = ""
	hashfunc.CRC32(filename)
	hashfunc.SHA256(filename)
	connection.Connect(conf_content)
	connection.Client.Disconnect()

}
