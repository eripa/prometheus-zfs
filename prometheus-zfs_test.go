package main

import (
	"fmt"
	"testing"
)

func TestExporter(t *testing.T) {
	zpool := zpool{
		name:     "zones",
		capacity: 65,
		healthy:  true,
		status:   "ONLINE",
		online:   8,
		faulted:  0,
	}
	exporter := exporter{zpool: &zpool}
	res := exporter.getCapacityMetric()
	expected := fmt.Sprintf("# HELP zpool_capacity_percentage Current zpool capacity level\n# TYPE zpool_capacity_percentage gauge\nzpool_capacity_percentage 65\n")
	if res != expected {
		t.Fatalf("Incorrect capacity metrics.\nExpected:\n%sGot:\n%s\n", expected, res)
	}
	res = exporter.getOnlineMetric()
	expected = fmt.Sprintf("# HELP zpool_online_providers_count Number of ONLINE zpool providers (disks)\n# TYPE zpool_online_providers_count gauge\nzpool_online_providers_count 8\n")
	if res != expected {
		t.Fatalf("Incorrect online metrics.\nExpected:\n%sGot:\n%s\n", expected, res)
	}
	res = exporter.getFaultedMetric()
	expected = fmt.Sprintf("# HELP zpool_faulted_providers_count Number of FAULTED/UNAVAIL zpool providers (disks)\n# TYPE zpool_faulted_providers_count gauge\nzpool_faulted_providers_count 0\n")
	if res != expected {
		t.Fatalf("Incorrect faulted metrics.\nExpected:\n%sGot:\n%s\n", expected, res)
	}
}
