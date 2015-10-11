package main

import (
	"errors"

	"github.com/boltdb/bolt"
)

var brain *Brain

type Brain struct {
	db *bolt.DB
}

func initBrain(c Config) {
	brain = new(Brain)
	var err error
	brain.db, err = bolt.Open(c.DatabaseFile, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func (brain *Brain) Close() {
	brain.db.Close()
}

func (brain *Brain) Put(namespace, key, value string) (err error) {
	err = brain.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(namespace))
		if err == nil {
			err = b.Put([]byte(key), []byte(value))
		}
		return err
	})
	return
}

func (brain *Brain) Delete(namespace, key string) (err error) {
	err = brain.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(namespace))
		if err == nil {
			err = b.Delete([]byte(key))
		}
		return err
	})
	return
}

func (brain *Brain) Get(namespace, key string) (value string, err error) {
	err = brain.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(namespace))
		if b != nil {
			value = string(b.Get([]byte(key)))
		} else {
			return errors.New("Bucket not found")
		}
		return nil
	})
	return
}

func (brain *Brain) GetAll(namespace string) (kv map[string]string, err error) {
	kv = make(map[string]string)

	err = brain.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(namespace))
		if b == nil {
			return errors.New("Bucket not found")
		}

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			kv[string(k)] = string(v)
		}
		return nil
	})
	return
}
