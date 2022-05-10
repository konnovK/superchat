package repository

import (
	"github.com/konnovK/superchat/internal/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TagRepository interface {
	Find(conditions *entity.Tag) (entity.Tags, error)
	FindAll() (entity.Tags, error)
	Create(target *entity.Tag) error
	Update(conditions *entity.Tag, target *entity.Tag) error
	Delete(target *entity.Tag) error
}

type Tag struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) TagRepository {
	return &Tag{
		db: db,
	}
}

func (t *Tag) Find(conditions *entity.Tag) (entity.Tags, error) {
	tags := []entity.Tag{}

	queryResult := t.db.Where(conditions).Find(&tags)
	if queryResult.Error != nil {
		return tags, queryResult.Error
	}

	return tags, nil
}

func (t *Tag) FindAll() (entity.Tags, error) {
	return t.Find(&entity.Tag{})
}

func (t *Tag) Create(target *entity.Tag) error {
	queryResult := t.db.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "title"}}, UpdateAll: true}).Omit("ID").Create(&target)
	if queryResult.Error != nil {
		return queryResult.Error
	}
	return nil
}

func (t *Tag) Update(conditions *entity.Tag, target *entity.Tag) error {
	queryResult := t.db.Where(&conditions).Updates(target)
	if queryResult.Error != nil {
		return queryResult.Error
	}

	return nil
}

func (t *Tag) Delete(target *entity.Tag) error {
	queryResult := t.db.Where(&target).Delete(&target)
	if queryResult.Error != nil {
		return queryResult.Error
	}

	return nil
}
