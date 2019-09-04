package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
)

func main() {

	var (
		addr = "127.0.0.1:8000"
		dir  = "/var/pickup"
	)

	flag.StringVar(&addr, "listen", addr, "host:port to listen for http requests")
	flag.StringVar(&dir, "pickup", dir, "directory used for pickup by powermta")
	flag.Parse()

	_, err := os.Stat(dir)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Emails will be placed in %s for pickup", dir)

	h := &handler{
		pickup: pickup(dir),
	}
	http.Handle("/messages", h)

	log.Printf("Post messages to http://%s/messages", addr)
	err = http.ListenAndServe(addr, nil)
	// One can use generate_cert.go in crypto/tls to generate cert.pem and key.pem.
	//err := http.ListenAndServeTLS(":8443", "cert.pem", "key.pem", nil)
	log.Fatal(err)

}

type handler struct {
	pickup pickup
}

func (h *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	msg := new(message)
	dec := json.NewDecoder(req.Body)
	err := dec.Decode(msg)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if msg.From == nil {
		log.Printf("\"from\" missing")
		http.Error(w, "\"from\" missing", http.StatusBadRequest)
		return
	}

	err = h.pickup.submit(msg)
	if err != nil {
		log.Printf("%v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
