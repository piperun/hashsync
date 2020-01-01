package netcom

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/piperun/hashsync/config"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type Connection struct {
	Client client
}

type client struct {
	ssh  *ssh.Client
	sftp *sftp.Client
}

func SFTPConnect(content config.Content) {
	user := getUser(content)

	// get host public key
	//hostKey := getHostKey(remote)

	config := &ssh.ClientConfig{
		User: user["user"],
		Auth: []ssh.AuthMethod{
			ssh.Password(user["password"]),
		},
		//Temporary solution
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		//HostKeyCallback: ssh.FixedHostKey(hostKey),
	}

	// connect
	conn, err := ssh.Dial("tcp", user["conn_ip"]+":"+user["conn_port"], config)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// open an SFTP session over an existing ssh connection.
	sftp, err := sftp.NewClient(conn)
	if err != nil {
		log.Fatal(err)
	}
	defer sftp.Close()

	// walk a directory
	w := sftp.Walk("/etc")
	for w.Step() {
		if w.Err() != nil {
			continue
		}
		log.Println(w.Path())
	}
}

// Local functions

func getUser(content config.Content) map[string]string {
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
