package database

import (
	"log"
	"os"
	"strings"
)

func migrate() {
	if strings.ToUpper(os.Getenv("MIGRATE")) != "YES" {
		return
	}
	log.Println("Trying to migrate old data")
	updateVotes()
	log.Println("Finished migrating old data")
}

func updateVotes() {
	tx, err := openTransaction()
	defer tx.Close()
	if err != nil {
		log.Fatalf(err.Error())
	}

	result, err := tx.Run(`MATCH (r:Result)
WHERE r.csv IS NULL RETURN r, elementId(r);`, nil)
	if err != nil {
		log.Fatalf("Failed to execute update Votes")
	} else if result.Peek() {
		for result.Next() {
			props := GetPropsMapForRecordPosition(result.Record(), 0)
			vote := AccountVotes{
				Question:        props.GetString("question"),
				IterableAnswers: props.GetArray("answers"),
				Anonymous:       props.GetBool("anonymous"),
				Type:            VoteType(props.GetInt("type")),
				AnswerAmount:    props.GetInt("amount"),
				IllegalVotes:    nil,
				Illegal:         props.GetArray("illegal"),
				List:            nil,
				Voter:           props.GetArray("voter"),
				Votes:           props.GetArray("votes"),
			}
			generateCSV(&vote)
			err = tx.RunWithoutResult(`MATCH (r:Result) WHERE elementId(r) = $id SET r.csv = $csv;`,
				map[string]any{
					"id":  result.Record().Values[1],
					"csv": vote.CSV,
				})
			if err != nil {
				log.Fatalf("Failed to execute update Votes")
			}
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Fatalf(err.Error())
	}
}
