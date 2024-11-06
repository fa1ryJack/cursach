package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func getData(w http.ResponseWriter, r *http.Request) {
    data, err := FetchMyLikes(r.URL.Query().Get("profile"))
	if err != nil{
    	w.Header().Set("Content-Type", "application/json")
		resp := make(map[string]string)
		resp["error"] = err.Error()
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonResp)

	}else{
    	w.Header().Set("Content-Type", "application/json")
		jsonResp, err := json.Marshal(data)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
    	w.WriteHeader(http.StatusOK)
		w.Write(jsonResp)
	}
}

func main() {
    http.HandleFunc("/data", getData)
	fmt.Println("localhost:8080/")
    log.Fatal(http.ListenAndServe(":8080", nil))
}