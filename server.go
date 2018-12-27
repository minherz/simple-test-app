package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const (
	version         = "2.0"
	refreshInterval = 30000
	pageTemplate    = `<html>
<head><title>%v</title></head>
<body bgcolor="white">
 <div style="max-width: 50%%; margin: auto; left: 1%%; right: 1%%; position: absolute;">
  <div id="info">
  </div>
  <div id="customize">
 	  <table>
 	   <tr><td>Custom Field Name:</td><td style="width: 5px"/><td><input type="text" id="fname"/><br></td></tr>
 	   <tr><td>Custom Field Value:</td><td style="width: 5px"/><td><input type="text" id="fvalue"/><br></td></tr>
 	   <tr><td><button type="button" onclick="setHeader()">Setup</button></td></tr>
 	  </table>
  </div>
 </div>
</body>
<script type="text/javascript">
//<![CDATA[
function refreshInfo() {
	var xhttp = new XMLHttpRequest();
	xhttp.onload = function() {
		if (this.readyState == 4 && this.status == 200) {
			document.getElementById("info").innerHTML = xhttp.responseText;
		}
	}
	xhttp.open("GET", location.href + "info", true)
	if (customField !== "") {
		xhttp.setRequestHeader(customField, customFieldValue)
	}
	xhttp.send()
	window.setTimeout(refreshInfo,%v);
}
var customField = "", customFieldValue = ""
function setHeader() {
	customField = document.getElementById("fname").value
	customFieldValue = document.getElementById("fvalue").value
}
refreshInfo();
//]]>
</script>
</html>`
	infoTable = `	<p style="line-height: 2; border-radius: 25px;padding: 50px;border: 2px solid #73AD21;background-color: #DAD5D4">
<table>
 <tr><td><b>Application:</b></td><td style="width: 5px"/><td>%v</td></tr>
 <tr><td><b>Server address:</b></td><td style="width: 5px"/><td>%v</td></tr>
 <tr><td><b>Server name:</b></td><td style="width: 5px"/><td>%v</td></tr>
 <tr><td><b>System time:</td><td style="width: 5px"/><td>%s</td></tr>
</table>
</p>
`
)

var title, ip, hostname string

func main() {
	port := 8282

	// parse args (TBD: use package)
	args := os.Args[1:]
	if len(args) > 0 {
		if args[0] == "--version" {
			fmt.Println(version)
			os.Exit(0)
		}
		n, err := strconv.ParseInt(args[0], 10, 16)
		if err != nil {
			fmt.Println("Invalid port value: ", args[0])
			os.Exit(1)
		}
		port = int(n)
	}

	ip, hostname = getNetValues()
	title = getTitle()

	// setup handlers
	http.HandleFunc("/", getIndex)
	http.HandleFunc("/info", getInfo)
	http.HandleFunc("/healthz", checkHealth)

	fmt.Printf("%v is started at %s on port %d\n", title, time.Now().UTC().Format("2006-01-02 15:04:05"), port)

	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
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

func getIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("x-app-version", version)
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, pageTemplate, fmt.Sprintf("Simple Test Application (%v)", version), refreshInterval)
}

func getInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("x-app-version", version)
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, infoTable, title, ip, hostname, time.Now().UTC().Format("2006-01-02 15:04:05"))
}

func checkHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
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

func getTitle() string {
	s := os.Getenv("TITLE")
	if len(s) == 0 {
		s = fmt.Sprintf("Simple Test Application (%v)", version)
	}
	return s
}
