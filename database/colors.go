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
}

var ColorPaletteMap = map[string]ColorPalette{}

func init() {
	_, exists := ColorPaletteMap[loc.StandardColorName]
	if !exists {
		ColorPaletteMap[loc.StandardColorName] = ColorPalette{
			Name:       loc.StandardColorName,
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
	value := ""
	err := postgresDB.QueryRow(`SELECT o.account_name FROM organisation_to_account 
    INNER JOIN ownership o on organisation_to_account.account_name = o.account_name
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
	ColorPaletteMap[color.Name] = *color
	return nil
}

func RemoveColorPalette(name string, acc *Account) (*ColorPalette, error) {
	if !HasPrivilegesForColorsDelete(acc) {
		return nil, NotAllowedError
	}
	result := ColorPaletteMap[name]
	delete(ColorPaletteMap, name)
	return &result, nil
}

func loadColorPalettesFromDB() {
	//Todo: move this into the migration function
	_, err := postgresDB.Exec(`CREATE TABLE colors (
    name TEXT PRIMARY KEY,
    background TEXT,
    text TEXT,
    link TEXT
)`)
	if err != nil {
		log.Fatalf("Could not create postgres color tabel: %v", err)
	}

	results, err := postgresDB.Query("SELECT name, background, text, link FROM colors;")
	if err != nil {
		log.Fatalf("Could not read postgres color tabel: %v", err)
	}

	for results.Next() {
		var color ColorPalette

		err = results.Scan(&color.Name, &color.Background, &color.Text, &color.Link)
		if err != nil {
			slog.Error("could not scan entry correctly:", "err", err)
		}

		ColorPaletteMap[color.Name] = color
	}
}

func saveColorPalettesToDB() {
	queryStmt := `
        INSERT INTO colors (name, background, text, link)
        VALUES ($1, $2, $3, $4) 
        ON CONFLICT (name) DO UPDATE SET background=$2, text = $3, link = $4;
    `
	for name := range ColorPaletteMap {
		color := ColorPaletteMap[name]
		_, err := postgresDB.Exec(queryStmt, &color.Name, &color.Background, &color.Text, &color.Link)
		if err != nil {
			slog.Error("While saving colors encountered an error: ", "err", err)
		}
	}
}
