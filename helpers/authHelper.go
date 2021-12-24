package helper

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func CheckUserType(c *gin.Context, role string) (err error) {
	userType := c.GetString("user_type")
	err = nil

	if userType != role {
		err = errors.New("Unauthorized to access this resource")
		return err
	}
	return err
}

func MatchUserTypeToUid(c *gin.Context, userid string) (err error) {
	userType := c.GetString("user_type")
	uid := c.GetString("uid")
	err = nil

	if userType == "USER" && uid != userid {
		err = errors.New("Unauthorized to access this resource")
		return err
	}

	err = CheckUserType(c, userType)
	return err
}

func MatchBookid(c *gin.Context, bookid string) (err error) {
	uid := c.GetString("_id")
	err = nil

	if uid != bookid {
		err = errors.New("Not found")
		return err
	}
	return err
}
