package filesystem

import (
	"log"
	"os"

	"github.com/pkg/sftp"
)

type RemoteFile struct {
	Name string
	Data *sftp.File
}

type LocalFile struct {
	Name string
	Data *os.File
}

// Local File methods

func (local *LocalFile) GetFile() *os.File {
	return local.Data
}

func (local *LocalFile) Open() {
	var err error
	local.Data, err = os.Open(local.Name)
	if err != nil {
		log.Print(err)
	}
}

func (local *LocalFile) Close() {
	local.Data.Close()
}

// Remote File methods

func (remote *RemoteFile) GetFile() *sftp.File {
	return remote.Data
}

func (remote *RemoteFile) Open(client *sftp.Client) {
	var err error
	remote.Data, err = client.Open(remote.Name)
	if err != nil {
		log.Fatal(err)
	}
}

func (remote *RemoteFile) Close() {
	remote.Data.Close()
}

// Global functions

func CheckFile(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}
