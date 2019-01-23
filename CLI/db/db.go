package db

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/boltdb/bolt"
)

// Task structure
type Task struct {
	Key   int
	Value string
}

var taskBucket = []byte("tasks")
var db *bolt.DB

// Init Initialize & open the database connection
func Init(dbPath string) error {
	var err error
	db, err = bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		fmt.Printf("could not open db, %v\n", err)
		return err

	}

	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(taskBucket)
		fmt.Printf("could not create bucket: %v\n", err)
		return err
	})
}

// CreateTask This function is used to add the task in database and it will return the taskno with error
func CreateTask(task string) (int, error) {
	var taskno int

	err := db.Update(func(tx *bolt.Tx) error {
		taskbuck := tx.Bucket(taskBucket)
		taskno64, _ := taskbuck.NextSequence()
		taskno = int(taskno64)
		key := itob(int(taskno))
		return taskbuck.Put(key, []byte(task))

	})

	fmt.Println("Added task")
	return taskno, err
}

// GetAllLists This function is used to list all the  tasks present in database and it will return the tasks with error
func GetAllLists() ([]Task, error) {

	var tasks []Task
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {

			tasks = append(tasks, Task{
				Key:   btoi(k),
				Value: string(v),
			})
		}
		return nil

	})
	return tasks, err
}

// DeleteTasks This function is used to delete the added task in database and it will return the error if task is not deleted from the db
func DeleteTasks(key int) error {

	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		return b.Delete(itob(key))
	})

}

// itob This function is used to convert integer value in byte

func itob(v int) []byte {

	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

// btoi This function is used to convert byte to integer

func btoi(b []byte) int {

	return int(binary.BigEndian.Uint64(b))
}
