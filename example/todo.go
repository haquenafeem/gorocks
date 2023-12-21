package main

import (
	"fmt"
	"net/http"

	"github.com/haquenafeem/gorocks"
)

type todo struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var todos = []todo{
	{
		Id:        "x1",
		Title:     "Title 1",
		Completed: false,
	},
	{
		Id:        "x2",
		Title:     "Title 1",
		Completed: false,
	},
	{
		Id:        "x3",
		Title:     "Title 1",
		Completed: false,
	},
}

func runTodo() {
	app := gorocks.New()

	app.Use(TimeTaken)
	// app.Use(gorocks.ApplicationJSON)
	// app.Use(gorocks.BasicAuth("user", "xp1"))

	app.Get("/todos", gorocks.ResponseWithHeaders(map[string]string{
		"key1":         "val1",
		"Content-Type": "application/json",
	})(getAll))

	app.Get("/todos/:id", getSingle)
	app.Post("/todos", post)
	app.Delete("/todos/:id", delete)
	app.Put("/todos/:id", update)
	app.Get("/index", func(app *gorocks.App) {
		app.SetHeader("Content-Type", "text/html")
		app.ResponseWriter().Write([]byte("<h1>Hello<h2>"))
	})

	app.Get("/", func(app *gorocks.App) {
		app.JSON(400, map[string]interface{}{
			"key": "value",
		})
	})

	app.PrintRoutes()
	fmt.Println("app starting at :3002")
	app.Run(":3002")
}

func update(app *gorocks.App) {
	id := app.Param("id")
	if id == "" {
		app.JSON(
			http.StatusInternalServerError,
			map[string]interface{}{
				"err": "id not found",
			},
		)

		return
	}
	fmt.Println("id==>", id)
	var td todo
	if err := app.BindJson(&td); err != nil {
		app.JSON(http.StatusInternalServerError, map[string]interface{}{
			"err": err.Error(),
		})

		return
	}

	for i := range todos {
		if todos[i].Id == id {
			todos[i].Completed = td.Completed
			todos[i].Title = td.Title
			app.JSON(
				http.StatusInternalServerError,
				map[string]interface{}{
					"is_success": true,
				},
			)

			return
		}
	}

	app.JSON(
		http.StatusNotFound,
		map[string]interface{}{
			"err": "resource not found",
		},
	)
}

func delete(app *gorocks.App) {
	id := app.Param("id")
	if id == "" {
		app.JSON(
			http.StatusInternalServerError,
			map[string]interface{}{
				"err": "id not found",
			},
		)

		return
	}
	fmt.Println("id==>", id)
	for i := range todos {
		if todos[i].Id == id {
			todos = append(todos[:i], todos[i+1:]...)
			app.JSON(
				http.StatusOK,
				map[string]interface{}{
					"is_success": true,
				},
			)

			return
		}
	}

	app.JSON(
		http.StatusNotFound,
		map[string]interface{}{
			"err": "resource not found",
		},
	)
}

func post(app *gorocks.App) {
	var td todo
	if err := app.BindJson(&td); err != nil {
		app.JSON(http.StatusInternalServerError, map[string]interface{}{
			"err": err.Error(),
		})

		return
	}

	todos = append(todos, td)
	app.JSON(http.StatusOK, map[string]interface{}{
		"is_success": true,
	})
}

func getSingle(app *gorocks.App) {
	id := app.Param("id")
	if id == "" {
		app.JSON(
			http.StatusInternalServerError,
			map[string]interface{}{
				"err": "id not found",
			},
		)

		return
	}
	fmt.Println("id==>", id)
	for i := range todos {
		if todos[i].Id == id {
			app.JSON(
				http.StatusOK,
				map[string]interface{}{
					"todo": todos[i],
				},
			)

			return
		}
	}

	app.JSON(
		http.StatusNotFound,
		map[string]interface{}{
			"err": "resource not found",
		},
	)
}

func getAll(app *gorocks.App) {
	// app.SetHeader("Content-Type", "application/json")
	app.JSON(
		http.StatusOK,
		map[string]interface{}{
			"data": todos,
		},
	)
}
