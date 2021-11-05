package frontend

import "github.com/labstack/echo/v4"

func NewHandler(e *echo.Echo) {
	e.Static("/src", "src/")

	e.File("/", "src/html/average_kdr.html")
	e.File("/kdr", "src/html/average_kdr.html")
	e.File("/dpm", "src/html/average_dpm.html")
	e.File("/hpm", "src/html/average_hpm.html")
}
