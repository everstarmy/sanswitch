package sanswitch_test

import (
	"context"
	"fmt"

	"github.com/everstarmy/sanswitch"
)

func ExampleNew() {
	client, err := sanswitch.New("switch.example",
		sanswitch.WithHTTP(),
		sanswitch.WithRetry(0),
	)
	if err != nil {
		return
	}
	defer func() { _ = client.Close() }()

	fmt.Println(client.Timeout())
	fmt.Println(client.RetryCount())
	// Output:
	// 30s
	// 0
}

func ExampleClient_consumerInterface() {
	var reader interface {
		Ports(context.Context) ([]sanswitch.Port, error)
	} = (*sanswitch.Client)(nil)

	fmt.Printf("%T\n", reader)
	// Output:
	// *sanswitch.Client
}
