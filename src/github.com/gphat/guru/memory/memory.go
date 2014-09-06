package memory

import (
  "bufio"
  "github.com/gphat/guru/defs"
  "log"
  "os"
  "strconv"
  "strings"
  "time"
)

func GetMetrics() defs.Response {

  timestamp := time.Now()
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

  metrics := make([]defs.Metric, 0)

  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
      memline := scanner.Text()
      parts := strings.Fields(memline)

      // Each of these lines is:
      // 0 - Name (MemTotal)
      // 1 - Value (12345)
      // 2 - Unit (kB)

      // Make sure we got something that looks correct in terms of fields
      if(len(parts) > 3 || len(parts) < 2) {
        // Weird. Don't know how to grok this line so spit it out and
        // move on
        log.Printf("Expected 2 or 3 fields, got something else: %v", memline)
        continue
      }

      // Make sure we can parse the memory value as a float 64, else
      // we'll skip.
      floatval, fconverr := strconv.ParseFloat(parts[1], 64)
      if fconverr != nil {
        log.Printf("Cannot parse memory value '%v' as float64, skipping\n", parts[1])
        continue
      }

      // Line looks good, make the info struct so we can send it back
      info := make(map[string]string)
      if(len(parts) == 3) {
        // If we have 3 then the last one is the unit
        info["unit"] = parts[2]
      } else {
        // If not then the unit is the # pages
        info["unit"] = "Page"
      }

      // Don't want that pesky :
      name := strings.Replace(parts[0], ":", "", 1)
      info["what"] = name
      info["target_type"] = "gauge"

      metrics = append(metrics, defs.Metric{
        Timestamp:  timestamp,
        Info:       info,
        Value:      floatval,
      })
  }

  if err := scanner.Err(); err != nil {
      log.Fatal(err)
  }

  return defs.Response{
    Metrics: metrics,
  }
}
