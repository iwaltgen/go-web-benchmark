package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/magefile/mage/sh"
)

var (
	http2    = true
	port     = 8080
	tps      = 10000
	sleep    = 100 * time.Millisecond
	duration = 10 * time.Second
	message  = []byte("hello world")
	tls      = struct {
		key, cert string
	}{
		key:  "localhost+2-key.pem",
		cert: "localhost+2.pem",
	}
)

func init() {
	flag.BoolVar(&http2, "http2", http2, "HTTP/2")
	flag.IntVar(&port, "port", port, "HTTP port")
	flag.DurationVar(&sleep, "sleep", sleep, "HTTP handler mock processing time")
	flag.DurationVar(&duration, "duration", duration, "stress test duration time")
	flag.IntVar(&tps, "tps", tps, "stress test request per second")
	flag.StringVar(&tls.cert, "tls-cert", tls.cert, "TLS cert file path")
	flag.StringVar(&tls.key, "tls-key", tls.key, "TLS key file path")

	flag.Parse()
}

func isHTTPS() bool {
	return tls.key != "" && tls.cert != ""
}

func mockProcessing() {
	if 0 < sleep {
		time.Sleep(sleep)
		return
	}

	runtime.Gosched()
}

func bench(name string) error {
	color.Green("run http engine...")
	defer color.Magenta("shutdown http engine...")

	switch name {
	case "echo":
		s := serveEcho()
		defer shutdownEcho(s)

	case "gin":
		s := serveGin()
		defer shutdownGin(s)

	case "fiber":
		s := serveFiber()
		defer shutdownFiber(s)

	default:
		s := serveDefault()
		defer shutdownDefault(s)

	}

	time.Sleep(time.Second)

	var hists []string
	for i := 0; i < 5; i++ {
		v := time.Duration(i * 10)
		hists = append(hists, fmt.Sprintf("%v", sleep+(v*time.Millisecond)))
	}

	color.Cyan("run warm-up...")
	cmd := fmt.Sprintf("echo 'GET https://:8080/hello' "+
		"| vegeta attack -name=%s -http2=%t -insecure -duration=%v "+
		"| vegeta report -type='hist[%s]'", name, http2, duration/10, strings.Join(hists, ","))
	if err := sh.RunV("sh", "-c", cmd); err != nil {
		log.Fatalf("warm-up error: %v", err)
	}

	color.Cyan("run stress tool...")
	id := fmt.Sprintf("%s-%v-%dtps", name, sleep, tps)
	cmd = fmt.Sprintf("echo 'GET https://:8080/hello' "+
		"| vegeta attack -name=%s -http2=%t -insecure -rate=%d -duration=%v "+
		"| tee report/%s.bin "+
		"| vegeta report", id, http2, tps, duration, id)
	return sh.RunV("sh", "-c", cmd)
}

func main() {
	if err := os.MkdirAll("report", os.ModePerm); err != nil {
		log.Fatalf("mkdir error: %v", err)
	}

	engines := []string{"echo", "gin", "fiber", "default"}
	var plotArgs []string
	for _, name := range engines {
		if err := bench(name); err != nil {
			log.Fatalf("benchmark error: %v", err)
		}

		id := fmt.Sprintf("%s-%v-%dtps", name, sleep, tps)
		plotArgs = append(plotArgs, fmt.Sprintf("report/%s.bin", id))
	}

	cmd := fmt.Sprintf("vegeta plot %s", strings.Join(plotArgs, " "))
	if err := sh.RunV("sh", "-c", fmt.Sprintf(cmd+" > report/plot-%v-%dtps.html", sleep, tps)); err != nil {
		log.Fatalf("plot error: %v", err)
	}
}
