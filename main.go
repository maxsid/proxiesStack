package main

import (
	"github.com/maxsid/proxiesStack/api"
	"github.com/maxsid/proxiesStack/cataloger"
	"github.com/maxsid/proxiesStack/config"
	"github.com/maxsid/proxiesStack/database"
	"github.com/maxsid/proxiesStack/grabber"
	"log"
	"time"
)

func main() {
	config.ParseEnv()
	database.InitDatabase()
	if !config.NoScan {
		go scan()
	}
	api.RunAPIServer()
}

func scan() {
	for {
		err := database.UnionSetsToUnionSet()
		if err != nil {
			log.Panic(err)
		}
		hostChan := make(chan string)
		doneSubChan := make(chan bool)
		go grabber.GrabProxies(hostChan, doneSubChan)
		go database.AddSetChanValue(database.UnionSetKey, hostChan, doneSubChan)

		<-doneSubChan

		go database.PopAllUnionSet(hostChan, doneSubChan)
		go cataloger.CategorizeProxies(hostChan, doneSubChan)

		<-doneSubChan

		time.Sleep(time.Duration(config.ScanInterval) * time.Second)
	}
}
