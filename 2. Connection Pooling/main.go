package main

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

type DBPool struct {
	// mu          sync.Mutex
	// channel     chan interface{}
	// connections []*sql.DB
	pool chan *sql.DB
}

func createNewDBConnection() (*sql.DB, error) {
	db_url := "postgresql://tapesh:tapesh@localhost:5432/hoteldb?sslmode=disable"
	db, err := sql.Open("postgres", db_url)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("Cannot reach the database: ", err)
	}
	return db, nil
}

func dummyDBOperation(db *sql.DB) error {
	query := `SELECT pg_sleep(0.01);` // Simulating a small delay
	_, err := db.Exec(query)
	return err
}

func runWithoutPooling(queriesCount int) (time.Duration, error) {
	start := time.Now()

	// start queriesCount threads to simulate multiple DB connections
	var wg sync.WaitGroup
	var once sync.Once
	var operErr error

	for i := 0; i < queriesCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			connection, err := createNewDBConnection()
			if err != nil {
				once.Do(func() {
					operErr = err
				})
				return
			}

			defer connection.Close()

			if err := dummyDBOperation(connection); err != nil {
				once.Do(func() {
					operErr = err
				})
			}
		}()
	}

	wg.Wait()
	return time.Since(start), operErr
}

func newPool(size int) *DBPool {
	p := &DBPool{
		pool: make(chan *sql.DB, size),
	}

	for i := 0; i < size; i++ {
		db, _ := createNewDBConnection()
		p.pool <- db
	}
	return p
}

func (p *DBPool) Get() *sql.DB {
	return <-p.pool
}

func (p *DBPool) Put(db *sql.DB) {
	p.pool <- db
}

func (p *DBPool) Close() {
	close(p.pool)
}

func runWithPooling(queriesCount int) (time.Duration, error) {
	pool := newPool(100)

	start := time.Now()
	var wg sync.WaitGroup
	var once sync.Once
	var operErr error

	for i := 0; i < queriesCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			connection := pool.Get()

			defer pool.Put(connection)
			if err := dummyDBOperation(connection); err != nil {
				once.Do(func() {
					operErr = err
				})
			}
		}()
	}

	wg.Wait()
	pool.Close()
	return time.Since(start), operErr
}

/* Time DIFF with or without Connection Pooling */
func main() {

	queriesCounts := [4]int{10, 50, 100, 300}

	// 1. Run without Connection Pooling
	fmt.Println("Running without Connection Pooling...")

	for _, count := range queriesCounts {
		fmt.Printf("Executing %d queries without connection pooling...\n", count)
		time_taken, err := runWithoutPooling(count)

		if err != nil {
			fmt.Printf("Error executing queries without pooling: %v\n", err)
			continue
		}
		fmt.Printf("Time taken without pooling for %d queries: %v\n", count, time_taken)
	}

	// 2. Run with Connection Pooling
	fmt.Println()
	fmt.Println("Running with Connection Pooling...")

	for _, count := range queriesCounts {
		fmt.Printf("Executing %d queries with connection pooling...\n", count)
		time_taken, err := runWithPooling(count)

		if err != nil {
			fmt.Printf("Error executing queries with pooling: %v\n", err)
			continue
		}
		fmt.Printf("Time taken with pooling for %d queries: %v\n", count, time_taken)
	}
}
