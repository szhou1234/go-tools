// type in http://localhost:1718 in browser to see what will happen.
package main

import (
	"flag"
	"html/template"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
)

const templateStr = `
<html>
<head>
<title>QR Link Generator</title>
</head>
<body>
{{if .}}
<img src="http://chart.apis.google.com/chart?chs=300x300&cht=qr&choe=UTF-8&chl={{.}}" />
<br>
{{.}}
<br>
<br>
{{end}}
<form action="/qr" name=f method="GET"><input maxLength=1024 size=70
name=s value="" title="Text to QR Encode"><input type=submit
value="Show QR" name=qr>
</form>
</body>
</html>
`

var addr = flag.String("addr", ":1718", "http service address") // Q=17, R=18

var templ = template.Must(template.New("qr").Parse(templateStr))

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stderr instead of stdout, could also be a file.
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

func main() {
	flag.Parse()

	// This means that access to http://<IP>:<PORT> will trigger execution of
	// function QR(). To avoid that, use "/qr" here and also in the html template.
	http.Handle("/qr", http.HandlerFunc(qrHandler))
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

// QR is the handler for equests.
func qrHandler(w http.ResponseWriter, req *http.Request) {
	log.Info("Serving request from ", req.RemoteAddr)
	templ.Execute(w, req.FormValue("s"))
}
