package utils

import "github.com/astaxie/beego/context"

func LoginFilter(ctx *context.Context)  {

	// 获取session
	id := ctx.Input.Session("id")
	if id == nil {   // 说明未登录
		ctx.Redirect(302,"/")
	}
}
