package hashfunc

import (
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"hash/crc32"
	"io"
	"log"
	"os"
)

func openFile(filename string, fp **os.File) {
	file, err := os.Open(filename)
	if err != nil {
		log.Print(err)
	}
	*fp = file
}

// CRC32 hash function
func CRC32(filename string, hashversion ...string) string {
	var (
		versions   = make(map[string]uint32)
		table      *crc32.Table
		hashstring string
		file       *os.File
	)

	openFile(filename, &file)

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

	hash := crc32.New(table)

	if _, err := io.Copy(hash, file); err != nil {
		log.Print(err)
	}

	hashstring = hex.EncodeToString(hash.Sum(nil)[:])
	if false {
		log.Print(hashstring)
	}

	defer file.Close()
	return hashstring

}

//SHA256 hash function
func SHA256(str string, hashversion ...string) string {
	var sha hash.Hash

	if hashversion == nil || len(hashversion[0]) <= 0 {

	}
	sha = sha256.New()
	sha.Write([]byte(str))
	if false {
		log.Print(sha.Sum(nil), hex.EncodeToString(sha.Sum(nil)[:]))
	}
	return hex.EncodeToString(sha.Sum(nil)[:])

}

//SHA256 file hash function
func SHA256FILE(file *os.File, hashversion ...string) {

}
