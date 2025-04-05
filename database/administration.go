package database

import (
	"github.com/lib/pq"
	"regexp"
	"strconv"
)

type AdministrationQuery struct {
	Rows  [][]string
	Error error
}

func (a *AdministrationQuery) HasRows() bool {
	return len(a.Rows) != 0
}

func (a *AdministrationQuery) HasError() bool {
	return a.Error != nil
}

func ExecuteQueryString(query string) *AdministrationQuery {
	// For this to work the returned value must be a single "SELECT array[x,y::TEXT,z] ..." whose internal values are cast to TEXT.
	admin := &AdministrationQuery{}
	result, err := postgresDB.Query(query)
	if err != nil {
		admin.Error = err
		return admin
	}
	defer closeRows(result)
	retVal := make([][]string, 0)
	row := make([]string, 0)
	for result.Next() {
		err = result.Scan(pq.Array(&row))
		if err != nil {
			admin.Error = err
			return admin
		}
		retVal = append(retVal, row)
	}
	admin.Error = nil
	admin.Rows = retVal
	return admin
}

func ExecuteNamedQuery(query string, parameters []any) *AdministrationQuery {
	var re = regexp.MustCompile(`(?m)((?:.|\s)+?;)\s*?(--PARAM=)([0-9]+):([0-9]+)`)
	matches := re.FindAllStringSubmatch(query, -1)
	admin := &AdministrationQuery{}

	if len(matches) == 0 {
		admin.Error = DbError("no commands found to execute")
		return admin
	}

	tx, err := postgresDB.Begin()

	if err != nil {
		admin.Error = err
		return admin
	}

	var start int
	var end int
	for matchPos := range len(matches) {
		start, _ = strconv.Atoi(matches[matchPos][3])
		end, _ = strconv.Atoi(matches[matchPos][4])
		if matchPos == len(matches)-1 {
			break
		}

		_, err = tx.Exec(matches[matchPos][1], parameters[start:end]...)
		if err != nil {
			admin.Error = DbError("error in Command " + strconv.Itoa(matchPos+1) + ": " + err.Error())
			_ = tx.Rollback()
			return admin
		}
	}

	result, err := tx.Query(matches[len(matches)-1][1], parameters[start:end]...)
	if err != nil {
		_ = tx.Rollback()
		admin.Error = DbError("error in Query Command: " + err.Error())
		return admin
	}

	defer closeRows(result)
	retVal := make([][]string, 0)
	row := make([]string, 0)
	for result.Next() {
		err = result.Scan(pq.Array(&row))
		if err != nil {
			_ = tx.Rollback()
			admin.Error = err
			return admin
		}
		retVal = append(retVal, row)
	}

	err = tx.Commit()
	if err != nil {
		_ = tx.Rollback()
		admin.Error = DbError("error in Commit: " + err.Error())
		return admin
	}

	admin.Error = nil
	admin.Rows = retVal
	return admin
}
