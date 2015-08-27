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
	fmt.Printf(e.getCapacityMetric())
	fmt.Printf(e.getOnlineMetric())
	fmt.Printf(e.getFaultedMetric())
}

func (e *exporter) getCapacityMetric() string {
	return fmt.Sprintf("# HELP %s %s\n# TYPE %s gauge\n%s %d\n", capacityLabel, capacityHelp, capacityLabel, capacityLabel, e.zpool.capacity)
}

func (e *exporter) getFaultedMetric() string {
	return fmt.Sprintf("# HELP %s %s\n# TYPE %s gauge\n%s %d\n", faultedLabel, faultedHelp, faultedLabel, faultedLabel, e.zpool.faulted)
}

func (e *exporter) getOnlineMetric() string {
	return fmt.Sprintf("# HELP %s %s\n# TYPE %s gauge\n%s %d\n", onlineLabel, onlineHelp, onlineLabel, onlineLabel, e.zpool.online)
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
