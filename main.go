package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/simulatedsimian/fffprintmon/gcode"
)

type Config struct {
	Help            bool
	PrinterHostName string
	PrinterPortNum  int
	MonitorFolder   string
}

var config Config

func init() {
	flag.BoolVar(&config.Help, "h", false, "display help")
	flag.IntVar(&config.PrinterPortNum, "p", 8899, "Printer Port number. Default: 8899")
	flag.StringVar(&config.PrinterHostName, "n", "192.168.1.101", "Printer Host Name/IP address.")
	flag.StringVar(&config.MonitorFolder, "m", ".", "Folder to monitor for uploads. Default: Current folder")

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: fffprintmon [options] ")
		flag.PrintDefaults()
	}
}

func exitOnError(err error, message string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v: %v", message, err)
		os.Exit(1)
	}

}

func main() {

	var a, b, c int
	fmt.Sscanf("Endstop: X-max: 1 Y-max: 2 Z-max: 3", "Endstop: X-max: %d Y-max: %d Z-max: %d", &a, &b, &c)
	fmt.Println(a, b, c)

	flag.Parse()

	if config.PrinterHostName == "" || config.Help {
		flag.Usage()
		os.Exit(1)
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", config.PrinterHostName, config.PrinterPortNum))
	exitOnError(err, "Failed to connect to printer")
	defer conn.Close()

	gcode := gcode.New(conn)

	run(gcode)
}

func run(g *gcode.GCode) {
	var resp gcode.CommandResponse
	var err error

	resp, err = g.CMDTakeControl()
	exitOnError(err, "")
	resp, err = g.CMDPrinterInfo()
	exitOnError(err, "")
	fmt.Println(resp.Params)

	resp, err = g.CMDHomePos()
	exitOnError(err, "")

	for {
		time.Sleep(1 *time.Second)
		resp, err = g.CMDPrinterStatus()
		exitOnError(err, "")
		fmt.Println(resp.Params)
	}
}
