package builder

import (
	"PoliSim/data/database"
	"encoding/json"
	"fmt"
	"os"
)

var Translation = make(map[string]string)
var Configuration = make(map[string]string)
var RawStartPageContent = ""

func ImportTranslation(lang string) {
	file, err := os.ReadFile("resources/" + lang + ".json")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "Error while trying to read the language json file:\n"+err.Error()+"\n")
		os.Exit(1)
	}

	err = json.Unmarshal(file, &Translation)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "Error while trying to parse the translation json:\n"+err.Error()+"\n")
		os.Exit(1)
	}
	_, _ = fmt.Fprintf(os.Stdout, "Imported language localisation\n")
	importConfiguration()
	importRawHTMLForStartPage()
	setTranslationsForDatabase()
}

func setTranslationsForDatabase() {
	database.VoteTranslation[database.SingleVote] = Translation["translationSingleVote"]
	database.VoteTranslation[database.MultipleVotes] = Translation["translationMultipleVotes"]
	database.VoteTranslation[database.RankedVotes] = Translation["translationRankedVotes"]
	database.VoteTranslation[database.ThreeCategoryVoting] = Translation["translationThreeCategoryVoting"]

	database.StatusTranslation[database.Public] = Translation["translationPublic"]
	database.StatusTranslation[database.Private] = Translation["translationPrivate"]
	database.StatusTranslation[database.Secret] = Translation["translationSecret"]
	database.StatusTranslation[database.Hidden] = Translation["translationHidden"]

	database.RoleTranslation[database.PressAccount] = Translation["translationPressAccount"]
	database.RoleTranslation[database.User] = Translation["translationUser"]
	database.RoleTranslation[database.MediaAdmin] = Translation["translationMediaAdmin"]
	database.RoleTranslation[database.Admin] = Translation["translationAdmin"]
	database.RoleTranslation[database.HeadAdmin] = Translation["translationHeadAdmin"]
}

func importConfiguration() {
	file, err := os.ReadFile("resources/" + os.Getenv("CONFIG") + ".json")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "Error while trying to read the config json file:\n"+err.Error()+"\n")
		os.Exit(1)
	}

	err = json.Unmarshal(file, &Configuration)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "Error while trying to parse the config json:\n"+err.Error()+"\n")
		os.Exit(1)
	}
	_, _ = fmt.Fprintf(os.Stdout, "Imported config\n")
}

func importRawHTMLForStartPage() {
	file, err := os.ReadFile("resources/startPage.html")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stdout, "Error while trying to read the the start page HTML file:\n"+err.Error()+"\n")
		os.Exit(1)
	}
	RawStartPageContent = string(file)
}
