package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const (
	toolVersion   = "0.2.0"
	capacityLabel = "zpool_capacity_percentage"
	capacityHelp  = "Current zpool capacity level"
	faultedLabel  = "zpool_faulted_providers_count"
	faultedHelp   = "Number of FAULTED/UNAVAIL zpool providers (disks)"
	onlineLabel   = "zpool_online_providers_count"
	onlineHelp    = "Number of ONLINE zpool providers (disks)"
)

type exporter struct {
	zpool *zpool
}

func (e *exporter) export() {
	e.zpool.getStatus()
	e.writeUsage()
	e.writeOnline()
	e.writeFaulted()
}

func (e *exporter) writeUsage() {
	fmt.Printf("# HELP %s %s\n", capacityLabel, capacityHelp)
	fmt.Printf("# TYPE %s gauge\n", capacityLabel)
	fmt.Printf("%s %d\n", capacityLabel, e.zpool.capacity)
}

func (e *exporter) writeFaulted() {
	fmt.Printf("# HELP %s %s\n", faultedLabel, faultedHelp)
	fmt.Printf("# TYPE %s gauge\n", faultedLabel)
	fmt.Printf("%s %d\n", faultedLabel, e.zpool.faulted)
}

func (e *exporter) writeOnline() {
	fmt.Printf("# HELP %s %s\n", onlineLabel, onlineHelp)
	fmt.Printf("# TYPE %s gauge\n", onlineLabel)
	fmt.Printf("%s %d\n", onlineLabel, e.zpool.online)
}

var (
	zfsPool      string
	versionCheck bool
)

func init() {
	const (
		defaultPool  = "tank"
		selectedPool = "what ZFS pool to monitor"
		versionUsage = "display current tool version"
		handleUsage  = "HTTP endpoint to export data on"
	)
	flag.StringVar(&zfsPool, "pool", defaultPool, selectedPool)
	flag.StringVar(&zfsPool, "p", defaultPool, selectedPool+" (shorthand)")
	flag.BoolVar(&versionCheck, "version", false, versionUsage)
	flag.Parse()
}

func main() {
	if versionCheck {
		fmt.Printf("prometheus-zfs v%s (https://github.com/eripa/prometheus-zfs)\n", toolVersion)
		os.Exit(0)
	}
	err := checkExistence(zfsPool)
	if err != nil {
		log.Fatal(err)
	}
	z := zpool{name: zfsPool}
	exporter := exporter{zpool: &z}
	exporter.export()
}
