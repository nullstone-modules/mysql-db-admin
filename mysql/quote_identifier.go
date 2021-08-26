package mysql

import (
	"fmt"
	"strings"
)

func QuoteIdentifier(identifier string) string {
	return fmt.Sprintf("`%s`", strings.Replace(identifier, "`", "``", -1))
}

func QuoteLiteral(literal string) string {
	return fmt.Sprintf("'%s'", strings.Replace(literal, "'", "''", -1))
}
