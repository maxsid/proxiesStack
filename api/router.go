package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/maxsid/proxiesStack/config"
	db "github.com/maxsid/proxiesStack/database"
	"log"
	"net/http"
)

func RunAPIServer() {
	router := mux.NewRouter()
	router.HandleFunc("/working/pop", popWorkingHandler).Methods(http.MethodGet)
	router.HandleFunc("/info", getInfo).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func popWorkingHandler(w http.ResponseWriter, _ *http.Request) {
	proxyHost, err := db.PopSetValue(db.WorkingSetKey)
	switch {
	case err != nil:
		log.Println(err)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
	case proxyHost == "":
		w.WriteHeader(http.StatusNoContent)
	default:
		err = db.AddSetValue(db.NotWorkingSetKey, proxyHost)
		if err != nil {
			log.Println(err)
		}
		replyJSON(w, proxyHost)
	}
}

func getInfo(w http.ResponseWriter, _ *http.Request) {
	type status struct {
		Working          int    `json:"working"`
		NotWorking       int    `json:"not_working"`
		Union            int    `json:"union"`
		ScanStatus       string `json:"scan_status"`
		ScanInterval     int    `json:"scan_interval"`
		TimeoutHTTP      int    `json:"timeout_http"`
		GrabPageAddress  string `json:"grab_page_address"`
		GrabPagePattern  string `json:"grab_page_pattern"`
		CheckPageAddress string `json:"check_page_address"`
		CheckPagePattern string `json:"check_page_pattern"`
	}
	var (
		setsCards [3]int
		err       error
	)
	for i, set := range []string{db.WorkingSetKey, db.NotWorkingSetKey, db.UnionSetKey} {
		setsCards[i], err = db.GetSetCard(set)
		if err != nil {
			log.Println(err)
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
	statusData := status{
		Working:          setsCards[0],
		NotWorking:       setsCards[1],
		Union:            setsCards[2],
		ScanStatus:       config.ScanStatus,
		ScanInterval:     config.ScanInterval,
		TimeoutHTTP:      config.TimeoutHTTP,
		GrabPageAddress:  config.GrabPageAddress,
		GrabPagePattern:  config.GrabPagePattern,
		CheckPageAddress: config.CheckPageAddress,
		CheckPagePattern: config.CheckPagePattern,
	}
	replyJSON(w, statusData)
}

func replyJSON(w http.ResponseWriter, content interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(content)
	if err != nil {
		log.Println(err)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
	}
}
