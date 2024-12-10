package main

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"
)

type cacheEntry struct {
	IP     string
	Expiry time.Time
}

var (
	cache      = make(map[string]cacheEntry)
	cacheMutex sync.Mutex
)

func main() {
	listner, err := net.Listen("tcp", ":9000")

	if err != nil {
		fmt.Println("Error in starting dns resolver", err)
		return
	}

	defer listner.Close()
	fmt.Println("dns resolver is running on 9000")

	for {
		conn, err := listner.Accept()

		if err != nil {
			fmt.Println("Error with connection the client", err)
			continue
		}

		go handleClient(conn)
	}

}

func handleClient(conn net.Conn) {
	defer conn.Close()
	var query string
	err := json.NewDecoder(conn).Decode(&query)
	if err != nil {
		fmt.Println("Error decoding client query:", err)
	}

	//chache mutex and search in cache
	cacheMutex.Lock()
	entry, found := cache[query]
	if found && entry.Expiry.After(time.Now()) {
		cacheMutex.Unlock()
		json.NewEncoder(conn).Encode(entry.IP)
		return
	}
	cacheMutex.Unlock()

	//query the authoritative server
	ip, err := queryAuthoritative(query)
	if err != nil {
		fmt.Printf("Error querying authritative server")
		return
	}

	//update the cache
	cacheMutex.Lock()
	cache[query] = cacheEntry{
		IP:     ip,
		Expiry: time.Now().Add(30 * time.Second),
	}
	cacheMutex.Unlock()

}

func queryAuthoritative(domain string) (string, error) {
	conn, err := net.Dial("tcp", "localhost:9001")

	if err != nil {
		fmt.Print("err", err)
	}

	defer conn.Close()

	json.NewEncoder(conn).Encode(domain)
	var ip string
	err = json.NewDecoder(conn).Decode(&ip)
	return ip, err
}
