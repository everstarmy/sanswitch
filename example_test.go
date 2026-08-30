package san_test

import (
	"context"
	"fmt"
	"os"

	san "github.com/everstarmy/sanswitch"
)

func ExampleNewClient() {
	client := san.NewClient("switch.example", "admin", os.Getenv("SAN_PASSWORD"),
		san.WithHTTP(),
		san.WithRetry(0),
	)
	defer func() { _ = client.Close() }()

	fmt.Println(client.Timeout())
	fmt.Println(client.RetryCount())
	// Output:
	// 30s
	// 0
}

func ExampleSwitchReader() {
	var reader interface {
		GetPortsWithContext(context.Context) ([]san.PortInfo, error)
	} = (*san.SANSwitch)(nil)

	fmt.Printf("%T\n", reader)
	// Output:
	// *san.SANSwitch
}
