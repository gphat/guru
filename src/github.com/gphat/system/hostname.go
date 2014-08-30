package system

 import (
  "github.com/gphat/defs"
)

// Sqrt returns an approximation to the square root of x.
func GetMetrics() defs.Response {
  info := make(map[string]string)

  info["hostname"] = "poop"

  resp := defs.Response{Metrics: []defs.Metric{defs.Metric{info, 1.2345}}}

  return resp
}
