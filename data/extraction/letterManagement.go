package extraction

import (
	"PoliSim/data/database"
	"errors"
	"gorm.io/gorm"
	"time"
)

func CreateLetter(letter *database.Letter) error {
	return database.DB.Create(&letter).Error
}

func GetByID(uuid string) (*database.Letter, error) {
	letter := &database.Letter{}
	err := database.DB.Where("uuid = ?", uuid).First(letter).Error
	return letter, err
}

func GetLettersAfter(publicationUUID string, amount int, accountId int64) (*database.ExtendedLetterList, bool, error) {
	letterList := &database.ExtendedLetterList{}
	exists, err := getLetters(letterList, publicationUUID, func(pub *database.Letter) *gorm.DB {
		return getBasicLetterQuery(pub.UUID, amount, accountId).Where("written < ?", pub.Written).Order("written desc")
	})
	return letterList, exists, err
}

func GetLettersBefore(publicationUUID string, amount int, accountId int64) (*database.ExtendedLetterList, bool, error) {
	letterList := &database.ExtendedLetterList{}
	exists, err := getLetters(letterList, publicationUUID, func(pub *database.Letter) *gorm.DB {
		return database.DB.Select("*").Table("(?) as X", getBasicLetterQuery(pub.UUID, amount, accountId).Where("written > ?", pub.Written).Order("written")).Order("X.written desc")
	})
	return letterList, exists, err
}

func GetModMailsAfter(publicationUUID string, amount int) (*database.ExtendedLetterList, bool, error) {
	letterList := &database.ExtendedLetterList{}
	exists, err := getLetters(letterList, publicationUUID, func(pub *database.Letter) *gorm.DB {
		return getBasicModmailQuery(pub.UUID, amount).Where("written < ?", pub.Written).Order("written desc")
	})
	return letterList, exists, err
}

func GetModMailsBefore(publicationUUID string, amount int) (*database.ExtendedLetterList, bool, error) {
	letterList := &database.ExtendedLetterList{}
	exists, err := getLetters(letterList, publicationUUID, func(pub *database.Letter) *gorm.DB {
		return database.DB.Select("*").Table("(?) as X", getBasicModmailQuery(pub.UUID, amount).Where("written > ?", pub.Written).Order("written")).Order("X.written desc")
	})
	return letterList, exists, err
}

func getLetters(letterList *database.ExtendedLetterList, publicationUUID string, query func(pub *database.Letter) *gorm.DB) (bool, error) {
	exists := true
	pub, err := GetByID(publicationUUID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		exists = false
		pub.Written = time.Now().UTC()
	} else if err != nil {
		return exists, err
	}
	err = query(pub).Find(letterList).Error
	return exists, err
}

func getBasicLetterQuery(uuid string, amount int, accountID int64) *gorm.DB {
	return database.DB.Joins("JOIN letter_account ON letters.uuid = letter_account.uuid").
		Where("letter_account.id = ?", accountID).Select("*").Where("letters.uuid != ?", uuid).Limit(amount).Table("letters")
}

func getBasicModmailQuery(uuid string, amount int) *gorm.DB {
	return database.DB.Where("letters.uuid != ? AND mod_message = true", uuid).Limit(amount).Table("letters")
}

func GetLetterByIDOnlyWithAccount(uuid string, accountID int64, isMod bool) (*database.Letter, error) {
	letter := &database.Letter{}
	err := database.DB.Preload("Viewer").Joins("JOIN letter_account ON letters.uuid = letter_account.uuid").
		Where("(letter_account.id = ? AND letters.uuid = ?) OR (? = true AND mod_message = true AND letters.uuid = ?)", accountID, uuid, isMod, uuid).First(letter).Error
	return letter, err
}

func SetLetterAsRead(uuid string, accountID int64) {
	_ = database.DB.Raw("UPDATE letter_account SET read = true WHERE uuid = ? and id = ?", uuid, accountID).Scan(&database.LetterAccount{}).Error
}

func SetAllLetterAsReadForAccount(accountID int64) error {
	return database.DB.Raw("UPDATE letter_account SET read = true WHERE read = false and id = ?", accountID).Scan(&database.LetterAccount{}).Error
}

func UpdateLetter(letter *database.Letter) error {
	return database.DB.Where("uuid = ?", letter.UUID).Updates(letter).Error
}

func GetCountOfLetters(accountID int64) (int64, error) {
	counter := int64(0)
	err := database.DB.Raw("SELECT accounts.id FROM accounts LEFT JOIN letter_account on accounts.id = letter_account.id WHERE (accounts.id = ? OR linked = ?) AND read = false",
		accountID, accountID).Count(&counter).Error
	return counter, err
}
