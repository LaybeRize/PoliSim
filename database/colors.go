package database

import (
	loc "PoliSim/localisation"
	"log"
	"log/slog"
)

type ColorPalette struct {
	Name       string
	Background string
	Text       string
	Link       string
	Permanent  bool
}

var ColorPaletteMap = map[string]ColorPalette{}

func HasPrivilegesForColorsAdd(acc *Account) bool {
	if acc.IsAtLeastAdmin() {
		return true
	}
	value := ""
	err := postgresDB.QueryRow(`SELECT account_name FROM organisation_linked 
	WHERE is_admin = true AND owner_name = $1
                    LIMIT 1;`, acc.GetName()).Scan(&value)
	if err != nil {
		return false
	}
	return true
}

func HasPrivilegesForColorsDelete(acc *Account) bool {
	return acc.IsAtLeastAdmin()
}

func AddColorPalette(color *ColorPalette, acc *Account) error {
	if !HasPrivilegesForColorsAdd(acc) {
		return NotAllowedError
	}
	if _, exists := ColorPaletteMap[color.Name]; exists {
		color.Permanent = ColorPaletteMap[color.Name].Permanent
	} else {
		color.Permanent = false
	}
	ColorPaletteMap[color.Name] = *color
	return nil
}

func RemoveColorPalette(name string, acc *Account) (*ColorPalette, error) {
	if !HasPrivilegesForColorsDelete(acc) {
		return nil, NotAllowedError
	}
	result := ColorPaletteMap[name]
	if result.Permanent {
		return nil, CanNotDeleteColor
	}
	delete(ColorPaletteMap, name)
	return &result, nil
}

func loadColorPalettesFromDB() {
	results, err := postgresDB.Query("SELECT name, background, text, link, permanent FROM colors;")
	if err != nil {
		log.Fatalf("Could not read postgres color tabel: %v", err)
	}
	hasPermanentColor := false
	for results.Next() {
		var color ColorPalette

		err = results.Scan(&color.Name, &color.Background, &color.Text, &color.Link, &color.Permanent)
		if err != nil {
			slog.Error("could not scan entry correctly:", "err", err)
		}
		hasPermanentColor = hasPermanentColor || color.Permanent
		ColorPaletteMap[color.Name] = color
	}

	if !hasPermanentColor {
		_, err = postgresDB.Exec(`INSERT INTO colors (name, background, text, link, permanent) 
VALUES ($1, '#000000', '#FFFFFF', '#9999FF', true)`, loc.StandardColorName)
		if err != nil {
			log.Fatalf("could not create color entry correctly: %v", err)
		}
		log.Println("Created new standard color")
		ColorPaletteMap[loc.StandardColorName] = ColorPalette{
			Name:       loc.StandardColorName,
			Background: "#000000",
			Text:       "#FFFFFF",
			Link:       "#9999FF",
			Permanent:  true,
		}
	}
}

func saveColorPalettesToDB() {
	queryStmt := `
        INSERT INTO colors (name, background, text, link, permanent)
        VALUES ($1, $2, $3, $4, $5) 
        ON CONFLICT (name) DO UPDATE SET background=$2, text = $3, link = $4;
    `
	for name := range ColorPaletteMap {
		color := ColorPaletteMap[name]
		_, err := postgresDB.Exec(queryStmt, &color.Name, &color.Background, &color.Text, &color.Link, &color.Permanent)
		if err != nil {
			slog.Error("While saving colors encountered an error: ", "err", err)
		}
	}
}
