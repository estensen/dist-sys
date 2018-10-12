package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var counter = 0

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "applications/json")
	w.WriteHeader(code)
	w.Write(response)
}

func incrementHandler(w http.ResponseWriter, r *http.Request) {
	counter++
	fmt.Println(counter)
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func startServer() {
	r := mux.NewRouter()
	r.HandleFunc("/increment", incrementHandler)
	log.Fatal(http.ListenAndServe(":8000", r))
}

func startClient() {
	scanner := bufio.NewScanner(os.Stdin)
	text := ""
	for text != "exit" {
		scanner.Scan()
		text := scanner.Text()
		fmt.Println(text)
	}
}

func main() {
	clusterSize := flag.Int("cluster-size", 3, "number of nodes in cluster")
	id := flag.Int("id", 1, "node id")
	flag.Parse()
	fmt.Println("Cluster size:", *clusterSize)
	fmt.Println("Node id:", *id)

	go startServer()
	startClient()
}
