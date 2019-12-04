package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/maxsid/proxiesStack/database"
	"log"
	"net/http"
)

func RunAPIServer() {
	router := mux.NewRouter()
	router.HandleFunc("/working/pop", PopWorkingHandler).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func PopWorkingHandler(w http.ResponseWriter, _ *http.Request) {
	proxyHost, err := database.PopSetValue(database.WorkingSetKey)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
	} else {
		replyJSON(w, proxyHost)
	}
}

func replyJSON(w http.ResponseWriter, content interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(content)
	if err != nil {
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
	}
}
