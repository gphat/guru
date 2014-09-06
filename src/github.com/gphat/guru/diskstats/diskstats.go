package diskstats

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
  file, err := os.Open("/proc/diskstats")
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
    //  1 - major number
    //  2 - minor mumber
    //  3 - device name
    //  4 - reads completed successfully
    //  5 - reads merged
    //  6 - sectors read
    //  7 - time spent reading (ms)
    //  8 - writes completed
    //  9 - writes merged
    // 10 - sectors written
    // 11 - time spent writing (ms)
    // 12 - I/Os currently in progress
    // 13 - time spent doing I/Os (ms)
    // 14 - weighted time spent doing I/Os (ms)

    // Make sure we got something that looks correct in terms of fields
    if(len(parts) != 14) {
      // Weird. Don't know how to grok this line so spit it out and
      // move on
      log.Printf("Expected 14 fields, got something else: %v", memline)
      continue
    }

    // Line looks good, loop over the fields
    for i := 3; i < 14; i++ {
      info := make(map[string]string)

      if(strings.HasPrefix(parts[2], "ram") || strings.HasPrefix(parts[2], "loop")) {
        // We're going to skip ramdisks and loop devices
        continue
      }

      // Make sure we can parse the value as a float 64, else
      // we'll skip.
      floatval, fconverr := strconv.ParseFloat(parts[i], 64)
      if fconverr != nil {
        log.Printf("Cannot parse diskstats value '%v' as float64, skipping\n", parts[1])
        continue
      }

      info["device"] = parts[2] // Device name

      // Switch on the index since each line is a different
      switch i {
        case 3:
          info["what"]        = "read_success"
          info["unit"]        = "Req"
          info["direction"]   = "out"
          info["target_type"] = "count"
        case 4:
          info["what"]        = "read_merged"
          info["unit"]        = "Req"
          info["direction"]   = "out"
          info["target_type"] = "count"
        case 5:
          info["what"]        = "read_sectors"
          info["unit"]        = "Sector"
          info["direction"]   = "out"
          info["target_type"] = "count"
        case 6:
          info["what"]        = "reading"
          info["unit"]        = "ms"
          info["direction"]   = "out"
          info["target_type"] = "count"
        case 7:
          info["unit"]        = "Req"
          info["direction"]   = "in"
          info["target_type"] = "count"
        case 8:
          info["what"]        = "write_merged"
          info["unit"]        = "Req"
          info["direction"]   = "in"
          info["target_type"] = "count"
        case 9:
          info["what"]        = "write_sectors"
          info["unit"]        = "Req"
          info["direction"]   = "in"
          info["target_type"] = "count"
        case 10:
          info["what"]        = "writing"
          info["unit"]        = "ms"
          info["direction"]   = "in"
          info["target_type"] = "count"
        case 11:
          info["what"]        = "io_in_progress"
          info["unit"]        = "Req"
          info["target_type"] = "gauge"
        case 12:
          info["what"]        = "io"
          info["unit"]        = "ms"
          info["target_type"] = "count"
        case 13:
          info["what"]        = "io_weighted"
          info["unit"]        = "ms"
          info["target_type"] = "count"
      }

      metrics = append(metrics, defs.Metric{
        Timestamp:  timestamp,
        Info:       info,
        Value:      floatval,
      })
    }
  }

  if err := scanner.Err(); err != nil {
      log.Fatal(err)
  }

  return defs.Response{
    Metrics: metrics,
  }
}
