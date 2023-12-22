# GOROCKS
A simple web framework.

## Notes
- Experimental
- Needs Work

## Example 
```go
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/haquenafeem/gorocks"
)

func main() {
	app := gorocks.New()
	// If you want your response to be of json type; "Content-Type":"application/json"
	app.Use(gorocks.ApplicationJSON)
	app.Use(TimeTaken)

	app.Get("/", func(a *gorocks.App) {
		app.JSON(http.StatusOK, gorocks.Map{
			"message": "ok",
		})
	})

	if err := app.Run(":9000"); err != nil {
		panic(err)
	}
}

// TimeTaken
// Middleware
func TimeTaken(next gorocks.HandlerFunc) gorocks.HandlerFunc {
	return func(app *gorocks.App) {
		now := time.Now()
		next(app)
		passed := time.Since(now)
		out := fmt.Sprintf("time taken to complete %s is %v", app.HttpRequest().URL.Path, passed.Seconds())
		fmt.Println(out)
	}
}

// test by curl localhost:9000
// output :
// {
//         "message": "ok"
// }
```
<p>
  More Example <a href="https://github.com/haquenafeem/gorocks/tree/main/example">Here</a>
</p>