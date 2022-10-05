package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/buaazp/fasthttprouter"
	"github.com/hashicorp/logutils"
	flags "github.com/jessevdk/go-flags"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
)

var opts struct {
	ListenNetwork   string `long:"listen-network" env:"LISTEN_NETWORK" default:"tcp4" description:"listen network"`
	ListenAddr      string `long:"listen-addr" env:"LISTEN_ADDR" default:":4400" description:"listen address"`
	Dbg             bool   `long:"debug" env:"DEBUG" description:"debug mode"`
	FilePath        string `long:"file-path" env:"FILE_PATH" description:"file path"`
	FileScanTimeout int    `long:"file-scan-timeout" env:"FILE_SCAN_TIMEOUT" default:"10" description:"scan file timeout in seconds"`
}

type application struct {
	HTTP   *fasthttp.Server
	IPList map[string][]IPItem
}

//IPItem struct
type IPItem struct {
	IP       net.IP
	IPNet    *net.IPNet
	isSubnet bool
}

func main() {

	p := flags.NewParser(&opts, flags.Default)
	if _, e := p.ParseArgs(os.Args[1:]); e != nil {
		os.Exit(1)
	}

	setupLog(opts.Dbg)

	log.Printf("[INFO] ip checker started ... ")

	list := make(map[string][]IPItem)
	var err error

	if list, err = readLines(opts.FilePath); err != nil {
		log.Printf("[ERROR] file with ips error: %+v", err)
		os.Exit(1)
	}

	app := &application{
		HTTP: &fasthttp.Server{
			Name: "IP-CHECKER",
		},
		IPList: list,
	}

	listener, err := reuseport.Listen(opts.ListenNetwork, opts.ListenAddr)
	if err != nil {
		log.Printf("[ERROR] error create reuseport listener, %+v", err)
		os.Exit(1)
	}

	mainTicker := time.NewTicker(time.Second * time.Duration(opts.FileScanTimeout))

	go func() {
		for range mainTicker.C {
			if list, err = readLines(opts.FilePath); err != nil {
				log.Printf("[ERROR] file with ips error: %+v", err)
			}
			app.IPList = list
		}
	}()

	//graceful stop
	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	go func() {
		log.Printf("[INFO] started graceful stop goroutine")
		sig := <-gracefulStop
		log.Printf("[INFO] sig graceful stop IP-CHECKER: %+v", sig)
		err := listener.Close()
		if err != nil {
			log.Printf("[ERROR] listener close error IP-CHECKER: %+v", err)
			return
		}
		return
	}()

	router := fasthttprouter.New()

	router.GET("/validate/:ip", app.checkIP)

	app.HTTP.Handler = router.Handler

	if err := app.HTTP.Serve(listener); err != nil {
		log.Printf("[ERROR] error serve")
		os.Exit(1)
	}
}

func setupLog(dbg bool) {
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel("INFO"),
		Writer:   os.Stdout,
	}

	log.SetFlags(log.Ldate | log.Ltime)

	if dbg {
		log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
		filter.MinLevel = logutils.LogLevel("DEBUG")
	}

	log.SetOutput(filter)
}
