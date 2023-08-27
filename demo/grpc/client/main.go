package main

import (
	"context"
	"fmt"
	rpc_impl "github.com/mats9693/study/demo/grpc/proto/impl"
	"google.golang.org/grpc"
	"log"
	"strconv"
	"time"
)

func main() {
	conn, err := grpc.Dial(":9693", grpc.WithInsecure())
	if err != nil {
		log.Fatalln("dial server failed, error:", err)
	}
	defer conn.Close()

	client := rpc_impl.NewIUserClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	for i := 0; i < 5; i++ {
		res, err := client.Login(ctx, &rpc_impl.LoginReq{
			UserName: "mario" + strconv.Itoa(i),
			Password: "mario" + strconv.Itoa(i),
		})
		if err != nil {
			log.Fatalln("send req failed, error:", err)
		}

		fmt.Println("res:", res.IsSuccess, res.Err)

		time.Sleep(time.Second)
	}
}
