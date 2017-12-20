package mysql

import "fmt"

//ForeignKeyCheck returns  `"SET FOREIGN_KEY_CHECKS = %s"`
func ForeignKeyCheck(turnOn bool) string {
	v := "0"
	if turnOn {
		v = "1"
	}
	return fmt.Sprintf("SET foreign_key_checks = %s", v)
}
