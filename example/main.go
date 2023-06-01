package main

import (
	"fmt"

	go253 "github.com/ching-xia/cl253-go"
)

func main() {
	client, err := go253.NewClient(
		go253.WithAccount("account"),
		go253.WithPassword("password"),
		go253.WithNodeType(go253.NodeShanghai),
	)
	if err != nil {
		panic(err)
	}
	defer client.Close()
	balance, err := client.Balance()
	if err != nil {
		panic(err)
	}
	fmt.Println(balance)
	// single send
	msg, err := go253.NewMessage(
		go253.WithMessage("mobile", "msg", "params1, params2"),
		go253.WithUid("uid"),
		go253.WithSenderID("senderID"),
	)
	if err != nil {
		panic(err)
	}
	record := client.SingleMessage(msg)
	fmt.Println(record)
	// batch send
	in := client.In()
	out := client.Out()
	go func() {
		for i := 0; i < 10; i++ {
			msg, err := go253.NewMessage(
				go253.WithMessage("mobile", "msg", "params1, params2"),
				go253.WithUid("uid"),
				go253.WithSenderID("senderID"),
			)
			if err != nil {
				panic(err)
			}
			in <- msg
		}
	}()
	go func() {
		for record := range out {
			fmt.Println(record)
			if record.Error != nil {
				panic(record.Error)
			}
		}
	}()
}
