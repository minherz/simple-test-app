package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var port = flag.Int("port", 8282, "Port number to serve test page.")

const (
	responseBody = "<html><head><title>%s</title></head>" +
		"<body bgcolor=\"white\"><div style=\"max-width: 50%%; margin: auto; left: 1%%; right: 1%%; position: absolute;\"><p style=\"line-height: 2; border-radius: 25px;padding: 50px;border: 2px solid #73AD21;background-color: #DAD5D4\">" +
		"<table><tr><td><b>Application:</b></td><td style=\"width: 5px\"/><td>%v (%v)</td></tr>" +
		"<tr><td><b>Tenant:</b></td><td style=\"width: 5px\"/><td>%v</td></tr>" +
		"<tr><td><b>Server address:</b></td><td style=\"width: 5px\"/><td>%v</td></tr>" +
		"<tr><td><b>Server name:</b></td><td style=\"width: 5px\"/><td>%v</td></tr>" +
		"<tr><td><b>System time:</td><td style=\"width: 5px\"/><td>%s</td></tr>" +
		"</table></p></div</body></html>"
	unknownVersion = "&lt;unknown&gt;"
)

func main() {

	flag.Parse()

	title := os.Getenv("TITLE")
	version := os.Getenv("VERSION")

	if title == "" {
		title = "Test web app"
	}
	if version == "" {
		version = unknownVersion
	}
	ip, name := getNetValues()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var tenant = r.Header.Get("x-tenant-id")
		if tenant == "" {
			tenant = unknownVersion
		}
		w.Header().Set("x-app-title", title)
		w.Header().Set("x-app-version", version)
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, responseBody, title, title, version, tenant, ip, name, time.Now().UTC().Format("2006-01-02 15:04:05"))
	})

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "ok")
	})

	fmt.Printf("%v is started at %s on port %d\n", title, time.Now().UTC().Format("2006-01-02 15:04:05"), *port)

	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", *port),
	}

	go func() {
		err := srv.ListenAndServe()
		if err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "could not start http server: %s\n", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	err := srv.Shutdown(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not graceful shutdown http server: %s\n", err)
	}
}

func getNetValues() (serverIP string, serverName string) {

	name, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	serverName = name

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		serverIP = "0.0.0.0"
	} else {
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					serverIP = ipnet.IP.String()
				}
			}
		}
	}
	return
}
