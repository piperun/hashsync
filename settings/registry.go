package settings

import (
	"flag"
	"runtime"

	"github.com/pelletier/go-toml"
)

// Register is the software's internal configuration register, this is to make it easier to allow send data in-between functions that need opinated outside data.
// Note: SSH/Paths configuration is excempt since the data is too opinated.
type Register struct {
	ConfigPath string // Where the config path shall be
	MThreads   bool   // If Multithreaded should be enabled or not
	System     System // Contains system specific variables needed
	Config     Content
}

type System struct {
	OS        string
	Separator string
}

const (
	UNIX_SEPARATOR     = "/"
	WINDOWS_SEPARATOR  = "\\"
	DEFAULT_CONFIGPATH = "./config.toml"
)

var std = New(DEFAULT_CONFIGPATH, true, setOS(), setConfig())

// Construct

func New(configpath string, mthreads bool, sys System, config Content) *Register {
	return &Register{ConfigPath: configpath, MThreads: mthreads, System: sys, Config: config}
}

func setOS() System {
	var sys System
	sys.OS = runtime.GOOS
	if sys.OS == "windows" {
		sys.Separator = WINDOWS_SEPARATOR
	} else {
		sys.Separator = UNIX_SEPARATOR
	}
	return sys
}

func setConfig() Content {
	var temp Content
	temp.Path = DEFAULT_CONFIGPATH
	temp.LoadTree()
	return temp

}

// Global functions

func HandleFlags() {
	const (
		mthread_desc   = "Enables or Disables Multithreading"
		mthread_defval = true
		mthread_flag   = "t"
	)

	flag.BoolVar(&std.MThreads, mthread_flag, mthread_defval, mthread_desc)
	flag.Parse()
}

func HandleArgs() {

}

// Set methods

func SetConfigPath(path string) {
	std.ConfigPath = path
}

func SetMthreads(state bool) {
	std.MThreads = state

}

// Get methods

func GetConfigPath() string {
	return std.ConfigPath
}

func GetMthreads() bool {
	return std.MThreads
}

func GetSysinfo() System {
	return std.System
}

func GetConfigContent_Struct() Config {
	std.Config.LoadStruct()
	return std.Config.Struct
}

func GetConfigContent_Tree() *toml.Tree {
	std.Config.LoadTree()
	return std.Config.Tree
}

func QueryTree(setting string) string {
	return std.Config.Query(setting)
}
