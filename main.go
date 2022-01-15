package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

var (
	hostname = flag.String("hostIP", "192.168.17.134", "Host IP address")
	port     = "22"
	username = "mylogin"
	password = "P@sswor0"
)

func main() {
	flag.Parse()

	//put any commands you want to be executed on remote hosts
	commands := []string{
		"sh mac add | inc 7563",
	}
	//this is for hosts with less privileges
	//enable + configuring mode password
	enableCli := []string{
		"en",
		"cisco",
	}

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		// Non-production only
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connect to host":"+port, config)
	client, err := ssh.Dial("tcp", *hostname+":"+port, config)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Create session
	sess, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer sess.Close()

	// StdinPipe for commands
	stdin, err := sess.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	// Uncomment to store output in variable
	// var b bytes.Buffer
	// sess.Stdout = &b
	//sess.Stderr = &b

	// Enable system stdout ==============
	// Comment these if you uncomment to store in variable
	sess.Stdout = os.Stdout
	sess.Stderr = os.Stderr

	// Start remote shell
	err = sess.Shell()
	if err != nil {
		log.Fatal(err)
	}
	for _, cmd := range enableCli {
		_, err = fmt.Fprintf(stdin, "%s\n", cmd)
		if err != nil {
			log.Fatal(err)
		}
	}
	for _, cmd := range commands {
		_, err = fmt.Fprintf(stdin, "%s\n", cmd)
		if err != nil {
			log.Fatal(err)
		}
	}
	// Wait for sess to finish
	err = sess.Wait()
	if err != nil {
		log.Fatal(err)
	}

	// Uncomment to store in variable
	fmt.Println(sess.Stdout)
	// out := b.String()
	// if len(out) > 0 {
	// 	fmt.Println(out)
	// } else {
	// 	fmt.Println("No output")
	// }

}
