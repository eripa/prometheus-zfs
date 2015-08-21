package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const (
	toolVersion = "0.1.0"
)

var zfsPool string
var versionCheck bool

func init() {
	const (
		defaultPool  = "tank"
		poolUsage    = "what ZFS pool to monitor"
		versionUsage = "display current tool version"
	)
	flag.StringVar(&zfsPool, "pool", defaultPool, poolUsage)
	flag.StringVar(&zfsPool, "p", defaultPool, poolUsage+" (shorthand)")
	flag.BoolVar(&versionCheck, "version", false, versionUsage)
	flag.Parse()
}

func main() {
	if versionCheck {
		fmt.Printf("prometheus-zfs v%s (https://github.com/eripa/prometheus-zfs)\n", toolVersion)
		os.Exit(0)
	}
	err := checkExistance(zfsPool)
	if err != nil {
		log.Fatal(err)
	}
	z := zpool{name: zfsPool}
	z.getStatus()
	fmt.Printf("%+v\n", z)
}
