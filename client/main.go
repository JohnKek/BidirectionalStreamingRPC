package main

import (
	api "chat/api/grpc"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"math/rand"
)

func main() {
	host := "localhost"
	port := "8080"
	fieldSize := int32(rand.Intn(6) + 1)
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%s", host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := api.NewBattleshipServiceClient(conn)
	stream, err := client.Game(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	err = stream.Send(&api.Request{Data: &api.Request_Start{Start: &api.StartGame{FieldSize: fieldSize}}})
	if err != nil {
		log.Fatal(err)
	}
outerLoop:
	for i := int32(0); i < fieldSize; i++ {
		for j := int32(0); j < fieldSize; j++ {
			err = stream.Send(&api.Request{Data: &api.Request_Coordinate{Coordinate: &api.AttackCoordinate{X: i, Y: j}}})
			if err != nil {
				if err == io.EOF {
					fmt.Println("End of stream")
					break outerLoop
				}
				log.Fatal(err)
			}
			resp, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					fmt.Println("End of stream")
					break outerLoop
				}
				log.Fatal(err)
			}
			log.Println("resp.Status")
			log.Println(resp.Status)
		}
	}
}
