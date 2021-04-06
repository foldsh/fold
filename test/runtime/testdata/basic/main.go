package main

import (
	"fmt"
	"os"

	"github.com/foldsh/fold/sdks/go/fold"
)

func main() {
	svc := fold.NewService()
	svc.Get("/hello/:name", func(req *fold.Request, res *fold.Response) {
		res.StatusCode = 200
		res.Body = map[string]interface{}{
			"greeting": fmt.Sprintf("Hello, %s!", req.PathParams["name"]),
		}
	})
	svc.Get("/crash", func(req *fold.Request, res *fold.Response) {
		os.Exit(1)
	})
	svc.Start()
}
