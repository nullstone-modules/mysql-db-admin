package mysql

import (
	"fmt"
	"strings"
)

func QuoteIdentifier(identifier string) string {
	return fmt.Sprintf("`%s`", strings.Replace(identifier, "`", "``", -1))
}
