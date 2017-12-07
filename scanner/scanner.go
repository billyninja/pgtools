package scanner

import (
	"fmt"
	"github.com/billyninja/pgtools/connector"
	"log"
	"strings"
)

type TableName string
type ColumnName string
type ConstraintType string

type Column struct {
	Name          ColumnName `db:"column_name"`
	CharMaxLength *uint16    `db:"character_maximum_length"`
	Type          string     `db:"data_type"`
	Default       *string    `db:"column_default"`
	Nullable      string     `db:"is_nullable"`
	RefersTo      *Table
}

func field_type_translation(column_type string) (string, string) {

	switch column_type {
	case "character varying":
		return "input", "text"
	case "timestamp without time zone":
		return "input", "date"
	case "timestamp with time zone":
		return "input", "date"
	case "numeric":
		return "input", "number"
	case "text":
		return "textarea", "textarea"
	case "boolean":
		return "input", "checkbox"
	default:
		log.Printf("\n\nUnmapped PSQL field type %s\n\n", column_type)
		return "input", "text"
	}

	return "err", "err"
}

func (cl Column) Input(value interface{}) string {
	archtype, fieldtype := field_type_translation(cl.Type)

	required := ""
	if cl.Nullable == "NO" {
		required = `required="required"`
	}

	max_length := ""
	if cl.CharMaxLength != nil {
		max_length = fmt.Sprintf(`max_length="%d"`, *cl.CharMaxLength)
	}

	value_str := ""
	if value != nil {

		switch v := value.(type) {
		case bool:
			if v == true {
				value_str = `checked="checked"`
			} else {
				value_str = `checked=""`
			}
			break
		case string:
			if archtype == "textarea" {
				value_str = v
			} else {
				value_str = `value="` + v + `"`
			}

		default:
			println("default")
		}
	}

	input := ""
	if archtype == "textarea" {
		input = fmt.Sprintf(`<textarea name="%s" %s>%s</textarea>`, cl.Name, required, value_str)
	} else {
		input = fmt.Sprintf(`<input type="%s" name="%s" %s %s %s/>`,
			fieldtype, cl.Name, max_length, required, value_str)
	}

	return input
}

type Constraint struct {
	Name    ColumnName     `db:"name"`
	Type    ConstraintType `db:"ctype"`
	FTable  *TableName     `db:"ftable"`
	Columns string         `db:"col"`
}

type Table struct {
	Name        TableName `db:"table_name"`
	TableType   string    `db:"table_type"`
	PkExp       string
	Columns     []*Column
	Constraints []*Constraint
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
		mapConstraints(conn, tb, allTables)
	}

	return allTables
}

func mapConstraints(conn *connector.Connector, tb *Table, allTables []*Table) {
	q_constraints := `
	SELECT
	  c.conname::varchar as name,
	  (select constraint_type from information_schema.table_constraints where constraint_name = c.conname) as ctype,
	  (select array_agg(attname::varchar)::varchar from pg_attribute where attrelid = c.conrelid and ARRAY[attnum] <@ c.conkey) as col,
	  (select r.relname from pg_class r where r.oid = c.confrelid) as ftable
	FROM pg_constraint c
	WHERE
	    c.conrelid = (select oid from pg_class where relname = '%s');
	`
	q_constraints = fmt.Sprintf(q_constraints, tb.Name)
	rows, err := conn.Sel(q_constraints)
	if err != nil {
		log.Panic("err parsing table struct:\n ", err)
	}

	for rows.Next() {
		ct := &Constraint{}
		err := rows.StructScan(ct)
		if err != nil {
			log.Panic("err parsing table constraints:\n ", err)
		}

		ct.Columns = ct.Columns[1 : len(ct.Columns)-1]

		if ct.Type == "PRIMARY KEY" {
			tb.PkExp = ""
			// Suporting composite keys
			for _, pkCl := range strings.Split(ct.Columns, ",") {
				tb.PkExp += fmt.Sprintf(`"%s", `, pkCl)
			}
			tb.PkExp = tb.PkExp[0 : len(tb.PkExp)-2]
		}

		if ct.Type == "FOREIGN KEY" && ct.FTable != nil {
			ctCol := getColumnByName(tb, ColumnName(ct.Columns))
			refTable := getTableByName(allTables, *ct.FTable)

			if ctCol == nil || refTable == nil {
				log.Panicf(
					"couldn't infere relationship between %s.%s and %s(table) | local column %s foreign table %s",
					tb.Name, ct.Name, *ct.FTable, ctCol, refTable)
			}
			ctCol.RefersTo = refTable
		}

		tb.Constraints = append(tb.Constraints, ct)
	}
}

func (tb *Table) ReverseRefQuery(rnd bool) string {
	q1 := fmt.Sprintf(`SELECT %s FROM %s `, tb.PkExp, tb.Name)
	if rnd {
		q1 += `ORDER BY random() LIMIT 1`
	}

	return q1
}

func getTableByName(allTables []*Table, name TableName) *Table {
	for _, at := range allTables {
		if at.Name == name {
			return at
		}
	}
	return nil
}

func getColumnByName(tb *Table, name ColumnName) *Column {
	for _, cl := range tb.Columns {
		if cl.Name == name {
			return cl
		}
	}
	return nil
}
