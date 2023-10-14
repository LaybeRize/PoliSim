package validation

type TitleModification struct {
	Name      string   `input:"name"`
	NewName   string   `input:"newName"`
	MainGroup string   `input:"mainGroup"`
	SubGroup  string   `input:"subGroup"`
	Flair     string   `input:"flair"`
	Holder    []string `input:"holder"`
}

func (form *TitleModification) CreateTitle() Message {
	return Message{}
}

func (form *TitleModification) SearchTitle() Message {
	return Message{}
}

func (form *TitleModification) ModifyTitle() Message {
	return Message{}
}

func (form *TitleModification) DeleteTitle() Message {
	return Message{}
}
