package public

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/go-pay/bm"
)

// 获取公众号全局唯一后台稳定版接口调用凭据（access_token）
// 公众号文档：https://developers.weixin.qq.com/doc/offiaccount/Basic_Information/getStableAccessToken.html
func (s *SDK) getStableAccessToken() (err error) {
	defer func() {
		if err != nil {
			// reset default refresh internal
			s.RefreshInternal = time.Second * 20
			if s.callback != nil {
				go s.callback("", "", 0, err)
			}
		}
	}()

	path := "/cgi-bin/stable_token"
	body := make(bm.BodyMap)
	body.Set("grant_type", "client_credential").
		Set("appid", s.Appid).
		Set("secret", s.Secret).
		Set("force_refresh", false)
	at := &AccessToken{}
	if _, err = s.doRequestPost(s.ctx, path, body, at); err != nil {
		return
	}
	if at.Errcode != Success {
		err = fmt.Errorf("errcode(%d), errmsg(%s)", at.Errcode, at.Errmsg)
		return
	}
	s.accessToken = at.AccessToken
	s.RefreshInternal = time.Second * time.Duration(at.ExpiresIn)
	if s.callback != nil {
		go s.callback(s.Appid, at.AccessToken, at.ExpiresIn, nil)
	}
	return nil
}

func (s *SDK) goAutoRefreshStableAccessToken() {
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 64<<10)
			buf = buf[:runtime.Stack(buf, false)]
			s.logger.Errorf("public_goAutoRefreshStableAccessTokenJob: panic recovered: %s\n%s", r, buf)
			time.Sleep(time.Second * 3)
			if err := s.getStableAccessToken(); err != nil {
				// 失败就不再自动刷新了
				return
			}
			s.goAutoRefreshStableAccessToken()
		}
	}()
	for {
		// every one hour, request new access token, default 10s
		time.Sleep(s.RefreshInternal / 2)
		err := s.getStableAccessToken()
		if err != nil {
			s.logger.Errorf("get stable access token error, after 10s retry: %+v", err)
			continue
		}
	}
}

// SetPublicAccessTokenCallback set public access token callback listener
func (s *SDK) SetPublicAccessTokenCallback(fn func(appid, accessToken string, expireIn int, err error)) {
	s.callback = fn
}

// GetPublicAccessToken get public access token string
func (s *SDK) GetPublicAccessToken() (at string) {
	return s.accessToken
}

// SetPublicAccessToken set public access token string
func (s *SDK) SetPublicAccessToken(accessToken string) {
	s.accessToken = accessToken
}

// 获取 Access Token
// 微信公众号文档：https://developers.weixin.qq.com/doc/offiaccount/Basic_Information/Get_access_token.html
func GetAccessToken(c context.Context, appid, secret string) (at *AccessToken, err error) {
	uri := HostDefault + "/cgi-bin/token?grant_type=client_credential&appid=" + appid + "&secret=" + secret
	at = &AccessToken{}
	if err = doRequestGet(c, uri, at); err != nil {
		return nil, err
	}
	if at.Errcode != Success {
		return nil, fmt.Errorf("errcode(%d), errmsg(%s)", at.Errcode, at.Errmsg)
	}
	return at, nil
}

// 获取 Stable Access Token
// 微信公众号文档：https://developers.weixin.qq.com/doc/offiaccount/Basic_Information/getStableAccessToken.html
func GetStableAccessToken(c context.Context, appid, secret string, forceRefresh bool) (at *AccessToken, err error) {
	url := HostDefault + "/cgi-bin/stable_token"
	body := make(bm.BodyMap)
	body.Set("grant_type", "client_credential").
		Set("appid", appid).
		Set("secret", secret).
		Set("force_refresh", forceRefresh)
	at = &AccessToken{}
	if err = doRequestPost(c, url, body, at); err != nil {
		return nil, err
	}
	if at.Errcode != Success {
		return nil, fmt.Errorf("errcode(%d), errmsg(%s)", at.Errcode, at.Errmsg)
	}
	return at, nil
}
