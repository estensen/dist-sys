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
var clusterSize int
var id int
var nodeURLs = make([]string, clusterSize)

func generateNodeURLs(clusterSize int) []string {
	nodeURLs := make([]string, clusterSize)
	for i := 1; i <= clusterSize; i++ {
		url := "http://localhost:" + strconv.Itoa(8000+i)
		nodeURLs[i-1] = url
	}
	fmt.Println("Available nodes:", len(nodeURLs))
	for i := range nodeURLs {
		fmt.Println(nodeURLs[i])
	}
	return nodeURLs
}

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
	fmt.Println("Local value:", counter)
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func sendIncrement() {
	for i := range nodeURLs {
		url := nodeURLs[i] + "/increment"
		_, err := http.Post(url, "application/json", nil)
		if err != nil {
			fmt.Printf("Node %d not reachable\n", i)
		} else {
			fmt.Printf("Node %d incremented!\n", i)
		}
	}
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

	nodeURLs = generateNodeURLs(*clusterSize)
	go startServer(*id)
	startClient()
}
