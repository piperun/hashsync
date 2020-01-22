package netcom

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/piperun/hashsync/settings"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type Connection struct {
	Client client
}

type client struct {
	ssh    *ssh.Client
	sftp   *sftp.Client
	config *ssh.ClientConfig
}

type authdata struct {
	user       map[string]string
	authmethod uint
}

var __config_content = settings.Content{}

// Construct

func InsertConfigContent(cc settings.Content) {
	__config_content = cc
}

// Connect opens sftp & sftp connection
func (connection *Connection) Connect() {
	var err error
	user := getUser(__config_content)
	auth := authdata{
		user: user,
	}

	connection.Client.SetClientConfig(auth)

	// connect
	connection.Client.ssh, err = ssh.Dial("tcp", user["conn_ip"]+":"+user["conn_port"], connection.Client.config)
	if err != nil {
		log.Fatal(err)
	}

	// open an SFTP session over an existing ssh connection.
	connection.Client.sftp, err = sftp.NewClient(connection.Client.ssh)
	if err != nil {
		log.Fatal(err)
	}

}

// Disconnect closes sftp & ssh connection
func (connection *Connection) Disconnect() {
	connection.Client.sftp.Close()
	connection.Client.ssh.Close()
}

func (client *client) SetClientConfig(auth authdata) {

	// get host public key
	//hostKey := getHostKey(remote)

	authmethod := []ssh.AuthMethod{}

	if auth.user["key"] != "" {
		// TODO
	}

	if auth.user["user"] != "" && auth.user["password"] != "" {
		authmethod = append(authmethod, ssh.Password(auth.user["password"]))
	}

	client.config = &ssh.ClientConfig{
		User: auth.user["user"],
		Auth: authmethod,
		//Temporary solution
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		//HostKeyCallback: ssh.FixedHostKey(hostKey),
	}
}

func (client *client) ListDir(path string) []string {
	var arr []string

	w := client.sftp.Walk(path)
	for w.Step() {
		if w.Err() != nil {
			continue
		}
		arr = append(arr, w.Path())
	}
	return arr
}

func (client *client) GetSFTPConnection() *sftp.Client {
	return client.sftp
}

// Local functions

func getUser(content settings.Content) map[string]string {
	var user = map[string]string{
		"user":      "",
		"password":  "",
		"conn_ip":   "",
		"conn_port": "",
	}

	for k, _ := range user {
		user[k] = content.Query("SSH." + k)
	}
	return user
}

func getHostKey(host string) ssh.PublicKey {
	// parse OpenSSH known_hosts file
	// ssh or use ssh-keyscan to get initial key
	file, err := os.Open(filepath.Join("C:", os.Getenv("HOMEPATH"), ".ssh", "known_hosts"))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var hostKey ssh.PublicKey
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		if len(fields) != 3 {
			continue
		}
		if strings.Contains(fields[0], host) {
			var err error
			hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
			if err != nil {
				log.Fatalf("error parsing %q: %v", fields[2], err)
			}
			break
		}
	}

	if hostKey == nil {
		log.Fatalf("no hostkey found for %s", host)
	}

	return hostKey
}
