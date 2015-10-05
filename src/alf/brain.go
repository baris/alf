package main

import (
	"errors"

	"github.com/boltdb/bolt"
)

type Brain struct {
	db *bolt.DB
}

func NewBrain(dbfile string) (brain *Brain) {
	brain = new(Brain)
	var err error
	brain.db, err = bolt.Open(dbfile, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	return brain
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