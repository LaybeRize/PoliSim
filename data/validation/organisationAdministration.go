package validation

type OrganisationModification struct {
	Name      string   `input:"name"`
	MainGroup string   `input:"mainGroup"`
	SubGroup  string   `input:"subGroup"`
	Status    string   `input:"status"`
	Flair     string   `input:"flair"`
	User      []string `input:"user"`
	Admins    []string `input:"admins"`
}
