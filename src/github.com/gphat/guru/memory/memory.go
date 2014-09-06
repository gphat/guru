package memory

import (
  "bufio"
  "github.com/gphat/guru/defs"
  "log"
  "os"
  "strings"
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
      memline := scanner.Text()
      parts := strings.Fields(memline)
      log.Printf("%v\n", memline)
      log.Printf("\t%v\n", len(parts))
  }

  if err := scanner.Err(); err != nil {
      log.Fatal(err)
  }

  return defs.Response{
    Metrics: make([]defs.Metric, 0),
  }
}
