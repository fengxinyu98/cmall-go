//Package service ...
/*
 * @Descripttion:
 * @Author: congz
 * @Date: 2020-06-17 14:45:17
 * @LastEditors: congz
 * @LastEditTime: 2020-07-18 14:31:47
 */
package service

import (
	"cmall/model"
	"cmall/pkg/e"
	"cmall/pkg/logging"
	"cmall/serializer"

	"github.com/jinzhu/gorm"
)

// AdminUpdateService 管理员修改密码的服务
type AdminUpdateService struct {
	UserName        string `form:"user_name" json:"user_name" binding:"required,min=5,max=30"`
	Password        string `form:"password" json:"password" binding:"required,min=8,max=40"`
	PasswordConfirm string `form:"password_confirm" json:"password_confirm" binding:"required,min=8,max=40"`
}

// Update 管理员修改密码函数
func (service *AdminUpdateService) Update() serializer.Response {
	var admin model.Admin
	code := e.SUCCESS

	if err := model.DB.Where("user_name = ?", service.UserName).First(&admin).Error; err != nil {
		//如果查询不到，返回相应错误
		if gorm.IsRecordNotFoundError(err) {
			logging.Info(err)
			code = e.ERROR_NOT_EXIST_USER
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
			}
		}
		code = e.ERROR_NOT_EXIST_USER
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	// 验证密码是否相同
	if service.PasswordConfirm != service.Password {
		code = e.ERROR_NOT_COMPARE_PASSWORD
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	// 加密密码
	if err := admin.SetPassword(service.Password); err != nil {
		logging.Info(err)
		code = e.ERROR_FAIL_ENCRYPTION
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	if err := model.DB.Save(&admin).Error; err != nil {
		logging.Info(err)
		code = e.ERROR_DATABASE
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	code = e.UPDATE_PASSWORD_SUCCESS
	//返回修改密码成功信息
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}
