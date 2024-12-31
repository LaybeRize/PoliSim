package database

type Title struct {
	Name     string
	MainType string
	SubType  string
	Flair    string
}

func (t *Title) Exists() bool {
	return t != nil
}

func CreateTitle(title *Title, holderNames []string) error {
	tx, err := openTransaction()
	defer tx.Close(ctx)
	if err != nil {
		return err
	}
	_, err = tx.Run(ctx,
		`CREATE (:Title {name: $name , main_type: $maintype , 
sub_type: $subtype , flair: $flair});`, map[string]any{
			"name":     title.Name,
			"maintype": title.MainType,
			"subtype":  title.SubType,
			"flair":    title.Flair})
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	_, err = tx.Run(ctx, `MATCH (a:Account), (t:Title) WHERE a.name IN $names  
AND t.name = $title CREATE (a)-[:HAS]->(t);`, map[string]any{
		"title": title.Name,
		"names": holderNames})
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	err = tx.Commit(ctx)
	return err
}

func UpdateTitle(oldtitle string, title *Title, holderNames []string) error {
	tx, err := openTransaction()
	defer tx.Close(ctx)
	if err != nil {
		return err
	}
	_, err = tx.Run(ctx,
		`MATCH (t:Title) WHERE name = $oldName 
SET name = $name , main_type = $maintype , 
sub_type = $subtype , flair = $flair;`,
		map[string]any{
			"oldName":  oldtitle,
			"name":     title.Name,
			"maintype": title.MainType,
			"subtype":  title.SubType,
			"flair":    title.Flair})
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	_, err = tx.Run(ctx, `MATCH (:Account)-[r]->(t:Title) WHERE t.name = $title 
DELETE r;`, map[string]any{"title": title.Name})
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	_, err = tx.Run(ctx, `MATCH (a:Account), (t:Title) WHERE a.name IN $names  
AND t.name = $title CREATE (a)-[:HAS]->(t);`, map[string]any{
		"title": title.Name,
		"names": holderNames})
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	err = tx.Commit(ctx)
	return err
}
