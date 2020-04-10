package models

import (
	"context"
	"time"

	"github.com/go-xorm/xorm"
)

type Comment struct {
	Id   int64
	Body string `xorm:"text"`
	//Commenter string
	PostId    int
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
}

func (Comment) GetById(c context.Context, id int64) ([]Comment, error) {
	var cs []Comment

	db := c.Value("DB").(*xorm.Session)
	err := db.Where("post_id=?", id).Find(&cs)
	if err != nil {
		return nil, err
	}

	return cs, nil
}

func (cm *Comment) Create(c context.Context) error {
	db := c.Value("DB").(*xorm.Session)
	_, err := db.Insert(cm)

	return err
}

func (Comment) Delete(c context.Context, id int) error {
	db := c.Value("DB").(*xorm.Session)
	_, err := db.ID(id).Delete(&Comment{})
	return err
}
