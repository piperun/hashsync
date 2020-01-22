package hashfunc

import (
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"hash/crc32"
	"io"
	"log"

	"github.com/piperun/hashsync/filesystem/file"
)

type Object struct {
	hash32 hash.Hash32
	hash   hash.Hash
	Source interface{}
}

type HashSum struct {
	Hex string
	Int uint32
}

// CRC32 hash function
func (object *Object) CRC32(hashversion ...string) HashSum {
	var (
		versions = make(map[string]uint32)
		table    *crc32.Table
		hashsum  HashSum
	)

	versions["C"] = 0x82F63B78
	versions["K"] = 0xEB31D82E
	versions["K2"] = 0x992C1A4C
	versions["Q"] = 0xD5828281
	versions["DEFAULT"] = 0xEDB88320
	if hashversion == nil || len(hashversion[0]) <= 0 {
		table = crc32.MakeTable(versions["DEFAULT"])
	} else {
		table = crc32.MakeTable(versions[hashversion[0]])
	}

	object.hash32 = crc32.New(table)
	object.convertCorrectType()

	hashsum.Hex = hex.EncodeToString(object.hash32.Sum(nil)[:])
	hashsum.Int = object.hash32.Sum32()
	if false {
		log.Print(hashsum.Hex)
	}

	return hashsum

}

//SHA256 hash function
func (object *Object) SHA256(hashversion ...string) HashSum {
	var (
		hashsum HashSum
	)

	if hashversion == nil || len(hashversion[0]) <= 0 {

	}
	object.hash = sha256.New()
	object.convertCorrectType()

	if false {
		log.Print(object.hash.Sum(nil), hex.EncodeToString(object.hash.Sum(nil)[:]))
	}
	hashsum.Hex = hex.EncodeToString(object.hash.Sum(nil)[:])
	return hashsum
}

// Local functions

func (object *Object) convertCorrectType() {

	// atm this ugly hack is the only to make it work.
	switch filetype := object.Source.(type) {
	case file.LocalFile:
		if _, err := io.Copy(object.hash32, filetype.Data); err != nil {
			log.Print(err)
		}
		filetype.Close()

	case file.RemoteFile:
		if _, err := io.Copy(object.hash32, filetype.Data); err != nil {
			log.Print(err)
		}
		filetype.Close()
	case string:

		object.hash.Write([]byte(filetype))
	case []byte:
		object.hash.Write(filetype)
	}
}
