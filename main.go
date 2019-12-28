package main

import (
	"github.com/piperun/hashsync/hashfunc"
	"github.com/piperun/hashsync/netcom"
)

type OSvars struct {
	separator string
}

func main() {
	var (
		filename string
	)

	filename = ""
	hashfunc.CRC32(filename)
	hashfunc.SHA256(filename)
	netcom.SFTPConnect()

}
