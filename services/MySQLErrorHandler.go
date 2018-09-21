package services

import (
	"github.com/go-sql-driver/mysql"
	"github.com/kubeam/kubeam/common"
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
				common.LogError.Println(err.Error())
			} else {
				common.LogError.Println(err.Error())
			}
		} else {
			common.LogInfo.Println(me.Message, me.Number)
		}
	}
}
