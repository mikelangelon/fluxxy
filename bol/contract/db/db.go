package db

import (
	"github.com/inconshreveable/log15"
	"gopkg.in/mgo.v2"
	"sync"
)

var counters map[string]int64
var lock sync.Mutex

const dbKey = "__db"

var LeakCount int64 = 2

func init() {
	counters = make(map[string]int64)
}

type Collection struct {
	*mgo.Collection
}

type Db struct {
	*mgo.Database
}

func inc(name string) {
	lock.Lock()
	counters[name] += 1
	if counters[name] > LeakCount {
		log15.Error("Leaking database sessions for key: ", "key", name, "count", counters[name])
	}
	lock.Unlock()
}

func decr(name string) {
	lock.Lock()
	counters[name] -= 1
	lock.Unlock()
}

func (session *Db) Close() {
	decr(dbKey)
	session.Session.Close()
}

func (session *Collection) Close() {
	decr(session.Name)
	session.Database.Session.Close()
}

func CollectionSupplier(session *mgo.Session, collection string) func() *Collection {

	return func() *Collection {
		inc(collection)
		return &Collection{session.Copy().DB("contract_management").C(collection)}
	}
}

func DbSupplier(session *mgo.Session) func() *Db {

	return func() *Db {
		inc(dbKey)
		return &Db{session.Copy().DB("")}
	}
}
