package service

import (
	"context"
	"fmt"
	"simpleTiktok/biz/model/basic/user"
	"simpleTiktok/biz/model/common"
	"simpleTiktok/pkg/constants"
	"simpleTiktok/pkg/errno"
	"simpleTiktok/pkg/utils"
	"sync"

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
	userInfo := &common.User{}
	errChan := make(chan error, 4)
	defer close(errChan)
	var wg sync.WaitGroup
	wg.Add(4)

	userInfo.Id = req.UserId
	go func(){
		flag := true
		user, err := db.QueryUserById(req.UserId)
		if err != nil {
			fmt.Println("get user info error: " + err.Error())
			errChan <- err
			flag = false
		}
		if *user == (db.User{}) {
			errChan <- errno.UserIsNotExistErr
			flag = false
		}
		if flag {
			userInfo.Name = user.UserName
			userInfo.Avatar = user.Avatar
			userInfo.BackgroundImage = user.BackgroundImage
			userInfo.Signature = user.Signature
		}
		wg.Done()
	}()

	go func(){
		favorite_count, err := db.GetFavoriteCount(req.UserId)
		if err != nil {
			errChan <- err
		} else {
			userInfo.FavoriteCount = favorite_count
		}
		wg.Done()
	}()

	go func(){
		total_favorite, err := db.GetTotalFavorited(req.UserId)
		if err != nil {
			errChan <- err
		} else {
			userInfo.TotalFavorited = total_favorite
		}
		wg.Done()
	}()

	go func(){
		work_count, err := db.GetWorkCount(req.UserId)
		if err != nil {
			errChan <- err
		}else{
			userInfo.WorkCount= work_count
		}
		wg.Done()
	}()
	
	wg.Wait()
	select{
		case err := <-errChan:
			return &common.User{}, err
		default:
	}
	userInfo.Id = req.UserId
	userInfo.FollowCount = 0
	userInfo.FollowerCount = 0
	userInfo.IsFollow = false

	return userInfo, nil
}

