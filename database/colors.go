package database

import (
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"os"
)

type ColorPalette struct {
	Name       string `json:"Name,omitempty"`
	Background string `json:"Background,omitempty"`
	Text       string `json:"Text,omitempty"`
	Link       string `json:"Link,omitempty"`
}

var ColorPaletteMap map[string]ColorPalette

const folderPath = "./data"
const filePath = folderPath + "/colors.json"

func init() {
	loadColorPalettesFromDisk()
}

// Todo: add logic for adding, removing and the html page and handler

func loadColorPalettesFromDisk() {
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(folderPath, 0750)
		if err != nil {
			log.Fatalf("Directioary can not be created: %v", err)
		}
		ColorPaletteMap = make(map[string]ColorPalette)
		return
	}
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Color file not found: %v", err)
	}
	err = json.NewDecoder(file).Decode(&ColorPaletteMap)

}

func saveColorPalettesToDisk() {
	file, err := os.Create(filePath)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	err = file.Truncate(0)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	_, err = file.Seek(0, 0)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	err = json.NewEncoder(file).Encode(&ColorPaletteMap)
	if err != nil {
		slog.Error(err.Error())
	}
}
