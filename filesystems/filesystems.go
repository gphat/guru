package filesystems

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/gphat/guru/defs"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func GetSizeMetric(ts time.Time, position int, mt string, parts []string) (defs.Metric, error) {
	info := make(map[string]string)
	info["target_type"] = "gauge"
	info["unit"] = "KiB"
	info["what"] = "disk_space"
	info["type"] = mt
	info["mountpoint"] = parts[6]
	info["device"] = parts[0]
	size, converr := strconv.ParseFloat(parts[position], 64)
	if converr != nil {
		return defs.Metric{}, errors.New(fmt.Sprintf("Cannot parse fs value '%v' as float64", parts[position]))
	}

	return defs.Metric{
		Timestamp: ts,
		Info:      info,
		Value:     size,
	}, nil
}

func ShouldSkip(fstype string) bool {
	switch fstype {
	case "tmpfs", "devtmpfs":
		return true
	default:
		return false
	}
}

func GetMetrics() (defs.Response, error) {

	timestamp := time.Now()
	out, err := exec.Command("df", "-T").Output()
	if err != nil {
		// That's weird. Oh well, we'll have to emit an error and return
		// empty work.
		return defs.Response{
			Metrics: make([]defs.Metric, 0),
		}, err
	}

	metrics := make([]defs.Metric, 0)

	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		dfline := scanner.Text()
		parts := strings.Fields(dfline)

		// Each of these lines is:
		// 0 - Filesystem
		// 1 - FS Type
		// 2 - 1K-blocks
		// 3 - Used
		// 4 - Available
		// 5 - Use%
		// 6 - Mountpoint

		// Make sure we got something that looks correct in terms of fields
		if len(parts) != 7 {
			// Weird. Don't know how to grok this line so spit it out and
			// move on
			log.Printf("Expected 7 fields, got something else: %v", dfline)
			continue
		}

		// We may not want to bother collecting…
		if ShouldSkip(parts[1]) {
			continue
		}

		tmet, tmerr := GetSizeMetric(timestamp, 3, "total", parts)
		if tmerr != nil {
			log.Printf("Error getting size metric for fs '%s': %s", parts[0], tmerr)
			continue
		} else {
			metrics = append(metrics, tmet)
		}

		umet, umerr := GetSizeMetric(timestamp, 4, "used", parts)
		if umerr != nil {
			log.Printf("Error getting size metric for fs '%s': %s", parts[0], umerr)
			continue
		} else {
			metrics = append(metrics, umet)
		}

		amet, amerr := GetSizeMetric(timestamp, 5, "available", parts)
		if amerr != nil {
			log.Printf("Error getting size metric for fs '%s': %s", parts[0], amerr)
			continue
		} else {
			metrics = append(metrics, amet)
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
