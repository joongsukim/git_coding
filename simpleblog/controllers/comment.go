package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"simpleblog/models"
	"github.com/labstack/echo"
)

type CommentController struct {
}

func (c CommentController) Init(g *echo.Group) {
	g.POST("/create", c.Create)
	g.GET("/delete/:cid", c.Delete)
}

func (CommentController) Create(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusOK, fmt.Sprintf("%s is not integer", c.Param("id")))
	}

	cm := new(models.Comment)
	cm.Body = c.FormValue("body")
	cm.PostId = id

	err = cm.Create(c.Request().Context())
	if err != nil {
		c.String(http.StatusOK, "can't craete comment")
	}
	return c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/posts/%d", id))
}

func (CommentController) Delete(c echo.Context) error {
	cid, err := strconv.Atoi(c.Param("cid"))
	if err != nil {
		c.String(http.StatusOK, fmt.Sprintf("%s is not integer", c.Param("cid")))
	}
	err = models.Comment{}.Delete(c.Request().Context(), cid)
	return c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("/posts/%s", c.Param("id")))
}
