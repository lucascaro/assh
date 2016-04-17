package main

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"
)

// RunSSH to multiple hosts, using the same user and key.
// Uses csshx to run.
// TODO: using cssh if linux.
func RunSSH(ips []string, user, sshKey string) {
	cmd := "csshx"
	if runtime.GOOS != "darwin" {
		cmd = "cssh"
	}
	args := []string{}
	fmt.Println("Running ssh with key: ", sshKey)
	for _, ip := range ips {
		args = append(args, user+"@"+ip)
	}
	args = append(args, "--ssh_args", "-i "+sshKey)
	fmt.Println(cmd, strings.Join(args, " "))
	out, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		fmt.Printf("ERROR: %s %s\n", err, out)
		log.Fatal(err)
	}
}
