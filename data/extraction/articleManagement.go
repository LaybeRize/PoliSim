package extraction

import "PoliSim/data/database"

func CreateArticle(article *database.Article) error {
	return database.DB.Create(&article).Error
}

func FindArticlesForPublicationUUID(uuid string) (*database.ArticleList, error) {
	list := &database.ArticleList{}
	err := database.DB.Where("publication = ?", uuid).Find(list).Error
	return list, err
}
