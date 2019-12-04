package database

import (
	"fmt"
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

func InitDatabase() {
	if config.RedisHost == "" {
		log.Panic("Your environment doesn't have REDIS variable with a host of the Redis server")
	}

	newPool, err := radix.NewPool("tcp", config.RedisHost, 10)
	if err != nil {
		log.Panicf("Redis init Error: %v\n", err)
	}
	pool = newPool
}

func PopAllUnionSet(hostChan chan string, done chan bool) {
	for {
		host, err := PopSetValue(UnionSetKey)
		if err != nil {
			log.Print(err)
			break
		}
		hostChan <- host
	}
	done <- true
}

func AddSetValue(key, value string) error {
	var count int
	err := pool.Do(radix.Cmd(&count, "SADD", key, value))
	switch {
	case err != nil:
		log.Printf("Redis RPUSH Error: %v\n", err)
		return err
	case count != 1:
		err := fmt.Errorf("Redis RPUSH Error: Count isn't 1, it's %d\n", count)
		log.Print(err)
		return err
	}
	return nil
}

func AddSetChanValue(key string, valueChan chan string, done chan bool) {
	for {
		select {
		case value := <-valueChan:
			err := AddSetValue(key, value)
			if err != nil {
				log.Print(err)
			}
		case <-done:
			return
		}
	}
}

func PopSetValue(key string) (string, error) {
	var value string
	err := pool.Do(radix.Cmd(&value, "SPOP", key))
	if err != nil {
		log.Printf("Redis LPOP Error: %v\n", err)
		return "", err
	}
	return value, nil
}

func UnionSets(toKey string, fromKeys ...string) error {
	cmdArgs := append([]string{toKey}, fromKeys...)
	err := pool.Do(radix.Cmd(nil, "SUNIONSTORE", cmdArgs...))
	if err != nil {
		log.Printf("Union Lists Error: %v", err)
		return err
	}
	return nil
}

// Future Feature
//func IsSetMember(key, value string) bool {
//	var answer int
//	err := pool.Do(radix.Cmd(&answer, "SISMEMBER", key, value))
//	if err != nil {
//		log.Print(err)
//		return false
//	}
//	return answer == 1
//}

func UnionSetsToUnionSet() error {
	return UnionSets(UnionSetKey, WorkingSetKey, NotWorkingSetKey)
}
