package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/lucascaro/assh/filecache"
)

var fc = filecache.New("${HOME}/.asshcache")

func main() {
	fmt.Println("assh")
	fmt.Print("\n\t ssh into auto scaling group\n\n")
	flag.Usage = func() {
		fmt.Printf("Usage:\n")
		fmt.Println("\t assh [flags] <ASG Name>")
		fmt.Println("\nFlags:")
		flag.PrintDefaults()
	}
	sshKey := flag.String("i", "${HOME}/.ssh/fpc-rcloud-ecs.pem", "your ssh key")
	user := flag.String("u", "ec2-user", "ssh user name")
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Println("Auto Scaling Group name is required")
		os.Exit(1)
	}
	asgName := flag.Arg(0)
	fmt.Println("fetching IPs for asg: ", asgName)

	ips := getIPAddresses(asgName)
	fmt.Println("IP Addresses:", ips)
	fmt.Println(*sshKey)

	runSSH(ips, *user, *sshKey)
}

func runSSH(ips []string, user, sshKey string) {
	cmd := "csshx"
	args := []string{}
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
func getIPAddresses(asg string) []string {
	ips := []string{}
	ids := getInstanceIds(asg)
	fmt.Println("Instance IDs:", ids)
	ips = getInstanceIPs(ids)
	return ips
}

func getInstanceIds(asg string) []string {
	ids := []string{}
	cmd := "aws"
	args := []string{"autoscaling",
		"describe-auto-scaling-groups",
		"--auto-scaling-group-names",
		asg,
		"--query",
		"AutoScalingGroups[0].Instances[]",
		"--output",
		"text",
	}
	out, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		fmt.Printf("ERROR: %s %s\n", err, out)
		log.Fatal(err)
	}

	// Split lines
	lines := strings.Split(string(out), "\n")
	fmt.Printf("ASG DESCRIBE: %s\n", out)

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) > 0 {
			ids = append(ids, fields[2])
		}
	}
	return ids
}

func getInstanceIPs(instanceIds []string) []string {
	var ips []string
	var missingIds []string
	for i, id := range instanceIds {
		if cached, ok := fc.Get(id); ok == nil {
			fmt.Println("Cached: ", id, cached)
			ips = append(ips, cached.Value)
		} else {
			missingIds = append(missingIds, instanceIds[i])
		}
	}
	instanceIds = missingIds
	fmt.Println("Get ips for", instanceIds)
	if len(instanceIds) > 0 {
		cmd := "aws"
		args := []string{
			"ec2",
			"describe-instances",
			"--instance-ids",
		}
		args = append(args, instanceIds...)
		args = append(args, "--query",
			"Reservations[].Instances[].PrivateIpAddress",
			"--output",
			"text",
		)
		// fmt.Println(args)
		out, err := exec.Command(cmd, args...).CombinedOutput()
		if err != nil {
			fmt.Printf("ERROR: %s %s\n", err, out)
			log.Fatal(err)
		}

		ips = strings.Fields(string(out))

		for i, ip := range ips {
			fc.Set(instanceIds[i], ip)
		}
		fc.Save()
		// fmt.Printf("INSTANCE DESCRIBE: %s\n", ips)
	}

	return ips
}
