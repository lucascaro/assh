// assh -- ssh into all instances in a AWS auto scaling group.
package main

import (
	"flag"
	"fmt"
	"os"
)

type commandOptions struct {
	SSHKey  string
	User    string
	ASGName string
}

func main() {
	printIntro()
	options := setupFlags()

	ips := GetIPAddresses(options.ASGName)
	fmt.Println("IP Addresses:", ips)

	RunSSH(ips, options.User, options.SSHKey)
}

func asshUsage() {
	fmt.Println("Usage:")
	fmt.Println("\t assh [flags] <ASG Name>")
	fmt.Println("Flags:")
	flag.PrintDefaults()
}

func printIntro() {
	fmt.Println("assh")
	fmt.Print("\t ssh into auto scaling group\n\n")
}

// set up and parse command line flags
func setupFlags() commandOptions {
	// Usage text
	flag.Usage = asshUsage

	// Flags
	sshKey := flag.String("i", "${HOME}/.ssh/id-rsa.pem", "your ssh key")
	user := flag.String("u", "ec2-user", "ssh user name")

	flag.Parse()

	// Required arguments
	if len(flag.Args()) < 1 {
		fmt.Println("Auto Scaling Group name is required")
		os.Exit(1)
	}

	// Return a nice struct with all flags and arguments
	return commandOptions{
		ASGName: flag.Arg(0),
		SSHKey:  *sshKey,
		User:    *user,
	}
}
