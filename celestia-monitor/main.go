package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	client "github.com/celestiaorg/celestia-openrpc"
	"github.com/celestiaorg/celestia-openrpc/types/blob"
	"github.com/celestiaorg/celestia-openrpc/types/share"
	"github.com/celestiaorg/rsmt2d"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	url := "ws://localhost:26658" // use ws for subscribing to headers
	token := ""                   // Use the appropriate token or skip-auth flag

	go func() {
		err := SubscribeHeaders(ctx, url, token)
		if err != nil {
			log.Fatalf("🚧 Error subscribing to headers: %v", err)
		}
	}()

	// Submit the initial blob
	err := SubmitBlob(ctx, url, token)
	if err != nil {
		log.Fatalf("🚧 Error submitting blob: %v", err)
	}

	select {}
}

// SubmitBlob submits a blob containing "Hello, World!" to the 0xC0DE namespace. It uses the default signer on the running node.
func SubmitBlob(ctx context.Context, url string, token string) error {
	client, err := client.NewClient(ctx, url, token)
	if err != nil {
		return err
	}

	// let's post to 0xC0DE namespace
	namespace, err := share.NewBlobNamespaceV0([]byte{0xC0, 0xDE})
	if err != nil {
		return err
	}

	// create a blob
	blobData := fmt.Sprintf("Hello, World! %v", time.Now())

	// Add 500kb of junk bytes
	//junkBytes := strings.Repeat("x", 500*1024)
	//blobData += junkBytes

	helloWorldBlob, err := blob.NewBlobV0(namespace, []byte(blobData))
	if err != nil {
		return err
	}

	// submit the blob to the network
	height, err := client.Blob.Submit(ctx, []*blob.Blob{helloWorldBlob}, blob.DefaultGasPrice())
	if err != nil {
		return err
	}

	fmt.Printf("🟢 Blob was included at height %d\n", height)

	// fetch the blob back from the network
	retrievedBlobs, err := client.Blob.GetAll(ctx, height, []share.Namespace{namespace})
	if err != nil {
		return err
	}

	fmt.Printf("🧐 Blobs are equal? %v\n", bytes.Equal(helloWorldBlob.Commitment, retrievedBlobs[0].Commitment))
	return nil
}

// SubscribeHeaders subscribes to new headers, fetches all blobs in the 0xC0DE namespace, retrieves the EDS at the new header's height, and submits new blobs.
func SubscribeHeaders(ctx context.Context, url string, token string) error {
	client, err := client.NewClient(ctx, url, token)
	if err != nil {
		return err
	}

	// create a namespace to filter blobs with
	namespace, err := share.NewBlobNamespaceV0([]byte{0xC0, 0xDE})
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
			fmt.Printf("🧊 New header received: Height %d, Hash %x\n", header.Height(), header.Hash())

			// fetch all blobs at the height of the new header
			blobs, err := client.Blob.GetAll(context.TODO(), header.Height(), []share.Namespace{namespace})
			if err != nil {
				fmt.Printf("🚧 Error fetching blobs: %v\n", err)
			} else {
				fmt.Printf("🟣 Found %d blobs at height %d in 0xC0DE namespace\n", len(blobs), header.Height())
			}

			// fetch the EDS at the height of the new header
			eds, err := GetEDS(ctx, url, token, header.Height())
			if err != nil {
				fmt.Printf("🚧 Error fetching EDS: %v\n", err)
			} else {
				fmt.Printf("🟩 EDS fetched at height %d: %v\n", header.Height(), eds)
			}

			// submit a new blob
			err = SubmitBlob(ctx, url, token)
			if err != nil {
				fmt.Printf("🚧 Error submitting new blob: %v\n", err)
			} else {
				fmt.Println("✅ New blob submitted and verified successfully")
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
