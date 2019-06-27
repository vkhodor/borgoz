package main

import (
	"fmt"
	"github.com/labstack/echo"
)

func backupNotOlderThen(c echo.Context) error {

	repo := c.Param("repo")
	time := c.Param("time")

	fmt.Printf("repo: %s\n", repo)
	fmt.Printf("minutes: %s\n", time)

	return nil
}

func main() {
	app := echo.New()
	app.GET("/check/backup/:repo/not_older_then/:time", backupNotOlderThen)

	app.Logger.Fatal(app.Start(":1323"))
}
