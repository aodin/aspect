package postgres

import (
	"fmt"
	"strings"

	"github.com/aodin/aspect"
)

type ExcludeConstraint struct {
	method  IndexMethod
	clauses []ExcludeClause
}

var _ aspect.Creatable = ExcludeConstraint{}
var _ aspect.TableModifier = ExcludeConstraint{}

// Create returns the proper syntax for CREATE TABLE commands.
func (exclude ExcludeConstraint) Create(d aspect.Dialect) (string, error) {
	// There must be at least one clause
	if len(exclude.clauses) == 0 {
		return "", fmt.Errorf(
			"postgres: an EXCLUDE constraint must contain at least one clause",
		)
	}

	compiled := "EXCLUDE"
	if exclude.method != "" {
		compiled += fmt.Sprintf(" USING %s", exclude.method)
	}
	clauses := make([]string, len(exclude.clauses))
	var err error
	for i, clause := range exclude.clauses {
		if clauses[i], err = clause.Compile(d, &aspect.Parameters{}); err != nil {
			return "", err
		}
	}

	compiled += fmt.Sprintf(" (%s)", strings.Join(clauses, ", "))
	return compiled, nil
}

// Modify adds the ExcludeConstraint to the table's creates
func (exclude ExcludeConstraint) Modify(table *aspect.TableElem) error {
	table.AddCreatable(exclude)
	return nil
}

// Using sets the ExcludeConstraints index method
func (exclude ExcludeConstraint) Using(method IndexMethod) ExcludeConstraint {
	exclude.method = method
	return exclude
}

// Exclude creates a new ExcludeConstraint
func Exclude(clauses ...ExcludeClause) ExcludeConstraint {
	return ExcludeConstraint{clauses: clauses}
}
