package main

import (
	"context"
	"html/template"
	"io"
	"net/http"
	"simpleblog/controllers"
	"simpleblog/models"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const (
	//driver  = "mysql"
	//connect = "root:pw@/table?charset=utf8"
	driver  = "mysql"
	connect = "root:twin2k@/hello?charset=utf8"
)

func main() {
	db, err := xorm.NewEngine(driver, connect)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Sync(new(models.Post))
	err = db.Sync(new(models.Comment))

	e := echo.New()
	e.Use(ContextDB(db))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	controllers.PostController{}.Init(e.Group("/posts"))
	controllers.CommentController{}.Init(e.Group("/posts/:id/comment"))

	t := &Template{
		templates: template.Must(template.ParseGlob("./views/*.html")),
	}
	e.Renderer = t

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Home")
	})

	e.Logger.Fatal(e.Start(":1323"))
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func ContextDB(db *xorm.Engine) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			session := db.NewSession()
			defer session.Close()

			req := c.Request()
			c.SetRequest(req.WithContext(context.WithValue(req.Context(), "DB", session)))

			switch req.Method {
			case "POST", "PUT", "DELETE":
				if err := session.Begin(); err != nil {
					return echo.NewHTTPError(500, err.Error())
				}
				if err := next(c); err != nil {
					session.Rollback()
					return echo.NewHTTPError(500, err.Error())
				}
				if c.Response().Status >= 500 {
					session.Rollback()
					return nil
				}
				if err := session.Commit(); err != nil {
					return echo.NewHTTPError(500, err.Error())
				}
			default:
				if err := next(c); err != nil {
					return echo.NewHTTPError(500, err.Error())
				}
			}

			return nil
		}
	}
}
