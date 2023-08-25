package jwt

import (
	"context"
	"fmt"
	"simpleTiktok/biz/dal/db"
	"simpleTiktok/biz/model/basic/user"
	"simpleTiktok/pkg/errno"
	"simpleTiktok/pkg/utils"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/jwt"
)

var JwtMiddleware *jwt.HertzJWTMiddleware
var identity = "user_id"

func InitJWT() {
	JwtMiddleware, _ = jwt.New(&jwt.HertzJWTMiddleware{
		Key:         []byte("secret key"),
		Timeout:     24 * time.Hour,
		MaxRefresh:  24 * time.Hour,
		IdentityKey: identity,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(int64); ok {
				return jwt.MapClaims{
					identity: v,
				}
			}
			return jwt.MapClaims{}
		},
		Authenticator: func(ctx context.Context, c *app.RequestContext) (interface{}, error) { //登录验证
			var req user.DouyinUserLoginRequest
			if err := c.BindAndValidate(&req); err != nil {
				return nil, err
			}
			user, err := db.QueryUserByName(req.Username)
			if err != nil {
				return nil, err
			}
			if success := utils.VerifyPassword(req.Password, user.Password); !success {
				err = errno.AuthorizationFailedErr
				return nil, err
			}
			c.Set("user_id", user.ID)
			c.Set("current_user_id", user.ID)
			return user.ID, nil
		},
		Authorizator: func(data interface{}, ctx context.Context, c *app.RequestContext) bool {
			value, ok := data.(float64)
			if ok {
				current_user_id := int64(value)
				c.Set("current_user_id", current_user_id)
				return true
			}
			return false
		},
		Unauthorized: func(ctx context.Context, c *app.RequestContext, code int, message string) {
			fmt.Println("jwt未通过")
			c.JSON(consts.StatusOK, user.DouyinUserLoginResponse{
				StatusCode: errno.AuthorizationFailedErr.ErrCode,
				StatusMsg:  message,
			})
		},
		LoginResponse: func(ctx context.Context, c *app.RequestContext, code int, token string, time time.Time) {
			c.Set("token", token)
		},
		HTTPStatusMessageFunc: func(e error, ctx context.Context, c *app.RequestContext) string {
			resp := utils.BuildBaseResp(e)
			return resp.StatusMsg
		},
	})
}
