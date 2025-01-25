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
	_, exists := ColorPaletteMap[""]
	if !exists {
		ColorPaletteMap[""] = ColorPalette{
			Name:       "Standard",
			Background: "#000000",
			Text:       "#FFFFFF",
			Link:       "#9999FF",
		}
	}
}

func HasPrivilegesForColorsAdd(acc *Account) bool {
	if acc.IsAtLeastAdmin() {
		return true
	}
	result, err := makeRequest(`MATCH (a:Account)-[:ADMIN|OWNER*..]->(o:Organisation) WHERE a.name = $name 
RETURN o.name;`, map[string]any{"name": acc.GetName()})
	if err != nil {
		return false
	} else if len(result.Records) == 0 {
		return false
	}
	return true
}

func HasPrivilegesForColorsDelete(acc *Account) bool {
	return acc.IsAtLeastAdmin()
}

func AddColorPalette(color *ColorPalette, acc *Account) error {
	if !HasPrivilegesForColorsAdd(acc) {
		return notAllowedError
	}
	ColorPaletteMap[color.Name] = *color
	return nil
}

func RemoveColorPalette(name string, acc *Account) (*ColorPalette, error) {
	if !HasPrivilegesForColorsDelete(acc) {
		return nil, notAllowedError
	}
	result := ColorPaletteMap[name]
	delete(ColorPaletteMap, name)
	return &result, nil
}

func loadColorPalettesFromDisk() {
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(folderPath, os.ModePerm)
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
