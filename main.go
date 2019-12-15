package main

import (
	"github.com/maxsid/proxiesStack/api"
	"github.com/maxsid/proxiesStack/classifier"
	"github.com/maxsid/proxiesStack/config"
	db "github.com/maxsid/proxiesStack/database"
	"github.com/maxsid/proxiesStack/grabber"
	"log"
	"time"
)

func main() {
	log.Println("Parse OS Environment")
	config.ParseEnv()
	log.Println("Init database")
	db.InitDatabase()
	if !config.NoScan {
		go scan()
	}
	api.RunAPIServer()
}

// Main scan function. Unites sets, finds hosts on GrabPage and checks them available
func scan() {
	for {
		log.Println("Scan is running")
		config.ScanStatus = "running"
		err := db.UniteSetsToUnionSet()
		if err != nil {
			log.Panic(err)
		}

		log.Println("Grab is starting")
		hostChan := make(chan string)
		go grabber.GrabProxies(hostChan)
		db.AddSetChanValue(db.UnionSetKey, hostChan)

		log.Println("Classification is starting")
		hostChan = make(chan string)
		workingChan, notWorkingChan := make(chan string), make(chan string)
		go db.AddSetChanValue(db.WorkingSetKey, workingChan)
		go db.AddSetChanValue(db.NotWorkingSetKey, notWorkingChan)
		go classifier.ClassifyProxies(hostChan, workingChan, notWorkingChan)
		db.PopAllUnionSet(hostChan)

		log.Println("Scan is sleeping")
		config.ScanStatus = "sleeping"
		time.Sleep(time.Duration(config.ScanInterval) * time.Second)
	}
}
