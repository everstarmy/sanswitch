package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/everstarmy/sanswitch"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	endpoint := os.Getenv("SAN_ENDPOINT")
	credentials := sanswitch.Credentials{
		Username: os.Getenv("SAN_USERNAME"),
		Password: os.Getenv("SAN_PASSWORD"),
	}
	if endpoint == "" || credentials.Username == "" || credentials.Password == "" {
		return fmt.Errorf("SAN_ENDPOINT, SAN_USERNAME, and SAN_PASSWORD are required")
	}

	client, err := sanswitch.Open(ctx, endpoint, credentials)
	if err != nil {
		return fmt.Errorf("open switch: %w", err)
	}
	defer func() { _ = client.Close() }()
	defer func() {
		logoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := client.Logout(logoutCtx); err != nil {
			log.Printf("logout: %v", err)
		}
	}()

	switches, err := client.FabricSwitches(ctx)
	if err != nil {
		return fmt.Errorf("list fabric switches: %w", err)
	}
	for _, sw := range switches {
		fmt.Printf("%s (domain %d, IP %s)\n", sw.Name, sw.DomainID, sw.IPAddress)
	}

	ports, err := client.Ports(ctx)
	if err != nil {
		return fmt.Errorf("list ports: %w", err)
	}
	for _, port := range ports {
		fmt.Printf("port %s: %s, speed %s\n", port.Name, port.OperationalStatusString, port.Speed)
	}
	return nil
}
