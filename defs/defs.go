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

func StringifyMetric(host string, meta map[string]string, metric Metric) string {

  var buffer bytes.Buffer
  for k, v := range metric.Info {
    buffer.WriteString(fmt.Sprintf("%v=%v ", k, v))
  }
  buffer.WriteString(fmt.Sprintf("host=%v ", host))
  // Add another space as we're separated from meta by a space
  buffer.WriteString(" ")

  // Now the meta tags
  for k, v := range meta {
    buffer.WriteString(fmt.Sprintf("%v=%v ", k, v))
  }

  buffer.WriteString(fmt.Sprintf("%v %v", metric.Value, metric.Timestamp.Unix()))

  return buffer.String()
}
