package main

import (
	"net/http"
	"net/rpc"
)

type Args struct {
	A, B int
}

type Arith int

func (t *Arith) Math(args * Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func main() {
	arith := new(Arith)
	rpc.Register(arith)
	rpc.HandleHTTP()

	http.ListenAndServe(":1234", nil)
}