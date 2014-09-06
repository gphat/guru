package system

import (
  "github.com/gphat/guru/defs"
)

// Sqrt returns an approximation to the square root of x.
func GetMetrics() defs.Response {
  info := make(map[string]string)

  info["hostname"] = "poop"

  metric := defs.Metric{
    info, 1.2345,
    []string{"ass"},
  }

  return defs.Response{
    Metrics: []defs.Metric{
      metric,
    },
  }
}
