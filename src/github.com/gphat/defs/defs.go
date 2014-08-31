package defs

type Metric struct {
  Info map[string]string
  Value float64
  Tags []string
}

type Response struct {
  Metrics []Metric
}

