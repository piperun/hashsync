package registry

import (
	"flag"
)

// Register is the software's internal configuration register, this is to make it easier to allow send data in-between functions that need opinated outside data.
// Note: SSH/Paths configuration is excempt since the data is too opinated.
type Register struct {
	ConfigPath string // Where the config path shall be
	MThreads   bool   // If Multithreaded should be enabled or not
}

var std = New("./config.toml", true)

// Construct

func New(configpath string, mthreads bool) *Register {
	return &Register{ConfigPath: configpath, MThreads: mthreads}
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

func SetConfigPath(path string) {
	std.ConfigPath = path
}

func SetMthreads(state bool) {
	std.MThreads = state

}

func GetConfigPath() string {
	return std.ConfigPath
}

func GetMthreads() bool {
	return std.MThreads
}
