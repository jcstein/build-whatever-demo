package main

import (
	"context"
	"fmt"
	"log"

	client "github.com/celestiaorg/celestia-openrpc"
	"github.com/celestiaorg/celestia-openrpc/types/share"
	"github.com/celestiaorg/rsmt2d"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	url := "ws://localhost:26658" // use ws for subscribing to headers
	token := ""                   // Use the appropriate token or skip-auth flag

	err := SubscribeHeaders(ctx, url, token)
	if err != nil {
		log.Fatalf("🚧 Error subscribing to headers: %v", err)
	}
}

// SubscribeHeaders subscribes to new headers, fetches all blobs in the 0xdeadbeef namespace, and retrieves the EDS at the new header's height.
func SubscribeHeaders(ctx context.Context, url string, token string) error {
	client, err := client.NewClient(ctx, url, token)
	if err != nil {
		return err
	}

	// create a namespace to filter blobs with
	namespace, err := share.NewBlobNamespaceV0([]byte{0xDE, 0xAD, 0xBE, 0xEF})
	if err != nil {
		return err
	}

	// subscribe to new headers using a <-chan *header.ExtendedHeader channel
	headerChan, err := client.Header.Subscribe(ctx)
	if err != nil {
		return err
	}

	fmt.Println("📡 Subscribed to headers. Waiting for new headers...")

	for {
		select {
		case header := <-headerChan:
			fmt.Printf("📝 New header received: Height %d, Hash %x\n", header.Height(), header.Hash())

			// fetch all blobs at the height of the new header
			blobs, err := client.Blob.GetAll(context.TODO(), header.Height(), []share.Namespace{namespace})
			if err != nil {
				fmt.Printf("🚧 Error fetching blobs: %v\n", err)
			} else {
				fmt.Printf("🟣 Found %d blobs at height %d in 0xdeadbeef namespace\n", len(blobs), header.Height())
			}

			// fetch the EDS at the height of the new header
			eds, err := GetEDS(ctx, url, token, header.Height())
			if err != nil {
				fmt.Printf("🚧️ Error fetching EDS: %v\n", err)
			} else {
				fmt.Printf("🟩 EDS fetched at height %d: %v\n", header.Height(), eds)
			}
		case <-ctx.Done():
			return nil
		}
	}
}

// GetEDS fetches the EDS at the given height.
func GetEDS(ctx context.Context, url string, token string, height uint64) (*rsmt2d.ExtendedDataSquare, error) {
	client, err := client.NewClient(ctx, url, token)
	if err != nil {
		return nil, err
	}

	// First get the header of the block you want to fetch the EDS from
	header, err := client.Header.GetByHeight(ctx, height)
	if err != nil {
		return nil, err
	}

	// Fetch the EDS
	return client.Share.GetEDS(ctx, header)
}
