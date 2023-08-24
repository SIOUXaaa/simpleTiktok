package service

import (
	"context"
	"simpleTiktok/biz/model/basic/user"
	"simpleTiktok/biz/model/common"
	"simpleTiktok/pkg/constants"
	"simpleTiktok/pkg/errno"
	"simpleTiktok/pkg/utils"

	"simpleTiktok/biz/dal/db"

	"github.com/cloudwego/hertz/pkg/app"
)

type UserService struct {
	ctx context.Context
	c   *app.RequestContext
}

func NewUserService(ctx context.Context, c *app.RequestContext) *UserService {
	return &UserService{ctx, c}
}

func (s *UserService) UserRegister(req *user.DouyinUserRegisterRequest) (user_id int64, err error) {
	user, err := db.QueryUserByName(req.Username)
	if err != nil {
		return 0, err
	}

	if *user != (db.User{}) {
		return 0, errno.UserAlreadyExistErr
	}

	password, _ := utils.Crypt(req.Password)
	user_id, _ = db.CreateUser(&db.User{
		UserName:        req.Username,
		Password:        password,
		Avatar:          constants.TestAva,
		BackgroundImage: constants.TestBackground,
		Signature:       constants.TestSign,
	})
	return user_id, nil
}

func (s *UserService) UserInfo(req *user.DouyinUserRequest) (*common.User, error) {
	return nil, nil
}

//todo: GetUserInfo
