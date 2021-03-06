package ifaces

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gphat/guru/defs"
)

type Ifaces struct{}

func NewIfaces() *Ifaces {
	return &Ifaces{}
}

func (p *Ifaces) GetMetrics() (defs.Response, error) {

	timestamp := time.Now()
	file, err := os.Open("/proc/net/dev")
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

		if strings.HasPrefix(memline, "Inter") || strings.HasPrefix(memline, " face") {
			// Skip headers
			continue
		}

		parts := strings.Fields(memline)

		// Each of these lines is:
		// 0 - interface name
		// Receive
		// 1 - bytes
		// 2 - packets
		// 3 - errs
		// 4 - drop
		// 5 - fifo
		// 6 - frame
		// 7 - compressed
		// 8 - multicast
		// Transmit
		// 9 - bytes
		// 10 - packets
		// 11 - errs
		// 12 - drop
		// 13 - fifo
		// 14 - colls
		// 14 - carrier
		// 15 - compressed

		iface := strings.Replace(parts[0], ":", "", 1)

		if strings.HasPrefix(iface, "lo") {
			// We're going to skip loop devices
			continue
		}

		// Make sure we got something that looks correct in terms of fields
		if len(parts) != 17 {
			// Weird. Don't know how to grok this line so spit it out and
			// move on
			log.Printf("Expected 17 fields, got something else: %v", memline)
			continue
		}

		// Line looks good, loop over the fields
		for i := 1; i < 16; i++ {

			info := make(map[string]string)
			// Make sure we can parse the value as a float 64, else
			// we'll skip.
			floatval, fconverr := strconv.ParseFloat(parts[i], 64)
			if fconverr != nil {
				log.Printf("Cannot parse net/dev value '%v' as float64, skipping\n", parts[i])
				continue
			}

			info["device"] = iface // Device name

			// Switch on the index since each line is a different
			// XXX What are the units for the things that aren't bytes & packets?
			// Using events for now
			switch i {
			case 1:
				info["unit"] = "B"
				info["direction"] = "in"
				info["target_type"] = "count"
			case 2:
				info["unit"] = "Pckt"
				info["direction"] = "in"
				info["target_type"] = "count"
			case 3:
				info["what"] = "errors"
				info["unit"] = "Event"
				info["direction"] = "in"
				info["target_type"] = "count"
			case 4:
				info["what"] = "drop"
				info["unit"] = "Event"
				info["direction"] = "in"
				info["target_type"] = "count"
			case 5:
				info["what"] = "fifo"
				info["unit"] = "Event"
				info["direction"] = "in"
				info["target_type"] = "count"
			case 6:
				info["what"] = "frame"
				info["unit"] = "Event"
				info["direction"] = "in"
				info["target_type"] = "count"
			case 7:
				info["what"] = "compressed"
				info["unit"] = "Event"
				info["direction"] = "in"
				info["target_type"] = "count"
			case 8:
				info["what"] = "multicast"
				info["unit"] = "Event"
				info["direction"] = "in"
				info["target_type"] = "count"
			case 9:
				info["unit"] = "B"
				info["direction"] = "out"
				info["target_type"] = "count"
			case 10:
				info["unit"] = "Pckt"
				info["direction"] = "out"
				info["target_type"] = "count"
			case 11:
				info["what"] = "errs"
				info["unit"] = "Event"
				info["direction"] = "out"
				info["target_type"] = "count"
			case 12:
				info["what"] = "drop"
				info["unit"] = "Event"
				info["direction"] = "out"
				info["target_type"] = "count"
			case 13:
				info["what"] = "fifo"
				info["unit"] = "Event"
				info["direction"] = "out"
				info["target_type"] = "count"
			case 14:
				info["what"] = "colls"
				info["unit"] = "Event"
				info["direction"] = "out"
				info["target_type"] = "count"
			case 15:
				info["what"] = "carrier"
				info["unit"] = "Event"
				info["direction"] = "out"
				info["target_type"] = "count"
			case 16:
				info["what"] = "compressed"
				info["unit"] = "Event"
				info["direction"] = "out"
				info["target_type"] = "count"
			}

			metrics = append(metrics, defs.Metric{
				Timestamp: timestamp,
				Info:      info,
				Value:     floatval,
			})
		}
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
