package main

import (
	"flag"
	"fmt"
	"encoding/json"
	"log"
	"net/http"

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
	
func main() {
	clusterSize := flag.Int("cluster-size", 3, "number of nodes in cluster")
	id := flag.Int("id", 1, "node id")
	flag.Parse()
	fmt.Println("Cluster size:", *clusterSize)
	fmt.Println("Node id:", *id)
	
	r := mux.NewRouter()
	r.HandleFunc("/increment", incrementHandler)
	log.Fatal(http.ListenAndServe(":8000", r))
}

