package scanner

import (
	"fmt"
	"github.com/billyninja/pgtools/connector"
	"log"
)


type TableName 		string
type ColumnName 	string
type ConstraintType string

type Column struct {
	Name          ColumnName  	`db:"column_name"`
	CharMaxLength *uint16 		`db:"character_maximum_length"`
	Type          string  		`db:"data_type"`
	Default       *string 		`db:"column_default"`
	Nullable      string  		`db:"is_nullable"`
}

type Constraint struct {
	Name          ColumnName  		`db:"name"`
	Type 		  ConstraintType	`db:"ctype"`
	FTable 	  	  *TableName		`db:"ftable"`
	Columns   	  []uint8 			`db:"col"`
}

type Table struct {
	Name      TableName 	`db:"table_name"`
	TableType string 		`db:"table_type"`
	Columns   		[]*Column
	Constraints   	[]*Constraint
}

func (t Table) String() string {
	out := fmt.Sprintf("Table: %s\n\tColumns:\n", t.Name)
	for _, c := range t.Columns {
		dft := "[no default]"
		if c.Default != nil {
			dft = *c.Default
		}
		out += fmt.Sprintf("\t\t%s (%s) dft: %s  Null? %s\n", c.Name, c.Type, dft, c.Nullable)
	}
	return out
}

func GetAllTables(conn *connector.Connector) []*Table {
	qTables := "SELECT table_name, table_type FROM information_schema.tables WHERE table_schema='public';"
	rows, err := conn.Sel(qTables)
	if err != nil {
		log.Panic("err:", err)
	}

	var allTables []*Table

	for rows.Next() {
		tb := &Table{}
		err := rows.StructScan(tb)
		if err != nil {
			log.Panic("---\n%v\n----", err)
		}
		allTables = append(allTables, tb)
	}

	for _, tb := range allTables {
		qColumns := fmt.Sprintf(`
            SELECT column_name, data_type, column_default, character_maximum_length, is_nullable FROM information_schema.columns
            WHERE table_schema = 'public' AND table_name = '%s'
        `, tb.Name)

		rows, err := conn.Sel(qColumns)
		if err != nil {
			log.Panic("err parsing table struct:\n %v", err)
		}

		for rows.Next() {
			cl := &Column{}
			err := rows.StructScan(cl)
			if err != nil {
				log.Panic("err parsing table struct:\n %v", err)
			}
			tb.Columns = append(tb.Columns, cl)
		}
		mapConstraints(conn, tb)
	}

	return allTables
}


func mapConstraints(conn *connector.Connector, tb *Table) {
	q_constraints := `
	SELECT
	  c.conname as name,
	  (select constraint_type from information_schema.table_constraints where constraint_name = c.conname) as ctype,
	  (select array_agg(attname) from pg_attribute where attrelid = c.conrelid and ARRAY[attnum] <@ c.conkey) as col,
	  (select r.relname from pg_class r where r.oid = c.confrelid) as ftable
	FROM pg_constraint c
	WHERE
	    c.conrelid = (
	        select oid from pg_class where relname = '%s'
	    );
	`
	q_constraints = fmt.Sprintf(q_constraints, tb.Name)
	println(q_constraints)

	rows, err := conn.Sel(q_constraints)
	if err != nil {
		log.Panic("err parsing table struct:\n %v", err)
	}

	for rows.Next() {
		ct := &Constraint{}
		err := rows.StructScan(ct)
		if err != nil {
			log.Panic("err parsing table constraints:\n %v", err)
		}
		tb.Constraints = append(tb.Constraints, ct)
	}
}
