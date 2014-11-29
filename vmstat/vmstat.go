package vmstat

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
	file, err := os.Open("/proc/vmstat")
	if err != nil {
		// That's weird. Oh well, we'll have to emit an error and return
		// empty work.
		return defs.Response{
			Metrics: make([]defs.Metric, 0),
		}, nil
	}
	defer file.Close()

	metrics := make([]defs.Metric, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		memline := scanner.Text()
		parts := strings.Fields(memline)

		// Each of these lines is:
		// 0 - Name (nr_free_pages)
		// 1 - Value (12345)

		// Make sure we got something that looks correct in terms of fields
		if len(parts) != 2 {
			// Weird. Don't know how to grok this line so spit it out and
			// move on
			log.Printf("Expected 2 fields, got something else: %v", memline)
			continue
		}

		// Make sure we can parse the memory value as a float 64, else
		// we'll skip.
		floatval, fconverr := strconv.ParseFloat(parts[1], 64)
		if fconverr != nil {
			log.Printf("Cannot parse vmstat value '%v' as float64, skipping\n", parts[1])
			continue
		}

		// Line looks good, make the info struct so we can send it back
		info := make(map[string]string)

		// Don't want that pesky :
		name := parts[0]
		info["what"] = name
		info["unit"] = "Page"
		info["target_type"] = "gauge"

		metrics = append(metrics, defs.Metric{
			Timestamp: timestamp,
			Info:      info,
			Value:     floatval,
		})
	}

	if err := scanner.Err(); err != nil {
		return defs.Response{
			Metrics: make([]defs.Metric, 0),
		}, err
	}

	return defs.Response{
		Metrics: metrics,
	}, nil
}
