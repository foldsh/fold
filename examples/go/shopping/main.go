package main

import (
	"github.com/foldsh/fold/sdks/go/fold"
)

var db map[string][]string

func main() {
	db = make(map[string][]string)
	svc := fold.NewService("hello")
	svc.Put("/items/:name", func(req *fold.Request, res *fold.Response) {
		n := req.PathParams["name"]
		item := req.Body["item"].(string)
		if value, ok := db[n]; ok {
			db[n] = append(value, item)
		} else {
			db[n] = []string{item}
		}
		res.StatusCode = 200
		res.Body = map[string]interface{}{
			"status": "success",
		}
	})
	svc.Get("/items/:name", func(req *fold.Request, res *fold.Response) {
		n := req.PathParams["name"]
		res.StatusCode = 200
		res.Body = map[string]interface{}{
			"items": db[n],
		}
	})
	svc.Start()
}
