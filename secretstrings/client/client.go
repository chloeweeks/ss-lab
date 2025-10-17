package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"sync"
	"time"

	"uk.ac.bris.cs/distributed2/secretstrings/stubs"
)

// initialising the servers that the client can use
// these are all to run: 'go run server.go 8030' for example
//var servers = []string{
//	"1172.31.20.29:8030",
//	"172.31.20.29:8040",
//	"172.31.20.29:8050",
//	"172.31.20.29:8060",
//}

func main() {

	// ************************** opening file and reading words **********************************
	file, er := os.Open("../wordlist") // opening word list file
	if er != nil {
		log.Fatal("failed to open wordlist", er) // error checking wordlist
	}
	defer file.Close()
	var words []string

	scanner := bufio.NewScanner(file) // going through each word in wordlist
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}
	if err := scanner.Err(); err != nil { // error checking
		log.Fatal(err)
	}

	// **********************************************************************************************

	// **************************** getting the list of available servers ***************************
	var servers []string

	for port := 8000; port <= 8090; port++ {
		address := "127.0.0.1:" + strconv.Itoa(port)
		conn, err := net.DialTimeout("tcp", address, 100*time.Millisecond)
		if err == nil {
			servers = append(servers, address)
			conn.Close()
		}
	}
	serverCount := len(servers) // initialising a count of servers
	fmt.Println("server count: " + strconv.Itoa(serverCount))

	// **********************************************************************************************

	wg := sync.WaitGroup{}

	for i, word := range words {
		server := servers[i%serverCount] // randomly allocating a server
		//fmt.Println("server with word " + server + word)

		wg.Add(1) // adding 1 to waitgroup to show go routine is about to be called
		go callReverser(server, word, &wg)

	}

	wg.Wait()
	fmt.Println("Done") // final done message

}

func callReverser(server string, word string, wg *sync.WaitGroup) {
	//fmt.Println("Calling reverser")
	defer wg.Done() // marking waitgroup as done

	client, err := rpc.Dial("tcp", server) // dialing tcp server

	if err != nil {
		fmt.Println("Skipping server with error", server, err)
		return
	}

	defer client.Close()

	// fmt.Println("Reversing word", word)                      // IT NEVER REACHES THIS HELP!!!!!!!!!!!!!!!!!
	req := stubs.Request{Message: word}                      // creating the request
	res := new(stubs.Response)                               // creating response
	err = client.Call(stubs.PremiumReverseHandler, req, res) // calling the reverse handler

	if err != nil {
		fmt.Println("RPC call failed on server with error:", server, err) // error checking
		return
	}
	//fmt.Printf("server %s responded: %s \n\n", server, res.Message) // printing which server responded with what
	fmt.Println(res.Message)
}
