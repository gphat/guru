package defs

import (
  "bytes"
  "fmt"
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

func StringifyMetric(metric Metric) string {

  var buffer bytes.Buffer
  for k, v := range metric.Info {
    buffer.WriteString(fmt.Sprintf("%v=%v ", k, v))
  }
  // Add another space as we're separated from meta by a space
  buffer.WriteString(" ")

  buffer.WriteString(fmt.Sprintf("%v %v", metric.Value, metric.Timestamp))

  return buffer.String()
}
