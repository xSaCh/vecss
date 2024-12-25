package main

import (
	"context"
	"fmt"

	"github.com/xSaCh/vecss/vus/pkg/repositories"
)

func main() {

	fac := repositories.RepositoryFactory{}
	storage := fac.NewStorageRepository()

	_ = storage

	// storage.T()
	res, err := storage.PutObject(context.TODO(), "bkt", "example.txt", 15*60)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", res.URL)

}
