package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

var counter = 0

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	resp, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "applications/json")
	w.WriteHeader(code)
	w.Write(resp)
}

func incrementHandler(w http.ResponseWriter, r *http.Request) {
	counter++
	fmt.Println(counter)
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func sendIncrement() {
	url := "http://localhost:8002/increment"
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Status)
}

func startServer(id int) {
	r := mux.NewRouter()
	r.HandleFunc("/increment", incrementHandler)
	port := ":" + strconv.Itoa(8000+id)
	log.Fatal(http.ListenAndServe(port, r))
}

func startClient() {
	scanner := bufio.NewScanner(os.Stdin)
	text := ""
	for text != "exit" {
		scanner.Scan()
		text := scanner.Text()
		switch text {
		case "exit":
			os.Exit(0)
		case "increment":
			sendIncrement()
		}
	}
}

func main() {
	clusterSize := flag.Int("cluster-size", 3, "number of nodes in cluster")
	id := flag.Int("id", 1, "node id")
	flag.Parse()
	fmt.Println("Cluster size:", *clusterSize)
	fmt.Println("Node id:", *id)

	go startServer(*id)
	startClient()
}
