package cpu

import (
  "bufio"
  "github.com/gphat/guru/defs"
  "github.com/gphat/guru/parser"
  "log"
  "os"
  "regexp"
  "strconv"
  "strings"
  "time"
)

func GetMetrics() defs.Response {

  var cpuLine = regexp.MustCompile(`^cpu[0-9]+`)

  timestamp := time.Now()
  file, err := os.Open("/proc/stat")
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

      // We'll use this to discern what we're looking at
      sigil := parts[0]

      // We'll go ahead and convert the first value since everything uses this
      floatval, fconverr := strconv.ParseFloat(parts[1], 64)
      if fconverr != nil {
        log.Printf("Cannot parse stat value '%v' as float64, skipping\n", parts[1])
        continue
      }

      // Line looks good, make the info struct so we can send it back
      info := make(map[string]string)

      // These are common
      info["target_type"] = "counter"

      // We can't do an exact switch because we need to do some matching
      // so we'll use a boolean
      switch {
        case sigil == "ctxt":
          info["what"] = "ctxt"
          info["unit"] = "Event"
        case cpuLine.MatchString(sigil):
          // This one needs to be first because later we'll check for cpu
          info["device"] = sigil // Put the cpu as the device
          continue
        case strings.HasPrefix(sigil, "cpu"):
          // This catches the total CPU line
          _, ferr := parser.ParseFloats(parts, 2, 10)
          if ferr != nil {
            log.Printf("Cannot parse stat value '%v' as float64, skipping\n", parts[1])
            continue
          }
          continue
        default:
          // Ignore the other stuff
          continue
      }

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

