package user_handler

import (
	"time"
)

// 增加用户Exp 并更新 用户信息
// nHExp 小时数
func addUserExp(uExp int64, nHExp int64) int64 {
	retExp := uExp
	timeNowUnix := time.Now().Unix()
	if uExp < timeNowUnix {
		retExp =
			time.Now().Add(
				time.Duration(nHExp) * time.Hour,
			).Unix()
	} else {
		retExp =
			uExp + nHExp*60*60
	}
	return retExp
}
