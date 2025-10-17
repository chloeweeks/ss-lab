package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"net/rpc"
	"os"
	"time"

	"uk.ac.bris.cs/distributed2/secretstrings/stubs"
)

/** Super-Secret `reversing a string' method we can't allow clients to see. **/
func ReverseString(s string, i int) string {
	time.Sleep(time.Duration(rand.Intn(i)) * time.Second)
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

type SecretStringOperations struct{}

func (s *SecretStringOperations) Reverse(req stubs.Request, res *stubs.Response) (err error) {
	res.Message = ReverseString(req.Message, 10)
	return // not returning by directly altering a function
}
func (s *SecretStringOperations) FastReverse(req stubs.Request, res *stubs.Response) (err error) {
	res.Message = ReverseString(req.Message, 2) // the 2 is the random delay
	return
}
func main() {
	port := "8030" // declaring initial port if no port is given

	if len(os.Args) > 1 { // adds port numbers to port
		port = os.Args[1]
		fmt.Println("port being used: " + port)
	}

	pAddr := flag.String("port", ":"+port, "Port to listen on")

	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	rpc.Register(&SecretStringOperations{})

	listener, err := net.Listen("tcp", *pAddr)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()
	rpc.Accept(listener)

}
