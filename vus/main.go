package main

import (
	"github.com/xSaCh/vecss/vus/pkg"
	"github.com/xSaCh/vecss/vus/pkg/repositories"
)

// func main() {

// 	fac := repositories.RepositoryFactory{}
// 	storage := fac.NewStorageRepository()

// 	_ = storage

// 	// storage.T()
// 	// res, err := storage.PutObject(context.TODO(), "bkt", "example.txt", 15*60)

// 	urls, err := storage.GenerateMultiPartPreSignedUrls(context.TODO(), "video.mp4", []int{1, 2, 3, 4, 5})
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("urls: %v\n", urls)
// 	for _, r := range urls.Urls {
// 		fmt.Printf("%v\n", r)
// 	}

// 	// tags := []string{"3204d31ca09c2493a1dd808f5eac79c7",
// 	// 	"597295474dbbaca27e9df459a6774b68",
// 	// 	"e6659f5e7b08206c2eaa9897ad334a0e",
// 	// 	"eef2f4f9815fe629edea6377c2243069",
// 	// 	"0b3c18e6c1e7cd136bcf45936ee0f7cf"}

// 	// err := storage.CombineMultiPartUploads(context.TODO(), "62b13101-7a98-49f8-9ef6-2462961d280c", tags, []int{1, 2, 3, 4, 5})
// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// }

func main() {
	fac := repositories.RepositoryFactory{}
	server := pkg.NewAPIServer(":8080", fac.NewStorageRepository())
	server.Run()

}
