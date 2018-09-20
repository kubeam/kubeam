package common

import (
	"github.com/go-sql-driver/mysql"
	errors "k8s.io/apimachinery/pkg/api/errors"
)

/*ErrorHandler is a generic error handler*/
func ErrorHandler(err error) {
	if err != nil {
		if me, ok := err.(*mysql.MySQLError); !ok {
			if errors.IsNotFound(err) ||
				errors.IsUnauthorized(err) ||
				errors.IsAlreadyExists(err) ||
				errors.IsForbidden(err) {
				LogError.Println(err.Error())
			} else {
				LogError.Println(err.Error())
			}
		} else {
			LogInfo.Println(me.Message, me.Number)
		}
	}
}
