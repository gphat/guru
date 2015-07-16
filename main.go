package main

import (
	"fmt"
	"github.com/gphat/guru/cpu"
	"github.com/gphat/guru/defs"
	"github.com/gphat/guru/diskstats"
	"github.com/gphat/guru/loadavg"
	"github.com/gphat/guru/memory"
	"github.com/gphat/guru/netstats"
	"github.com/gphat/guru/vmstat"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

func main() {

	var (
		interval = kingpin.Flag("interval", "Interval to collect metrics.").Default("5s").Short('t').Duration()
	)
	kingpin.Parse()

	plugins := map[string]func() (defs.Response, error){
		"cpu":       cpu.GetMetrics,
		"diskstats": diskstats.GetMetrics,
		"loadavg":   loadavg.GetMetrics,
		"memory":    memory.GetMetrics,
		"netstats":  netstats.GetMetrics,
		"vmstat":    vmstat.GetMetrics,
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

			for plugin_name, f := range plugins {
				fmt.Printf("Running: %v\n", plugin_name)
				resp, err := f()

				// XXX We don't have the hostname yet. It seems better to add
				// "global" values to the Metric's Info field. Some example:
				//  * server=hostname
				//  * guru module?=memory or whatever
				// meta value for agent (guru)
				if err != nil {
					log.Printf("Failed to execute plugin '%v': %v\n", plugin_name, err)
				} else if len(resp.Metrics) > 0 {
					for _, met := range resp.Metrics {
						fmt.Fprintf(conn, defs.StringifyMetric(hostname, meta, met))
						log.Println(defs.StringifyMetric(hostname, meta, met))
					}
				} else {
					log.Printf("Plugin '%v' returned 0 metrics.\n", plugin_name)
				}
				log.Println("Ticker at", t)
			}
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			fmt.Println("Got signal %s", sig)
			log.Println("Stopping ticker")
			ticker.Stop()
			os.Exit(0)
		}
	}()

	select {}
}
