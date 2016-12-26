package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/gphat/guru/cpu"
	"github.com/gphat/guru/defs"
	"github.com/gphat/guru/disks"
	"github.com/gphat/guru/ifaces"
	"github.com/gphat/guru/loadavg"
	"github.com/gphat/guru/memory"
	"github.com/gphat/guru/vmstat"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Plugin is a Guru plugin
type Plugin interface {
	GetMetrics() (defs.Response, error)
}

func main() {

	var (
		interval = kingpin.Flag("interval", "Interval to collect metrics.").Default("5s").Short('i').Duration()
	)
	kingpin.Parse()

	plugins := map[string]Plugin{
		"cpu":     cpu.NewCPU(),
		"disks":   disks.NewDisks(),
		"loadavg": loadavg.NewLoadAvg(),
		"memory":  memory.NewMemory(),
		"ifaces":  ifaces.NewIfaces(),
		"vmstat":  vmstat.NewVMStat(),
	}

	conn, err := net.Dial("udp", "localhost:8125")
	if err != nil {
		// blah
	}

	// Collect some metadata to use with the metrics
	meta := make(map[string]string)
	meta["agent"] = "guru"

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Printf("Error fetching hostname: %v\n", err)
	}
	fmt.Printf("Hello, from %v\n", hostname)

	ticker := time.NewTicker(*interval)
	go func() {
		for t := range ticker.C {

			for pluginName, p := range plugins {
				fmt.Printf("Running: %v\n", pluginName)
				resp, err := p.GetMetrics()

				// XXX We don't have the hostname yet. It seems better to add
				// "global" values to the Metric's Info field. Some example:
				//  * server=hostname
				//  * guru module?=memory or whatever
				// meta value for agent (guru)
				if err != nil {
					log.Printf("Failed to execute plugin '%v': %v\n", pluginName, err)
				} else if len(resp.Metrics) > 0 {
					for _, met := range resp.Metrics {
						fmt.Fprintf(conn, defs.StringifyMetric(hostname, meta, met))
						log.Println(defs.StringifyMetric(hostname, meta, met))
					}
				} else {
					log.Printf("Plugin '%v' returned 0 metrics.\n", pluginName)
				}
				log.Println("Ticker at", t)
			}
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			fmt.Printf("Got signal %s", sig)
			log.Println("Stopping ticker")
			ticker.Stop()
			os.Exit(0)
		}
	}()

	select {}
}
