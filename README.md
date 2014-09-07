# Guru

Guru is a Linux [system monitoring](http://en.wikipedia.org/wiki/System_monitoring)
daemon that collects metrics about the running system written in [Go](http://golang.org/).

# Goals

* Work on modern Linuxes (>= 2.6.11)
* Adheres to [Metrics 2.0](http://metrics20.org/) [standards](http://metrics20.org/spec/)
* No dependencies for basic operation.

# Capabilities

Guru collects the following information

* Disk statistics from `/proc/diskstats`
* Load averages from `/proc/loadavg`
* Running and total runnable threads from `/proc/loadavg`
* Memory information from `/proc/meminfo`
* Network from `/proc/net/dev`
* Virtual memory from `/proc/vminfo`
* Per-CPU and context switches from `/proc/stat`

# TODO

* df -k
* iostat?
