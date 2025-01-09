package service

import (
	"encoding/json"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/oauth"
)

var (
	appId     = "wxbb7b02e8aaffb2e4"
	appSecret = "b81d6464d1a8576b44fcf7f4452795b4"
	token     = "TOKEN"
)

// LoginHandler 微信登录
func LoginHandler(c *gin.Context) {
	code := c.Query("code")
	user, err := getWechatOfficialNickname(appId, appSecret, token, code, cache.NewMemory())
	if err != nil {
		c.HTML(500, err.Error(), nil)
		c.Abort()
	}
	user4json, err := json.Marshal(user)
	if err != nil {
		c.HTML(500, err.Error(), nil)
		c.Abort()
	}
	session := sessions.Default(c)
	session.Set("openid4wechat", user.OpenID)
	session.Save()
	c.HTML(200, string(user4json), nil)
}

func getWechatOfficialNickname(appId, appSecret, token, code string, cache cache.Cache) (userinfo oauth.UserInfo, err error) {
	cfg := &config.Config{
		AppID:     appId,
		AppSecret: appSecret,
		Token:     token,
		Cache:     cache,
	}
	wc := wechat.NewWechat()
	officialAccount := wc.GetOfficialAccount(cfg)
	oauth := officialAccount.GetOauth()
	accessToken, err := oauth.GetUserAccessToken(code)
	if err != nil {
		return
	}
	userinfo, err = oauth.GetUserInfo(accessToken.AccessToken, accessToken.OpenID, "")
	if err != nil {
		return
	}
	return
}
