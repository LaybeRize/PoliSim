package validation

type TitleModification struct {
	Name      string   `input:"name"`
	NewName   string   `input:"newName"`
	MainGroup string   `input:"mainGroup"`
	SubGroup  string   `input:"subGroup"`
	Holder    []string `input:"holder"`
}
