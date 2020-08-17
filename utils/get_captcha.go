package utils

import (
	"github.com/mojocn/base64Captcha"
	"image/color"
)

type Captcha struct {
	Id string
	BS64 string
	Code int
}
var store = base64Captcha.DefaultMemStore

func GetCaptcha() (id string,base64 string, err error) {
	rgbaColor := color.RGBA{0,0,0,0}
	fonts := []string{"wqy-microhei.ttc"}
	// 生成driver，高、宽,背景文字的干扰，画线的条数,背景颜色的指针
	driver := base64Captcha.NewDriverMath(50,140,0,0,&rgbaColor,fonts)


	// 使用前面的store和driver 生成验证码的实例
	captcha := base64Captcha.NewCaptcha(driver,store)

	// 生成验证码
	id, base64, err = captcha.Generate()
	return id,base64,err

}

func VerityCaptcha(id string,ret_captcha string) bool {
	return store.Verify(id,ret_captcha,true)
}

