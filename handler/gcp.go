package main

import (
	"fmt"

	"cloud.google.com/go/storage"
	"golang.org/x/net/context"
)

const (
	BUCKET_ROOT string = "rebug-global-bucket"
	WASM_FOLDER string = "wasm"
)

type BucketWriter struct {
	cl     *storage.Client
	ctx    context.Context
	bucket string
	path   string
}

var bucketWriter *BucketWriter

func NewBucketWriter() (*BucketWriter, error) {
	br := new(BucketWriter)
	br.ctx = context.Background()
	client, err := storage.NewClient(br.ctx)
	if err != nil {
		return nil, err
	}
	br.cl = client
	br.bucket = BUCKET_ROOT
	br.path = WASM_FOLDER
	return br, nil
}

func (br *BucketWriter) WriteToBucket(file string, content []byte) error {
	wc := br.cl.Bucket(br.bucket).Object(fmt.Sprintf("%s/%s/main.wasm", br.path, file)).NewWriter(br.ctx)
	defer func() {
		err := wc.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()
	wc.ContentType = "application/wasm"
	if _, err := wc.Write(content); err != nil {
		return err
	}
	return nil
}
