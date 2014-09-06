package main

import (
  "fmt"
  "log"
  "os"
  "time"
  "github.com/gphat/guru/defs"
  "github.com/gphat/guru/diskstats"
  "github.com/gphat/guru/memory"
  "github.com/gphat/guru/netstats"
  "github.com/gphat/guru/vmstat"
)

type HostInfo struct {
  hostname string
}

func main() {

  plugins := map[string]func() defs.Response{
    "diskstats": diskstats.GetMetrics,
    "memory": memory.GetMetrics,
    "netstats": netstats.GetMetrics,
    "vmstat": vmstat.GetMetrics,
  }

  ticker := time.NewTicker(time.Millisecond * 1000)
  go func() {
    for t := range ticker.C {
      for plugin_name, f := range plugins {
        hostname, err := os.Hostname()
        if err != nil {
          fmt.Printf("Error fetching hostname: %v\n", err)
        }
        var hi = HostInfo{hostname: hostname}

        fmt.Printf("Hello, from %v\n", hi.hostname)
        fmt.Printf("Running: %v\n", plugin_name)
        resp := f()

        // XXX We don't have the hostname yet. It seems better to add
        // "global" values to the Metric's Info field. Some example:
        //  * server=hostname
        //  * guru module?=memory or whatever
        // meta value for agent (guru)
        if(len(resp.Metrics) > 0) {
          for _,met := range resp.Metrics {
            log.Println(met)
          }
        } else {
          log.Printf("Plugin '%v' returned 0 metrics.\n", plugin_name)
        }
        fmt.Println("Ticker at", t)
      }
    }
  }()

  time.Sleep(time.Millisecond * 5000)
  ticker.Stop()
  fmt.Println("Ticker stopped")
}
