package db

import (
	"simpleTiktok/pkg/constants"
	"simpleTiktok/pkg/errno"
)

type User struct {
	ID              int64  `json:"id"`
	UserName        string `json:"user_name"`
	Password        string `json:"password"`
	Avatar          string `json:"avatar"`
	BackgroundImage string `json:"background_image"`
	Signature       string `json:"signature"`
}

func (User) TableName() string {
	return constants.UserTableName
}

func CreateUser(user *User) (int64, error) {
	err := DB.Create(user).Error
	if err != nil {
		return 0, err
	}
	return user.ID, err
}

func QueryUserById(user_id int64) (*User, error) {
	var user User
	if err := DB.Model(&User{}).Where("id = ?", user_id).Find(&user).Error; err != nil {
		return nil, err
	}
	if user == (User{}){	//用户不存在
		err := errno.UserIsNotExistErr
		return nil, err
	}
	return &user, nil
}

func QueryUserByName(userName string) (*User, error){
	var user = User{}
	if err := DB.Where("user_name = ?", userName).Find(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func VerifyUser(userName, password string) (int64, error) {
	var user User
	if err := DB.Where("user_name = ? AND password = ?", userName, password).Find(&user).Error; err != nil {
		return 0, err
	}
	if user.ID == 0 {	//密码错误
		err := errno.PasswordIsNotVerified
		return user.ID, err
	}
	return user.ID, nil
}

func CheckUserExitsById(user_id int64) (bool, error) {
	var user User 
	if err := DB.Where("id = ?", user_id).Find(&user).Error; err != nil {
		return false, err 	//查询失败
	}
	if user == (User{}) { 	//用户不存在
		return false, nil
	}
	return true, nil
}
