package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	toolVersion = "0.1.0"
)

// Exporter collects zpool stats from the given zpool and exports them using
// the prometheus metrics package.
type Exporter struct {
	mutex sync.RWMutex

	poolUsage, providersFaulted, providersOnline prometheus.Gauge
	hdFailures                                   prometheus.Counter
	zpool                                        *zpool
}

// NewExporter returns an initialized Exporter.
func NewExporter(zp *zpool) *Exporter {
	// Init and return our exporter.
	return &Exporter{
		zpool: zp,
		poolUsage: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "zpool_capacity_percentage",
			Help: "Current zpool capacity level",
		}),
		providersOnline: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "zpool_online_providers_count",
			Help: "Number of ONLINE zpool providers (disks)",
		}),
		providersFaulted: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "zpool_faulted_providers_count",
			Help: "Number of FAULTED/UNAVAIL zpool providers (disks)",
		}),
		hdFailures: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "hd_errors_total",
			Help: "Number of hard-disk errors.",
		}),
	}
}

// Describe describes all the metrics ever exported by the zpool exporter. It
// implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.poolUsage.Desc()
	ch <- e.providersOnline.Desc()
	ch <- e.providersFaulted.Desc()
	ch <- e.hdFailures.Desc()
}

// Collect fetches the stats from configured ZFS pool and delivers them
// as Prometheus metrics. It implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.zpool.getStatus()

	e.hdFailures.Inc()
	e.poolUsage.Set(float64(e.zpool.capacity))
	e.providersOnline.Set(float64(e.zpool.disks))
	e.providersFaulted.Set(float64(e.zpool.faulted))

	e.mutex.Lock() // To protect metrics from concurrent collects.
	defer e.mutex.Unlock()

	ch <- e.poolUsage
	ch <- e.providersOnline
	ch <- e.providersFaulted
	ch <- e.hdFailures
}

var (
	zfsPool       string
	listenPort    string
	metricsHandle string
	versionCheck  bool
)

func init() {
	const (
		defaultPool   = "tank"
		selectedPool  = "what ZFS pool to monitor"
		versionUsage  = "display current tool version"
		defaultPort   = "8080"
		portUsage     = "Port to listen on"
		defaultHandle = "metrics"
		handleUsage   = "HTTP endpoint to export data on"
	)
	flag.StringVar(&zfsPool, "pool", defaultPool, selectedPool)
	flag.StringVar(&zfsPool, "p", defaultPool, selectedPool+" (shorthand)")
	flag.StringVar(&listenPort, "port", defaultPort, portUsage)
	flag.StringVar(&metricsHandle, "endpoint", defaultHandle, handleUsage)
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

	exporter := NewExporter(&z)
	prometheus.MustRegister(exporter)

	fmt.Printf("Starting zpool metrics exporter on :%s/%s\n", listenPort, metricsHandle)
	http.Handle("/"+metricsHandle, prometheus.Handler())
	http.ListenAndServe(":"+listenPort, nil)

}
