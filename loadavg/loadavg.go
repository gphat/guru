package loadavg

import (
	"bufio"
	"github.com/gphat/guru/defs"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func GetMetrics() (defs.Response, error) {

	timestamp := time.Now()
	file, err := os.Open("/proc/loadavg")
	if err != nil {
		// That's weird. Oh well, we'll have to emit an error and return
		// empty work.
		return defs.Response{
			Metrics: make([]defs.Metric, 0),
		}, err
	}
	defer file.Close()

	metrics := make([]defs.Metric, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		memline := scanner.Text()
		parts := strings.Fields(memline)

		// Each of these lines is:
		// 0 - 1 min
		// 1 - 5 min
		// 2 - 15 min
		// 3 - executing/scheduling
		// 4 - last pid

		// Make sure we got something that looks correct in terms of fields
		if len(parts) != 5 {
			// Weird. Don't know how to grok this line so spit it out and
			// move on
			log.Printf("Expected 5 fields, got something else: %v", memline)
			continue
		}

		for index, element := range parts {
			// Line looks good, make the info struct so we can send it back
			switch index {
			case 4:
				// Don't care about the last pid
				continue
			case 3:
				schedparts := strings.Split(element, "/")

				runninginfo := make(map[string]string)
				runninginfo["unit"] = "Thread"
				runninginfo["what"] = "scheduled"
				runninginfo["target_type"] = "gauge"

				runningval, rconverr := strconv.ParseFloat(schedparts[0], 64)
				if rconverr != nil {
					log.Printf("Cannot parse loadavg scheduler value '%v' as float64, skipping\n", schedparts[0])
					continue
				}
				// Append the running threads number
				metrics = append(metrics, defs.Metric{
					Timestamp: timestamp,
					Info:      runninginfo,
					Value:     runningval,
				})

				totalinfo := make(map[string]string)
				totalinfo["unit"] = "Thread"
				totalinfo["what"] = "schedulable"
				totalinfo["target_type"] = "gauge"

				totalval, tconverr := strconv.ParseFloat(schedparts[1], 64)
				if tconverr != nil {
					log.Printf("Cannot parse loadavg scheduler value '%v' as float64, skipping\n", schedparts[1])
					continue
				}
				// Append the running threads number
				metrics = append(metrics, defs.Metric{
					Timestamp: timestamp,
					Info:      totalinfo,
					Value:     totalval,
				})

			default:
				info := make(map[string]string)
				info["unit"] = "Load"
				info["target_type"] = "gauge"
				// Make sure we can parse the memory value as a float 64, else
				// we'll skip.
				floatval, fconverr := strconv.ParseFloat(element, 64)
				if fconverr != nil {
					log.Printf("Cannot parse loadavg value '%v' as float64, skipping\n", element)
					continue
				}

				switch index {
				case 0:
					info["what"] = "load1"
				case 1:
					info["what"] = "load5"
				case 2:
					info["what"] = "load15"
				}

				metrics = append(metrics, defs.Metric{
					Timestamp: timestamp,
					Info:      info,
					Value:     floatval,
				})
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return defs.Response{
			Metrics: make([]defs.Metric, 0),
		}, err
	}

	return defs.Response{
		Metrics: metrics,
	}, err
}
