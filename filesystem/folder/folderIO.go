package folder

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/karrick/godirwalk"
	"github.com/piperun/hashsync/filesystem/file"
	"github.com/piperun/hashsync/hashdb"
	"github.com/piperun/hashsync/hashfunc"
	"github.com/piperun/hashsync/netcom"
)

type HostData map[string]interface{}

// folder datatype uses the folder path as key (might be changed to hash) and another datatype: folderData as it's value
type Folder map[string]folderData

type folderData struct {
	Files map[string]string
	Hash  string
}

var hostID = make(hashdb.HostID)

func SavetoCache(IP string) {
	var (
		collection hashdb.DBContent
		doc        = make(map[string]interface{})
	)
	collection.LoadCollection("Cache")
	doc[IP] = hostID[IP]
	collection.AddDocument(doc)
}

func LocalWalk(root string) {
	const hostIP = "127.0.0.1"
	var (
		dirpath, dirname, prevdir string
		local_file                file.LocalFile
		file_hashsum, dir_hashsum hashfunc.HashSum
		HostFS                    hashdb.DBContent
		hostData                  = make(HostData)
		dir                       = make(Folder)
		sub_files                 = make(map[string]string)
	)
	HostFS.LoadCollection("HostFS")

	err := godirwalk.Walk(root, &godirwalk.Options{
		Callback: func(currpath string, de *godirwalk.Dirent) error {
			if de.IsDir() == true {
				if dirpath != currpath && dirpath != "" {
					hostData[hostIP] = createhostData(dirpath, dir_hashsum.Hex, sub_files)
					prevdir = dirpath

					hostID[hostIP] = append(hostID[hostIP], HostFS.AddDocument(hostData))
					delete(dir, dirpath)
					sub_files = make(map[string]string)
					debug.FreeOSMemory()
				}
				dirpath = currpath
				dirname = de.Name()
				dir_hashsum = createhashSum(dirname)

			} else if de.IsRegular() == true && de.IsSymlink() == false {

				local_file.Name = currpath
				local_file.Open()
				file_hashsum = createhashSum(local_file)
				local_file.Close()

				sub_files[file_hashsum.Hex] = de.Name()

			}

			return nil
		},
		ErrorCallback: func(currpath string, err error) godirwalk.ErrorAction {

			// For the purposes of this example, a simple SkipNode will suffice,
			// although in reality perhaps additional logic might be called for.
			return godirwalk.SkipNode
		},
		Unsorted: true, // (optional) set true for faster yet non-deterministic enumeration (see godoc)
	})
	fmt.Print(err)
	if prevdir == "" {

		hostData[hostIP] = createhostData(dirpath, dir_hashsum.Hex, sub_files)
		hostID[hostIP] = append(hostID[hostIP], HostFS.AddDocument(hostData))
	}
	SavetoCache(hostIP)

}

func RemoteWalk(sftpcon netcom.Connection, root string, hostIP string) {
	var (
		prevdir, dirname, dirpath string
		remote_file               file.RemoteFile
		file_hashsum, dir_hashsum hashfunc.HashSum
		HostFS                    hashdb.DBContent
		hostData                  = make(HostData)
		dir                       = make(Folder)
		sub_files                 = make(map[string]string)
	)
	sftp := sftpcon.Client.GetSFTPConnection()
	HostFS.LoadCollection("HostFS")

	walk := sftp.Walk(root)

	for walk.Step() {

		if perm := walk.Stat().Mode().Perm(); perm&(1<<2) == 0 {
			continue
		}

		if walk.Stat().IsDir() {
			if dirpath != walk.Path() && dirpath != "" {
				hostData[hostIP] = createhostData(dirpath, dir_hashsum.Hex, sub_files)
				prevdir = dirpath

				hostID[hostIP] = append(hostID[hostIP], HostFS.AddDocument(hostData))
				delete(dir, dirpath)
				sub_files = make(map[string]string)
				debug.FreeOSMemory()
			}

			dirname = walk.Stat().Name()

			dirpath = walk.Path()
			dir_hashsum = createhashSum(dirname)
		} else if walk.Stat().Mode().IsRegular() && walk.Stat().Mode()&os.ModeSymlink == 0 {

			remote_file.Name = walk.Path()
			remote_file.Open(sftp)
			file_hashsum = createhashSum(remote_file)
			remote_file.Close()

			sub_files[file_hashsum.Hex] = walk.Stat().Name()

		}

	}
	if prevdir == "" {

		hostData[hostIP] = createhostData(dirpath, dir_hashsum.Hex, sub_files)
		hostID[hostIP] = append(hostID[hostIP], HostFS.AddDocument(hostData))
	}
	SavetoCache(hostIP)
}

func createhashSum(fsObject interface{}) (sum hashfunc.HashSum) {
	var hashobject hashfunc.Object
	hashobject.Source = fsObject
	switch fsObject.(type) {
	case string:
		sum = hashobject.SHA256()
	case file.RemoteFile:
		sum = hashobject.CRC32()
	case file.LocalFile:
		sum = hashobject.CRC32()
	}
	return
}

func createhostData(path, folder_hash string, sub_files map[string]string) (dir Folder) {
	dir = make(Folder)
	dir[path] = folderData{
		Files: sub_files,
		Hash:  folder_hash,
	}
	return

}
