package postgres

import (
	"fmt"
	"strings"

	"github.com/aodin/aspect"
)

// TODO Replace with a more robust binary clause?
type WithClause struct {
	Name     string
	Operator Operator
}

func (clause WithClause) String() string {
	output, _ := clause.Compile(&PostGres{}, aspect.Params())
	return output
}

func (clause WithClause) Compile(d aspect.Dialect, params *aspect.Parameters) (string, error) {
	return fmt.Sprintf(`"%s" WITH %s`, clause.Name, clause.Operator), nil
}

type ExcludeConstraint struct {
	method  IndexMethod
	clauses []WithClause
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
	// There must be at least one clause
	if len(exclude.clauses) == 0 {
		return fmt.Errorf(
			"postgres: an EXCLUDE constraint must contain at least one clause",
		)
	}

	// All excluded elements must exist in the given table
	for _, clause := range exclude.clauses {
		_, exists := table.C[clause.Name]
		if !exists {
			return fmt.Errorf("No column with the name '%s' exists in the table '%s'. Is it declared after the PrimaryKey declaration?", clause.Name, table.Name)
		}
	}

	table.AddCreatable(exclude)
	return nil
}

// Using sets the ExcludeConstraints index method
func (exclude ExcludeConstraint) Using(method IndexMethod) ExcludeConstraint {
	exclude.method = method
	return exclude
}

// Exclude creates a new ExcludeConstraint
func Exclude(clauses ...WithClause) ExcludeConstraint {
	return ExcludeConstraint{clauses: clauses}
}
