package controllers

import (
	"fmt"
	"net/http"
	"simpleblog/models"
	"strconv"

	"github.com/labstack/echo"
)

type PostController struct {
}

func (c PostController) Init(g *echo.Group) {
	g.GET("/", c.Index)
	g.GET("/:id", c.GetById)
	g.GET("/create", func(c echo.Context) error {
		return c.Render(http.StatusOK, "Create.html", nil)
	})
	g.POST("/create", c.Create)
	g.GET("/:id/delete", c.Delete)
}

func (PostController) Index(c echo.Context) error {
	ps, err := models.Post{}.Index(c.Request().Context())
	if err != nil {
		c.String(http.StatusOK, "read DB fail")
	}

	return c.Render(http.StatusOK, "Index.html", ps)
}

func (PostController) GetById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusOK, fmt.Sprintf("%s is not integer", c.Param("id")))
	}
	p, err := models.Post{}.GetById(c.Request().Context(), id)
	if err != nil {
		c.String(http.StatusOK, fmt.Sprintf("%d post isn't exist", id))
	}
	return c.Render(http.StatusOK, "Post.html", p)
}

func (PostController) Create(c echo.Context) error {
	p := new(models.Post)
	p.Title = c.FormValue("title")
	p.Body = c.FormValue("body")

	err := p.Create(c.Request().Context())
	if err != nil {
		return c.String(http.StatusOK, err.Error())
	}
	return c.Redirect(http.StatusMovedPermanently, "/posts/")
}

func (PostController) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusOK, fmt.Sprintf("%s is not integer", c.Param("id")))
	}

	models.Post{}.Delete(c.Request().Context(), id)
	return c.Redirect(http.StatusMovedPermanently, "/posts/")
}
