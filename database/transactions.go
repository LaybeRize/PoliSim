package database

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log/slog"
	"time"
)

type dbTransaction struct {
	tx neo4j.ExplicitTransaction
}

func (d *dbTransaction) Close() {
	_ = d.tx.Close(ctx)
}

func (d *dbTransaction) Run(query string, parameter map[string]any) (*dbResult, error) {
	var err error
	result := dbResult{}
	result.res, err = d.tx.Run(ctx, query, parameter)
	if err != nil {
		slog.Error(err.Error())
		_ = d.tx.Rollback(ctx)
	}
	return &result, err
}

func (d *dbTransaction) RunWithoutResult(query string, parameter map[string]any) error {
	_, err := d.tx.Run(ctx, query, parameter)
	if err != nil {
		slog.Debug(err.Error())
		_ = d.tx.Rollback(ctx)
	}
	return err
}

func (d *dbTransaction) Commit() error {
	err := d.tx.Commit(ctx)
	return err
}

type dbResult struct {
	res neo4j.ResultWithContext
}

func (d *dbResult) Peek() bool {
	return d != nil && d.res.Peek(ctx)
}

func (d *dbResult) Next() bool {
	return d != nil && d.res.Next(ctx)
}

func (d *dbResult) Record() *neo4j.Record {
	return d.res.Record()
}

func openTransaction() (*dbTransaction, error) {
	var err error
	result := dbTransaction{}
	result.tx, err = driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: ""}).BeginTransaction(ctx)
	return &result, err
}

func makeRequest(query string, parameter map[string]any) ([]*neo4j.Record, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver, query, parameter,
		neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		return nil, err
	}
	return result.Records, err
}

type PropsMap map[string]any

func GetPropsMapForRecordPosition(record *neo4j.Record, position int) PropsMap {
	if position >= len(record.Values) {
		return nil
	}
	v := record.Values[position]
	switch v.(type) {
	case neo4j.Node:
		return v.(neo4j.Node).Props
	case neo4j.Relationship:
		return v.(neo4j.Relationship).Props
	default:
		return nil
	}
}

func (p PropsMap) GetString(field string) string {
	v, ok := p[field].(string)
	if !ok {
		slog.Error("Failed to pass Field to String", "map", p, "field name", field)
	}
	return v
}

func (p PropsMap) GetInt(field string) int {
	v, ok := p[field].(int64)
	if !ok {
		slog.Error("Failed to pass Field to Int", "map", p, "field name", field)
	}
	return int(v)
}

func (p PropsMap) GetBool(field string) bool {
	v, ok := p[field].(bool)
	if !ok {
		slog.Error("Failed to pass Field to Int", "map", p, "field name", field)
	}
	return v
}

func (p PropsMap) GetTime(field string) time.Time {
	v, ok := p[field].(time.Time)
	if !ok {
		slog.Error("Failed to pass Field to Time", "map", p, "field name", field)
	}
	return v
}

func (p PropsMap) GetArray(field string) []any {
	v, ok := p[field].([]any)
	if !ok {
		slog.Error("Failed to pass Field to Array", "map", p, "field name", field)
	}
	return v
}
