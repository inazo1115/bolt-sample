package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/boltdb/bolt"
)

const (
	DBName     = "my.db"
	BucketName = "MyBucket"
)

var db *bolt.DB

func pingHandleFunc(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong")
}

func putHandleFunc(w http.ResponseWriter, r *http.Request) {

	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	if len(key) == 0 || len(value) == 0 {
		fmt.Fprintf(w, "required parameters are: key, value\n")
		return
	}

	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BucketName))
		err := b.Put([]byte(key), []byte(value))
		return err
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Fprintf(w, "ok\n")
}

func getHandleFunc(w http.ResponseWriter, r *http.Request) {

	key := r.URL.Query().Get("key")
	if len(key) == 0 {
		fmt.Fprintf(w, "required parameter is: key\n")
		return
	}

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BucketName))
		value := b.Get([]byte(key))
		if len(value) == 0 {
			fmt.Fprintf(w, "not found\n")
		} else {
			fmt.Fprintf(w, "key=%s, value=%s\n", key, value)
		}
		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func listHandleFunc(w http.ResponseWriter, r *http.Request) {

	prefix := []byte(r.URL.Query().Get("prefix"))

	err := db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(BucketName)).Cursor()
		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
			fmt.Fprintf(w, "key=%s, value=%s\n", k, v)
		}
		return nil
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func backupHandleFunc(w http.ResponseWriter, r *http.Request) {
	err := db.View(func(tx *bolt.Tx) error {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", `attachment; filename="my.db"`)
		w.Header().Set("Content-Length", strconv.Itoa(int(tx.Size())))
		_, err := tx.WriteTo(w)
		return err
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {

	var err error
	db, err = bolt.Open(DBName, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(BucketName))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/ping", pingHandleFunc)
	http.HandleFunc("/put", putHandleFunc)
	http.HandleFunc("/get", getHandleFunc)
	http.HandleFunc("/list", listHandleFunc)
	http.HandleFunc("/backup", backupHandleFunc)
	http.ListenAndServe(":8080", nil)
}
