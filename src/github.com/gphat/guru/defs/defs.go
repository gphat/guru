package defs

type Metric struct {
  Info  map[string]string
  Value float64
}

type Response struct {
  Metrics []Metric
}
