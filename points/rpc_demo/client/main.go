package main

import (
	"context"
	"fmt"
	"github.com/mats9693/utils/toy_server/rpc_demo/proto/impl"
	"google.golang.org/grpc"
	"log"
	"strconv"
	"time"
)

func main() {
	conn, err := grpc.Dial("localhost:9693", grpc.WithInsecure())
	if err != nil {
		log.Fatalln("dial server failed, error:", err)
	}
	defer conn.Close()

	client := rpc_impl.NewIUserClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	for i := 0; i < 5; i++ {
		res, err := client.Login(ctx, &rpc_impl.LoginReq{
			UserName: strconv.Itoa(i),
			Password: strconv.Itoa(i),
		})
		if err != nil {
			log.Fatalln("send req failed, error:", err)
		}

		fmt.Println("res:", res.GetIsSuccess(), res.GetErrMsg())

		time.Sleep(time.Second)
	}
}
