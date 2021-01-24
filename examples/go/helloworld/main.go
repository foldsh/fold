package main

import (
	"fmt"

	"github.com/foldsh/fold/sdks/go/fold"
)

func main() {
	svc := fold.NewService("hello")
	svc.Get("/hello/:name", func(req *fold.Request, res *fold.Response) {
		res.StatusCode = 200
		res.Body = map[string]interface{}{
			"greeting": fmt.Sprintf("Hello, %s!", req.PathParams["name"]),
		}
	})
	svc.Start()
}
