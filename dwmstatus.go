package main

// #cgo LDFLAGS: -lX11
// #include <X11/Xlib.h>
import "C"

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

var dpy = C.XOpenDisplay(nil)

func getBatteryPercentage(path string) (perc int, err error) {
	energy_now, err := ioutil.ReadFile(fmt.Sprintf("%s/energy_now", path))
	if err != nil {
		perc = -1
		return
	}
	energy_full, err := ioutil.ReadFile(fmt.Sprintf("%s/energy_full", path))
	if err != nil {
		perc = -1
		return
	}
	var enow, efull int
	fmt.Sscanf(string(energy_now), "%d", &enow)
	fmt.Sscanf(string(energy_full), "%d", &efull)
	perc = enow * 100 / efull
	return
}

func getWatts(path string) (watts float64, err error) {
	power_now, err := ioutil.ReadFile(fmt.Sprintf("%s/power_now", path))
	if err != nil {
		watts = -1.0
		return
	}
	var pnow int
	fmt.Sscanf(string(power_now), "%d", &pnow)
	watts = float64(pnow) / 1e6
	return
}


func getLoadAverage(file string) (lavg string, err error) {
	loadavg, err := ioutil.ReadFile(file)
	if err != nil {
		return "Couldn't read loadavg", err
	}
	lavg = strings.Join(strings.Fields(string(loadavg))[:3], " ")
	return
}

func setStatus(s *C.char) {
	C.XStoreName(dpy, C.XDefaultRootWindow(dpy), s)
	C.XSync(dpy, 1)
}

func formatStatus(format string, args ...interface{}) *C.char {
	status := fmt.Sprintf(format, args...)
	return C.CString(status)
}

func main() {
	if dpy == nil {
		log.Fatal("Can't open display")
	}
	for {
		t := time.Now().Format("Mon, 2 Jan • 15:04")
		b, err := getBatteryPercentage("/sys/class/power_supply/BAT0")
		if err != nil {
			log.Println(err)
		}
		l, err := getLoadAverage("/proc/loadavg")
		if err != nil {
			log.Println(err)
		}
		w, err := getWatts("/sys/class/power_supply/BAT0")
		if err != nil {
			log.Println(err)
		}
		s := formatStatus(" %s • %.3gW • %d%% • %s", l, w, b, t)
		setStatus(s)
		time.Sleep(time.Second * 5)
	}
}
