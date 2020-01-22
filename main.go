package main

import (
	"log"

	"github.com/piperun/hashsync/hashdb"

	"github.com/piperun/hashsync/filesystem/folder"

	"github.com/piperun/hashsync/netcom"
	"github.com/piperun/hashsync/settings"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	var (
		conf_content settings.Content
	)
	connection := netcom.Connection{}

	settings.HandleFlags()

	conf_content.Setup(settings.GetConfigPath())
	conf_content.LoadTree()

	netcom.InsertConfigContent(conf_content)
	connection.Connect()

	initDB("HostFS", "Cache")
	//folder.LocalWalk(conf_content.Query("Paths.root"))
	folder.RemoteWalk(connection, conf_content.Query("Paths.remoteroot"), conf_content.Query("SSH.conn_ip"))
	hashdb.PrintCollection("HostFS")
	hashdb.PrintCollection("Cache")
	//log.Print(hashobj.CRC32())

	connection.Disconnect()

}

func initDB(collections ...string) {
	for _, col := range collections {
		if !hashdb.CheckCollectionExists(col) {
			hashdb.CreateCollection(col)
		}
	}

}
