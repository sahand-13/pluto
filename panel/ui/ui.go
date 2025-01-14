package ui

import (
	"embed"
	"pluto"

	"github.com/labstack/echo/v4"
)

var (
	//go:embed all:out
	UI          embed.FS
	UIPublicDir = echo.MustSubFS(UI, "out")

	//go:embed out/index.html
	PageHome embed.FS

	//go:embed out/404.html
	PageNotFound embed.FS

	//go:embed out/auth.html
	PageAuth embed.FS

	//go:embed out/logsview.html
	PageLogsView embed.FS

	//go:embed out/pipelines.html
	PagePipelines embed.FS

	//go:embed out/pipelines/create.html
	PagePipelinesCreate embed.FS

	//go:embed out/pipelines/edit.html
	PagePipelinesEdit embed.FS
)

func init() {
	p := pluto.FindHTTPHost(pluto.PanelSubdomain)
	p.StaticFS("/", UIPublicDir)
	p.FileFS("/", "index.html", echo.MustSubFS(PageHome, "out"))
	p.FileFS("/404", "404.html", echo.MustSubFS(PageNotFound, "out"))
	p.FileFS("/auth", "auth.html", echo.MustSubFS(PageAuth, "out"))
	p.FileFS("/logsview", "logsview.html", echo.MustSubFS(PageLogsView, "out"))
	p.FileFS("/pipelines", "pipelines.html", echo.MustSubFS(PagePipelines, "out"))
	p.FileFS("/pipelines/create", "pipelines/create.html", echo.MustSubFS(PagePipelinesCreate, "out"))
	p.FileFS("/pipelines/edit", "pipelines/edit.html", echo.MustSubFS(PagePipelinesEdit, "out"))
}
