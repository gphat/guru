package defs

import (
  "time"
)

type Metric struct {
  Timestamp time.Time
  Info      map[string]string
  Value     float64
}

type Response struct {
  Metrics []Metric
}
