package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	pgx "github.com/jackc/pgx/v4"
	echo "github.com/labstack/echo/v4"
	middleware "github.com/labstack/echo/v4/middleware"
)

func handler(c echo.Context) error {

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var name string
	var weight int64
	query := fmt.Sprintf("select name, weight from widgets where id='%s'", c.QueryParam("source"))
	err = conn.QueryRow(context.Background(), query).Scan(&name, &weight)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	return c.String(http.StatusOK, "ok")
	// return c.JSON(http.StatusOK, users)
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", handler)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
