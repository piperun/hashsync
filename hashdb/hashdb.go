package hashdb

import (
	"fmt"
	"log"

	"github.com/HouzuoGuo/tiedot/data"
	"github.com/HouzuoGuo/tiedot/db"

	_ "github.com/mattn/go-sqlite3"
)

const default_path = "./database"

var database = NewDB(default_path)

type DBContent struct {
	Collection *db.Col
	Name       string
}

type HostID map[string][]int

/*
	Default values for Config is:
		DefaultDocMaxRoom = 2 * 1048576 // DefaultDocMaxRoom is the default maximum size a single document may never exceed.
		DocHeader         = 1 + 10      // DocHeader is the size of document header fields.
		EntrySize         = 1 + 10 + 10 // EntrySize is the size of a single hash table entry.
		BucketHeader      = 10          // BucketHeader is the size of hash table bucket's header fields.

		DocMaxRoom:    DefaultDocMaxRoom,
		ColFileGrowth: COL_FILE_GROWTH,
		PerBucket:     16,
		HTFileGrowth:  HT_FILE_GROWTH,
		HashBits:      HASH_BITS,

*/

func (dbcontent *DBContent) LoadCollection(name string) {

	dbcontent.Collection = database.Use(name)
	dbcontent.Name = name

}

func (dbcontent *DBContent) AddDocument(document map[string]interface{}) int {
	ID, err := dbcontent.Collection.Insert(document)
	if err != nil {
		log.Fatal(err)
	}
	return ID
}

func (dbcontent *DBContent) AddDocID(document map[string]interface{}, id int) {
	err := dbcontent.Collection.InsertRecovery(id, document)
	if err != nil {
		log.Fatal(err)
	}
}

func (dbcontent *DBContent) RemoveCollection(name string) {
	database.Drop(name)
}

func initConfig(path string) {
	var config data.Config

	config = data.Config{
		DocMaxRoom:    data.DefaultDocMaxRoom,
		ColFileGrowth: data.COL_FILE_GROWTH,
		PerBucket:     16,
		HTFileGrowth:  1 * 1048576,
		HashBits:      11,
	}
	data.SetCustomConfig(config)
}

func NewDB(path string) *db.DB {

	initConfig(path)
	database, err := db.OpenDB(path)

	if err != nil {
		log.Fatal(err)
	}
	if database == nil {
		log.Fatal("Database file is nil")
	}
	return database

}

func CreateCollection(name string) {
	err := database.Create(name)
	if err != nil {
		log.Fatal(err)
	}
}

func CheckCollectionExists(name string) bool {
	return database.ColExists(name)
}

func AddDocument(document map[string]interface{}, collection string) {
	hashdb := database.Use(collection)

	_, err := hashdb.Insert(document)
	if err != nil {
		log.Fatal(err)
	}
}

func PrintCollection(collection string) {
	hashdb := database.Use(collection)
	hashdb.ForEachDoc(func(id int, docContent []byte) (willMoveOn bool) {
		fmt.Println("Document", id, "is", string(docContent))
		return true // move on to the next document OR
		//return false // do not move on to the next document
	})
}
