package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"github.com/wadahana/memu"
	"github.com/wadahana/memu/log"
	"github.com/wadahana/ga/core"
	"github.com/wadahana/ga/display"
	"github.com/wadahana/ga/encoders"
	"github.com/wadahana/ga/rtc"
	"github.com/wadahana/ga/api"
)

const (
	defaultHttpPort   = "9500"
	defaultStunServer = "stun:stun.ideasip.com"
	defaultTurnServer = "turn:turn.ideasip.com"
	defaultConfigFile = "./ga.yaml"
	version = "v0.1.0"
)

var (
	dumpHelp   bool
	configFile string
)


func usage() {

	fmt.Printf("ga version: ga/%s\r\n", version)
	fmt.Printf("Usage: ga [-ch]\r\n")
	fmt.Printf("           -h print this message\r\n")
	fmt.Printf("           -c config file path\r\n")
}

func init() {
	flag.BoolVar(&dumpHelp,         "h", false,     "show usage")
	flag.StringVar(&configFile, "c", "./ga.yaml",   "config file")
	flag.Usage = usage
}

func main() {
	flag.Parse()
	if dumpHelp {
		usage()
		os.Exit(0)
	}
	err := core.ParseYamlFile(configFile) 
	if err != nil {
		panic(err);
	}
	config := &memu.MEmuConfig{
		MEmuPath:   core.Setting.MEmu.Path,
		LoggerFile: "console",
	}
	memu.Init(config)

	log.Infof("ga-server start...")

	var disp display.Service
	disp, err = display.NewProvider()
	if err != nil {
		log.Fatalf("Can't init display: %v", err)
	}

	var enc encoders.Service = &encoders.EncoderService{}
	if err != nil {
		log.Fatalf("Can't create encoder service: %v", err)
	}

	var webrtc rtc.Service
	webrtc = rtc.NewRemoteScreenService(disp, enc)

	for i, s := range core.Setting.Rtc.IceServers {
		log.Debugf("ice[%d]: %s, %s, %s\n", i, s.Url, s.Username, s.Password);
	}
	
	mux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", api.MakeApiHandler(webrtc, disp)))
	mux.Handle("/vm/", http.StripPrefix("/vm", api.MakeVMHandler()))
	mux.Handle("/rdp/", http.StripPrefix("/rdp", api.MakeRdpHandler()))
	mux.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./html"))))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Debugf("url: %s", r.URL);
		if r.URL.Path != "/" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		http.ServeFile(w, r, "./html/index.html")
	})

	errors := make(chan error, 2)
	go func() {
		log.Infof("Starting signaling server on port %d", core.Setting.Web.Port)
		errors <- http.ListenAndServe(fmt.Sprintf(":%d", core.Setting.Web.Port), mux)
	}()

	go func() {
		interrupt := make(chan os.Signal)
		signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
		errors <- fmt.Errorf("Received %v signal", <-interrupt)
	}()

	err = <-errors
	log.Infof("%s, exiting.", err)
}
