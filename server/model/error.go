package model

import "errors"

// 根据业务逻辑需要，自定义一些错误

var (
	ERROE_USER_NOTEXISTS = errors.New("用户不存在")
	ERROE_USER_EXISTS    = errors.New("用户存在")
	ERROE_USER_PWD       = errors.New("密码不正确")
)
