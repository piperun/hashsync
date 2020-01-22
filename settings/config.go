package settings

import (
	"log"
	"os"

	"github.com/pelletier/go-toml"
	"github.com/piperun/hashsync/filesystem/file"
)

// SSH options
type SSH struct {
	IP       string `toml:"conn_ip"`
	Port     string `toml:"conn_port"`
	Key      string `toml:"key" commented:"true" comment:"Placeholder not Implemented yet"`
	Password string `toml:"password"`
	User     string `toml:"user"`
}

// Path options
type Paths struct {
	Database   string `toml:"database"`
	Root       string `toml:"root"`
	RemoteRoot string `toml:"remoteroot"`
}

// General options
type General struct {
	MThreaded      bool     `toml:"MultiThread"`
	BlacklistedDir []string `toml:"Blacklist_dir"`
}

// Body of the config file
type Config struct {
	General General `toml:"General" comment:"General configuration"`
	SSH     SSH     `toml:"SSH" comment:"SSH configuration"`
	Paths   Paths   `toml:"Paths" comment:"File path configuration"`
}

/*
	Content holds the content of the config file + it's path which is set either by the user or by Setup()
	Path: Takes a string to represent the config location, it'll be defaulted to ./config.toml if none is given by the user
	Data: Uses a struct to store the data inside config.toml
	Tree: Uses a pointer to toml.Tree to store the data inside config.toml
*/
type Content struct {
	Path   string
	Struct Config     // ONLY TO BE USED WITH LoadStruct()
	Tree   *toml.Tree // ONLY TO BE USED WITH Load()
}

func (content *Content) Query(object string) string {
	confobject := content.Tree.Get(object).(string)
	if confobject != "" {
		return confobject
	}
	return ""
}

//Loads the file as a toml.Tree
func (content *Content) LoadTree() {
	var err error
	content.Tree, err = toml.LoadFile(content.Path)
	if err != nil {
		log.Fatal(err)
	}
}

func (content *Content) Setup(path string) {
	content.Path = path
	if !file.CheckFile(path) && len(path) > 0 {
		content.createConf()
	} else {

	}
}

// Local functions

func (content *Content) createConf() {
	file, err := os.OpenFile(content.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	config := Config{
		General{MThreaded: true},
		SSH{User: "user", Password: "password", Key: "", IP: "0.0.0.0", Port: "22"},
		Paths{Root: "", RemoteRoot: ""},
	}

	encoder := toml.NewEncoder(file)
	encoder.Encode(config)

}

/*
	b, err := toml.Marshal(config)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
	Add this to test to later
*/

// LoadStruct() Will load the config file as a Config{} struct
// Cannot be used with Query()
func (content *Content) LoadStruct() {
	file, err := os.OpenFile(content.Path, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal(content.Path, err)
	}
	defer file.Close()

	decoder := toml.NewDecoder(file)
	err = decoder.Decode(&content.Struct)
	if err != nil {
		log.Fatal(err)
	}
}
