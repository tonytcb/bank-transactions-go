package repository

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-sql-driver/mysql"
)

// ErrDuplicateEntry represents a duplicate entry error
type ErrDuplicateEntry struct {
	field string
	value string
}

// NewErrDuplicatedEntry creates an ErrDuplicateEntry struct
func NewErrDuplicatedEntry(field string, value string) *ErrDuplicateEntry {
	return &ErrDuplicateEntry{field: field, value: value}
}

// Error returns the formatted error message
func (e ErrDuplicateEntry) Error() string {
	return fmt.Sprintf(`duplicate entry '%s' for field '%s'`, e.Value(), e.Field())
}

// Field returns the duplicate value
func (e ErrDuplicateEntry) Value() string {
	return e.value
}

// Field returns the field value
func (e ErrDuplicateEntry) Field() string {
	return e.field
}

// --

// ErrRegisterNotFound represents an error when the register was not found
type ErrRegisterNotFound struct {
	field string
	value string
}

// NewErrRegisterNotFound creates an ErrDuplicateEntry struct
func NewErrRegisterNotFound(field, value string) *ErrRegisterNotFound {
	return &ErrRegisterNotFound{field: field, value: value}
}

// Field returns the field not found
func (e ErrRegisterNotFound) Field() string {
	return e.field
}

// Field returns the duplicate value
func (e ErrRegisterNotFound) Value() string {
	return e.value
}

// Error returns the formatted error message
func (e ErrRegisterNotFound) Error() string {
	return fmt.Sprintf(`'%s' '%s' not found`, e.Field(), e.Value())
}

// --

// ErrRegisterNotFound represents an error when occurred an error to load data
type ErrLoadInvalidData struct {
	field string
}

// NewErrLoadInvalidData creates an ErrLoadInvalidData struct
func NewErrLoadInvalidData(field string) *ErrLoadInvalidData {
	return &ErrLoadInvalidData{field: field}
}

// Error returns the formatted error message
func (e ErrLoadInvalidData) Error() string {
	return fmt.Sprintf(`error to load values from '%s'`, e.field)
}

// --

type ErrForeignKeyConstraint struct {
	table      string
	constraint string
	foreignKey string
	references string
}

func (e ErrForeignKeyConstraint) ForeignKey() string {
	return e.foreignKey
}

func NewErrForeignKeyConstraint(table string, constraint string, foreignKey string, references string) *ErrForeignKeyConstraint {
	return &ErrForeignKeyConstraint{table: table, constraint: constraint, foreignKey: foreignKey, references: references}
}

// Error returns the formatted error message
func (e ErrForeignKeyConstraint) Error() string {
	return fmt.Sprintf(`'%s' not found`, e.ForeignKey())
}

// --

func translateMySqlErrors(err *mysql.MySQLError) error {
	const (
		duplicateEntryErrorCode       = 1062
		foreignKeyConstraintErrorCode = 1452
	)

	if err.Number == duplicateEntryErrorCode {
		duplicatedErrorRegex := regexp.MustCompile(`Duplicate entry '(\w+)' for key '(\w+)'`)

		if match := duplicatedErrorRegex.FindAllStringSubmatch(err.Message, -1); len(match) > 0 {
			return NewErrDuplicatedEntry(match[0][2], match[0][1])
		}
	}

	if err.Number == foreignKeyConstraintErrorCode {
		foreignKeyErrorRegex := regexp.MustCompile(`Cannot add or update a child row: a foreign key constraint fails \('([a-z0-9-_]+)'.'([a-z0-9-_]+)', CONSTRAINT '([a-z0-9-_]+)' FOREIGN KEY \('([a-z0-9-_]+)'\) REFERENCES '([a-z0-9-_]+)' \('([a-z0-9-_]+)'\)\)`)

		message := strings.ReplaceAll(err.Message, "`", "'")

		if match := foreignKeyErrorRegex.FindAllStringSubmatch(message, -1); len(match) > 0 {
			r := match[0]

			if len(r) >= 7 {
				return NewErrForeignKeyConstraint(r[2], r[3], r[4], r[6])
			}
		}
	}

	return err
}
