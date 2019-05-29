package main

import (
	"fmt"
	"github.com/wpxun/gomicro/service"
	"log"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	redis := service.GetRedis()
	val, err := redis.Incr("count").Result()
	if err != nil {
		panic(err)
	}
	host, _ := os.Hostname()
	fmt.Fprintln(w, "hello world "+ host +", visitors = " + strconv.FormatInt(val, 10) )
}


type Args struct {
	A, B int
}

func MathHandler(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	a, _ := strconv.Atoi(r.Form["a"][0])
	b, _ := strconv.Atoi(r.Form["b"][0])

	client, err := rpc.DialHTTP("tcp", service.Conf.MathRPC.Address)
	if err != nil {
		log.Fatal("Dial error:", err)
	}

	args := Args{a,b}
	var reply int

	err = client.Call("Arith.Math", args, &reply)

	if err != nil {
		log.Fatal("arith error", err)
	}
	fmt.Fprintf(w, "Arith: %d*%d=%d\n ", args.A, args.B, reply)
}

func main()  {
	http.Handle("/index", http.HandlerFunc(IndexHandler))
	http.Handle("/math", http.HandlerFunc(MathHandler))
	http.ListenAndServe(":80", nil)
}