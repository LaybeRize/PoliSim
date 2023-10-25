package extraction

import "PoliSim/data/database"

func CreateArticle(article *database.Article) error {
	return database.DB.Create(&article).Error
}
