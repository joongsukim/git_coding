package models

import (
	"context"
	"time"

	"github.com/go-xorm/xorm"
)

type Post struct {
	Id        int64
	Title     string
	Body      string `xorm:"text"`
	Comments  []Comment
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
}

func (Post) Index(c context.Context) ([]Post, error) {
	var ps []Post

	db := c.Value("DB").(*xorm.Session)
	rows, err := db.Rows(&Post{})
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		p := Post{}
		err = rows.Scan(&p)
		if err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}

	return ps, nil
}

func (Post) GetById(c context.Context, id int) (*Post, error) {
	p := new(Post)

	db := c.Value("DB").(*xorm.Session)
	ok, err := db.Id(id).Get(p)
	if err != nil {
		return nil, err
	} else if !ok {
		return nil, nil
	}
	p.Comments, err = Comment{}.GetById(c, p.Id)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Post) Create(c context.Context) error {
	db := c.Value("DB").(*xorm.Session)
	_, err := db.Insert(p)
	return err
}

func (Post) Delete(c context.Context, id int) error {
	db := c.Value("DB").(*xorm.Session)
	_, err := db.ID(id).Delete(&Post{})
	return err
}
