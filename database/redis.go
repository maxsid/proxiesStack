package database

import (
	"github.com/maxsid/proxiesStack/config"
	"github.com/mediocregopher/radix"
	"log"
)

var pool *radix.Pool

const (
	WorkingSetKey    = "working"
	NotWorkingSetKey = "not_working"
	UnionSetKey      = "union"
)

// Creates connection to Redis database. Must be run before working with DB.
func InitDatabase() {
	if config.RedisHost == "" {
		log.Panic("Your environment doesn't have REDIS_HOST variable with a host of the Redis server")
	}

	newPool, err := radix.NewPool("tcp", config.RedisHost, 10)
	if err != nil {
		log.Panicf("Redis init Error: %v\n", err)
	}
	pool = newPool
}

// Fills Redis Set by values from valueChan
func AddSetChanValue(key string, valueChan <-chan string) {
	for value := range valueChan {
		err := AddSetValue(key, value)
		if err != nil {
			log.Print(err)
		}
	}
}

// Sends all records from the Union Set to hostChan
func PopAllUnionSet(hostChan chan<- string) {
	for {
		host, err := PopSetValue(UnionSetKey)
		if err != nil {
			log.Print(err)
			break
		}
		if host == "" {
			log.Println("Union set is ended")
			break
		}
		hostChan <- host
	}
	close(hostChan)
}

// Adds record to Set with 'key' name
func AddSetValue(key, value string) error {
	var count int
	err := pool.Do(radix.Cmd(&count, "SADD", key, value))
	switch {
	case err != nil:
		return err
	case count != 1:
		log.Printf("'%s' value hasn't added to '%s' set, "+
			"because this value is already exists or empty\n", value, key)
	}
	return nil
}

// Gets a number of elements of the set
func GetSetCard(key string) (int, error) {
	var count int
	err := pool.Do(radix.Cmd(&count, "SCARD", key))
	if err != nil {
		return count, err
	}
	return count, nil
}

// Pops a value from set
func PopSetValue(key string) (string, error) {
	var value string
	err := pool.Do(radix.Cmd(&value, "SPOP", key))
	if err != nil {
		return "", err
	}
	return value, nil
}

// Unites Working and NotWorking sets to Union. Working and NotWorking sets will be deleted.
func UniteSetsToUnionSet() error {
	log.Println("Uniting Sets")
	err := pool.Do(radix.Cmd(nil, "SUNIONSTORE", UnionSetKey, WorkingSetKey, NotWorkingSetKey))
	if err != nil {
		return err
	}
	log.Println("Deleting 'working' and 'not_working' sets")
	err = pool.Do(radix.Cmd(nil, "DEL", WorkingSetKey, NotWorkingSetKey))
	if err != nil {
		return err
	}
	return nil
}
