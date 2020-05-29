package callbacks

import (
	"fmt"
	"reflect"

	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/clause"
	"github.com/jinzhu/gorm/schema"
)

func BeforeDelete(db *gorm.DB) {
	if db.Statement.Schema != nil && db.Statement.Schema.BeforeDelete {
		callMethod := func(value interface{}) bool {
			if db.Statement.Schema.BeforeDelete {
				if i, ok := value.(gorm.BeforeDeleteInterface); ok {
					i.BeforeDelete(db)
					return true
				}
			}
			return false
		}

		if ok := callMethod(db.Statement.Dest); !ok {
			switch db.Statement.ReflectValue.Kind() {
			case reflect.Slice, reflect.Array:
				for i := 0; i <= db.Statement.ReflectValue.Len(); i++ {
					callMethod(db.Statement.ReflectValue.Index(i).Interface())
				}
			case reflect.Struct:
				callMethod(db.Statement.ReflectValue.Interface())
			}
		}
	}
}

func Delete(db *gorm.DB) {
	if db.Statement.Schema != nil && !db.Statement.Unscoped {
		for _, c := range db.Statement.Schema.DeleteClauses {
			db.Statement.AddClause(c)
			fmt.Println(db.Statement.SQL.String())
		}
	}

	if db.Statement.SQL.String() == "" {
		db.Statement.AddClauseIfNotExists(clause.Delete{})

		if db.Statement.Schema != nil {
			_, queryValues := schema.GetIdentityFieldValuesMap(db.Statement.ReflectValue, db.Statement.Schema.PrimaryFields)
			column, values := schema.ToQueryValues(db.Statement.Schema.PrimaryFieldDBNames, queryValues)

			if len(values) > 0 {
				db.Statement.AddClause(clause.Where{Exprs: []clause.Expression{clause.IN{Column: column, Values: values}}})
			}

			if db.Statement.Dest != db.Statement.Model && db.Statement.Model != nil {
				_, queryValues = schema.GetIdentityFieldValuesMap(reflect.ValueOf(db.Statement.Model), db.Statement.Schema.PrimaryFields)
				column, values = schema.ToQueryValues(db.Statement.Schema.PrimaryFieldDBNames, queryValues)

				if len(values) > 0 {
					db.Statement.AddClause(clause.Where{Exprs: []clause.Expression{clause.IN{Column: column, Values: values}}})
				}
			}
		}

		if _, ok := db.Statement.Clauses["WHERE"]; !ok {
			db.AddError(gorm.ErrMissingWhereClause)
			return
		}

		db.Statement.AddClauseIfNotExists(clause.From{})
		db.Statement.Build("DELETE", "FROM", "WHERE")
	}

	result, err := db.Statement.ConnPool.ExecContext(db.Statement.Context, db.Statement.SQL.String(), db.Statement.Vars...)

	if err == nil {
		db.RowsAffected, _ = result.RowsAffected()
	} else {
		db.AddError(err)
	}
}

func AfterDelete(db *gorm.DB) {
	if db.Statement.Schema != nil && db.Statement.Schema.AfterDelete {
		callMethod := func(value interface{}) bool {
			if db.Statement.Schema.AfterDelete {
				if i, ok := value.(gorm.AfterDeleteInterface); ok {
					i.AfterDelete(db)
					return true
				}
			}
			return false
		}

		if ok := callMethod(db.Statement.Dest); !ok {
			switch db.Statement.ReflectValue.Kind() {
			case reflect.Slice, reflect.Array:
				for i := 0; i <= db.Statement.ReflectValue.Len(); i++ {
					callMethod(db.Statement.ReflectValue.Index(i).Interface())
				}
			case reflect.Struct:
				callMethod(db.Statement.ReflectValue.Interface())
			}
		}
	}
}