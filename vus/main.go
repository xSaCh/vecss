package main

import (
	"context"

	"github.com/xSaCh/vecss/vus/pkg/repositories"
)

func main() {

	fac := repositories.RepositoryFactory{}
	storage := fac.NewStorageRepository()

	_ = storage

	// storage.T()
	// res, err := storage.PutObject(context.TODO(), "bkt", "example.txt", 15*60)

	// res, err := storage.GenerateMultiPartPreSignedUrls(context.TODO(), []int{1, 2, 3, 4, 5})
	// if err != nil {
	// 	panic(err)
	// }
	// for _, r := range res {

	// 	fmt.Printf("%v\n", r.URL)
	// }

	tags := []string{"3204d31ca09c2493a1dd808f5eac79c7",
		"597295474dbbaca27e9df459a6774b68",
		"e6659f5e7b08206c2eaa9897ad334a0e",
		"eef2f4f9815fe629edea6377c2243069",
		"0b3c18e6c1e7cd136bcf45936ee0f7cf"}

	err := storage.CombineMultiPartUploads(context.TODO(), "62b13101-7a98-49f8-9ef6-2462961d280c", tags, []int{1, 2, 3, 4, 5})
	if err != nil {
		panic(err)
	}

}
