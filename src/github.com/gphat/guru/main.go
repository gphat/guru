package main

import (
  "fmt"
  "log"
  "os"
  "time"
  "github.com/gphat/guru/defs"
  "github.com/gphat/guru/memory"
)

type HostInfo struct {
  hostname string
}

func main() {

  plugins := map[string]func() defs.Response{
    "memory": memory.GetMetrics,
  }

  ticker := time.NewTicker(time.Millisecond * 1000)
  go func() {
    for t := range ticker.C {
      for plugin_name, f := range plugins {
        hostname, err := os.Hostname()
        if err != nil {
          fmt.Printf("Shit, error: %v\n", err)
        }
        var hi = HostInfo{hostname: hostname}

        fmt.Printf("Hello, from %v\n", hi.hostname)
        fmt.Printf("Running: %v\n", plugin_name)
        resp := f()

        // XXX We don't have the hostname yet. It seems better to add
        // "global" values to the Metric's Info field. Some example:
        //  * server=hostname
        // meta value for agent (guru)
        if(len(resp.Metrics) > 0) {
          fmt.Println(resp.Metrics[0])
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
