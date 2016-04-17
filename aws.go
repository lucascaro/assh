package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/lucascaro/assh/filecache"
)

var fc = filecache.New("${HOME}/.asshcache")

// GetIPAddresses for an auto scaling group.
func GetIPAddresses(asg string) []string {
	fmt.Println("fetching IPs for asg: ", asg)
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

	if string(out) == "None\n" {
		log.Fatalf("Auto Scaling Group '%s' not found", asg)
	}

	fmt.Printf("Instances in ASG: %s\n", out)

	// Split lines and parse for instance ids
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 3 {
			log.Fatalf("Auto Scaling Group '%s' not found", asg)
		}
		if len(fields) > 0 {
			ids = append(ids, fields[2])
		}
	}
	return ids
}

func getInstanceIPs(instanceIds []string) []string {
	ips, missingIds := getCachedIPs(instanceIds)

	if len(missingIds) > 0 {
		fmt.Println("Get ips for", missingIds)
		ips = append(ips, getUncachedInstanceIPs(missingIds)...)

		for i, ip := range ips {
			fc.Set(missingIds[i], ip)
		}
		fc.Save()
	}

	return ips
}

func getUncachedInstanceIPs(ids []string) []string {
	cmd := "aws"
	args := []string{
		"ec2",
		"describe-instances",
		"--instance-ids",
	}
	args = append(args, ids...)
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

	return strings.Fields(string(out))
}
func getCachedIPs(instanceIds []string) (ips, missingIds []string) {
	for i, id := range instanceIds {
		if cached, ok := fc.Get(id); ok == nil {
			fmt.Println("Cached: ", id, cached)
			ips = append(ips, cached.Value)
		} else {
			missingIds = append(missingIds, instanceIds[i])
		}
	}
	return ips, missingIds
}
