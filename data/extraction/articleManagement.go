package extraction

import (
	"PoliSim/data/database"
	"errors"
)

func CreateArticle(article *database.Article) error {
	return database.DB.Create(&article).Error
}

func FindArticlesForPublicationUUID(uuid string) (*database.ArticleList, error) {
	list := &database.ArticleList{}
	err := database.DB.Where("publication = ?", uuid).Find(list).Error
	return list, err
}

func FindHiddenArticle(uuid string) (*database.Article, error) {
	art := &database.Article{}
	err := database.DB.Where("uuid = ?", uuid).First(art).Error
	if err != nil {
		return art, err
	}
	pub := &database.Publication{}
	err = database.DB.Where("uuid = ?", art.Publication).First(pub).Error
	//if the publication doesn't conform with the wished visiblity return a new error
	if err == nil && pub.Publicated {
		return art, errors.New("article already published")
	}
	return art, err
}

func DeleteArticle(article *database.Article) error {
	return database.DB.Delete(article).Error
}

func UpdatePublication(oldPub string, newPub string) error {
	return database.DB.Model(&database.Article{}).Select("publication").Where("publication = ?", oldPub).Updates(database.Article{Publication: newPub}).Error
}
