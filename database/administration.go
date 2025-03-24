package database

import "github.com/lib/pq"

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
