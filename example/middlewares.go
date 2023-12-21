package main

import (
	"fmt"
	"time"

	"github.com/haquenafeem/gorocks"
)

func TimeTaken(next gorocks.HandlerFunc) gorocks.HandlerFunc {
	return func(app *gorocks.App) {
		now := time.Now()
		next(app)
		passed := time.Since(now)
		out := fmt.Sprintf("time taken to complete %s is %v", app.HttpRequest().URL.Path, passed.Seconds())
		fmt.Println(out)
	}
}
