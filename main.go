package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type PiholeClientRequest struct {
	Password string `json:"password"`
	Groups   []int  `json:"groups"`
}

var (
	piholeURL      string
	piholePassword string
)

func main() {
	var ok bool

	piholeURL, ok = os.LookupEnv("PIHOLE_URL")
	if !ok {
		panic("PIHOLE_URL not set")
	}

	piholePassword, ok = os.LookupEnv("PIHOLE_PASSWORD")
	if !ok {
		panic("PIHOLE_PASSWORD not set")
	}

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8001"
	}

	http.HandleFunc("/on", setBlocking(true))
	http.HandleFunc("/off", setBlocking(false))

	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func setBlocking(on bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientName := r.Header.Get("X-Forwarded-For")
		endpointURL := fmt.Sprintf("%s/api/clients/%s", piholeURL, clientName)

		clientsRequest := PiholeClientRequest{
			Password: piholePassword,
			Groups:   []int{},
		}

		if on {
			clientsRequest.Groups = []int{0}
		}

		log.Printf("Setting blocking=%v for %s", on, clientName)

		w.Header().Add("Content-Type", "application/json")

		json, err := json.Marshal(clientsRequest)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "Failed to build API payload"}`))
			return
		}

		req, err := http.NewRequest("PUT", endpointURL, bytes.NewBuffer(json))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "Failed to build API request"}`))
			return
		}

		req.Header.Set("Content-Type", "application/json")
		httpClient := &http.Client{}
		res, err := httpClient.Do(req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"message": "Failed to send API request"}`))
			return
		}
		defer res.Body.Close()

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Success"}`))
	}
}
