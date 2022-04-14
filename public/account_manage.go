package public

import (
	"context"

	"github.com/go-pay/wechat-sdk/pkg/bmap"
)

// QRCodeCreate 生成带参数的二维码
//	注意：errcode = 0 为成功
//	注意：expire_seconds 字段不传，代表永久二维码。
//	文档：https://developers.weixin.qq.com/doc/offiaccount/Account_Management/Generating_a_Parametric_QR_Code.html
func (s *SDK) QRCodeCreate(c context.Context, body bmap.BodyMap) (qr *QRCodeRsp, err error) {
	path := "/cgi-bin/qrcode/create?access_token=" + s.Conf.AccessToken
	qr = &QRCodeRsp{}
	if err = s.doRequestPost(c, path, body, qr); err != nil {
		return nil, err
	}
	return
}

// ShortKeyGen 生成短key托管
//	注意：errcode = 0 为成功
//	文档：https://developers.weixin.qq.com/doc/offiaccount/Account_Management/KEY_Shortener.html
func (s *SDK) ShortKeyGen(c context.Context, body bmap.BodyMap) (skg *ShortKeyGenRsp, err error) {
	path := "/cgi-bin/shorten/gen?access_token=" + s.Conf.AccessToken
	skg = &ShortKeyGenRsp{}
	if err = s.doRequestPost(c, path, body, skg); err != nil {
		return nil, err
	}
	return
}

// ShortKeyFetch 获取托管的短key
//	注意：errcode = 0 为成功
//	shortKey：短key
//	文档：https://developers.weixin.qq.com/doc/offiaccount/Account_Management/KEY_Shortener.html
func (s *SDK) ShortKeyFetch(c context.Context, shortKey string) (skf *ShortKeyFetchRsp, err error) {
	path := "/cgi-bin/shorten/fetch?access_token=" + s.Conf.AccessToken
	body := make(bmap.BodyMap)
	body.Set("short_key", shortKey)
	skf = &ShortKeyFetchRsp{}
	if err = s.doRequestPost(c, path, body, skf); err != nil {
		return nil, err
	}
	return
}
