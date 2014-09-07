package parser

import (
  "errors"
  "fmt"
  "strconv"
)

func ParseFloats(toConvert []string, start int, end int) ([]float64, error) {

  floats := make([]float64, (end - start) + 1)

  for i := start; i <= end; i++ {
    floatval, fconverr := strconv.ParseFloat(toConvert[i], 64)
    if fconverr != nil {
      return nil, errors.New(fmt.Sprintf("Cannot parse value '%v' as float64, skipping\n", toConvert[i]))
      continue
    }
    floats = append(floats, floatval)
  }

  return floats, nil
}
