package memory

import (
  "bufio"
  "fmt"
  "github.com/gphat/guru/defs"
  "log"
  "os"
)

func GetMetrics() defs.Response {

  file, err := os.Open("/proc/meminfo")
  if err != nil {
    // That's weird. Oh well, we'll have to emit an error and return
    // empty work.
    log.Fatal(err)
    return defs.Response{
      Metrics: make([]defs.Metric, 0),
    }
  }
  defer file.Close()

  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
      fmt.Println(scanner.Text())
  }

  if err := scanner.Err(); err != nil {
      log.Fatal(err)
  }

  return defs.Response{
    Metrics: make([]defs.Metric, 0),
  }
}
