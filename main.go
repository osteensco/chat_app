package main

import (
	"github.com/aofei/air"
)

func main() {
	air.Default.GET("/", func(req *air.Request, res *air.Response) error {
		return res.WriteString("Hello world")
	})
	air.Default.Serve()
}
