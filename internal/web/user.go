package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandle struct {
	emilExp     *regexp.Regexp
	passwordExp *regexp.Regexp
}

func NewUserHandle() *UserHandle {
	const (
		//正则表达式 简单的邮箱验证以及至少需要八位且含有一个特殊字符的密码验证
		emailRegexPattern    = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		passwordRegexPattern = `^(?=.*[!@#$%^&*()_+\-=$$$${};':"\\|,.<>/?])(?=.*[a-zA-Z0-9]).{9,}$`
	)
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	return &UserHandle{
		emilExp:     emailExp,
		passwordExp: passwordExp,
	}
}

func (user *UserHandle) SignalUP(ctx *gin.Context) {
	//1.解析数据
	type SignUp struct {
		Email           string `json:"email"`
		ConfirmPassword string `json:"confirm_password"`
		Password        string `json:"password"`
	}
	var req SignUp
	if err := ctx.Bind(&req); err != nil {
		return
	}
	//2.校验 使用正则表达式
	//2.1 邮箱校验
	ok, err := user.emilExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "邮箱格式错误")
		return
	}
	//2.2密码验证
	if req.ConfirmPassword != req.Password {
		ctx.String(http.StatusOK, "两次密码不一致")
		return
	}
	ok, err = user.passwordExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "密码必须大于八位包含特殊字符")
		return
	}
	//3.数据库操作

	//4.返回结果
	ctx.JSON(200, gin.H{
		"email":            req.Email,
		"password":         req.Password,
		"confirm_password": req.ConfirmPassword,
	})
}

func (user *UserHandle) Login(ctx *gin.Context) {

}

func (user *UserHandle) Edit(ctx *gin.Context) {

}
func (user *UserHandle) Profile(ctx *gin.Context) {

}
