package extraction

import (
	"PoliSim/data/database"
	"fmt"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"testing"
)

func TestInterfacing(t *testing.T) {
	var err error
	database.DB, err = gorm.Open(postgres.Open(fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_ADDRESS"),
		os.Getenv("DB_NAME"))), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	assert.Nil(t, err)
	doc := SpecificDocument{}
	err = FindDocumentInterfaceByUUID("5f245374-9d9c-4b50-882f-8663cbc02331", &doc)
	assert.Nil(t, err)
	assert.Equal(t, "5f245374-9d9c-4b50-882f-8663cbc02331", doc.UUID)
	assert.Equal(t, "Testabstimmungsformular Nr. 2", doc.Title)
	assert.Equal(t, 2, len(doc.Viewer))
	assert.Equal(t, int64(1), doc.Viewer[0].ID)
	assert.Equal(t, "Lennard Kirchenberg", doc.Viewer[0].DisplayName)
	assert.Equal(t, int64(4), doc.Viewer[1].ID)
	assert.Equal(t, "test1234", doc.Viewer[1].DisplayName)
}
