package web

import (
	"Project/internal/domain"
	"Project/internal/service"
	myjwt "Project/internal/web/jwt"
	"errors"
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

type UserHandle struct {
	codeSVC     service.CodeService
	svc         service.UserService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
	myjwt.Handle
}

func NewUserHandle(userSvc service.UserService,
	codeSvc service.CodeService, jwtHdl myjwt.Handle) *UserHandle {
	const (
		//正则表达式 简单的邮箱验证以及至少需要八位且含有一个特殊字符的密码验证
		emailRegexPattern    = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		passwordRegexPattern = `^(?=.*[!@#$%^&*()_+\-=$$$${};':"\\|,.<>/?])(?=.*[a-zA-Z0-9]).{9,}$`
	)
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	return &UserHandle{
		codeSVC:     codeSvc,
		svc:         userSvc,
		emailExp:    emailExp,
		passwordExp: passwordExp,
		Handle:      jwtHdl,
	}
}

func (user *UserHandle) Register(server *gin.Engine) {
	userHandle := server.Group("/user")
	{
		userHandle.POST("/signup", user.SignalUP)
		userHandle.POST("/login", user.LoginJWT)
		userHandle.POST("/edit", user.Edit)
		userHandle.GET("/profile", user.Profile)
		userHandle.POST("/login_sms/code/send", user.SendSMSCode)
		userHandle.POST("/login_sms", user.LoginBySMS)
		userHandle.POST("/refresh_token", user.RefreshToken)
	}
}

func (user *UserHandle) RefreshToken(ctx *gin.Context) {
	RefreshToken := user.ExtractToken(ctx)
	var rcKey myjwt.RefreshClaims
	token, err := jwt.ParseWithClaims(RefreshToken, &rcKey, func(token *jwt.Token) (interface{}, error) {
		return myjwt.RtKey, nil
	})
	if err != nil || !token.Valid {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	err = user.CheckSession(ctx, rcKey.SsID)
	if err != nil {
		_ = ctx.AbortWithError(401, fmt.Errorf("redis出错或者退出登录"))
		return
	}
	err = user.SetJWTToken(ctx, rcKey.Uid, rcKey.SsID)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "ok",
	})
}

func (user *UserHandle) SignalUP(ctx *gin.Context) {
	//1.解析数据
	type SignUp struct {
		Email           string `json:"email"`
		ConfirmPassword string `json:"confirm_password"`
		Phone           string `json:"phone"`
		Password        string `json:"password"`
	}
	var req SignUp
	if err := ctx.Bind(&req); err != nil {
		return
	}
	//2.校验 使用正则表达式
	//2.1 邮箱校验
	ok, err := user.emailExp.MatchString(req.Email)
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
	err = user.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Phone:    req.Phone,
		Password: req.Password,
	})
	if errors.Is(err, service.SVCErrUserDuplicated) {
		ctx.String(http.StatusOK, "邮箱已经被注册")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "svc的问题造成的系统异常")
		return
	}
	//4.返回结果
	ctx.JSON(200, gin.H{
		"email":            req.Email,
		"password":         req.Password,
		"confirm_password": req.ConfirmPassword,
	})
}

func (user *UserHandle) LoginJWT(ctx *gin.Context) {
	type LoginReq struct {
		Email    string
		Password string
	}
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		//不需要填写返回信息，因为 gin绑定错误会自动返回错误信息 400
		return
	}
	//验证 通过一层层下传数据验证
	//调用service层的方法
	usvc, err := user.svc.Login(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		ctx.String(http.StatusOK, "账号或者密码不对")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "login系统错误")
		return
	}

	if err := user.Handle.SetLoginToken(ctx, usvc.ID); err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"email":    usvc.Email,
		"password": usvc.Password,
	})
}

func (user *UserHandle) SendSMSCode(ctx *gin.Context) {
	const biz = "loginOrVerify"

	type Reqs struct {
		Phone string `json:"phone"`
	}
	var req Reqs
	if err := ctx.Bind(&req); err != nil {
		return
	}
	err := user.codeSVC.SendCode(ctx, biz, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
			Data: err,
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Code: 1,
		Msg:  "发送成功",
	})
}

func (user *UserHandle) LoginBySMS(ctx *gin.Context) {
	const biz = "loginOrVerify"
	type Reqs struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Reqs
	if err := ctx.Bind(&req); err != nil {
		return
	}
	err := user.codeSVC.VerifyCode(ctx, biz, req.Phone, req.Code)
	switch {
	case err == nil:
		ctx.JSON(http.StatusOK, Result{
			Code: 1,
			Msg:  "校验成功",
		})
	case errors.Is(err, service.ErrCodeVerifyToManyTimes):
		ctx.JSON(http.StatusOK, Result{
			Code: 3,
			Msg:  "发送太频繁，请稍后再试",
		})
	default:
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
			Data: err,
		})
		return
	}

	//给他注册
	u, err := user.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "FindOrCreate系统错误",
		})
		return
	}
	err = user.Handle.SetLoginToken(ctx, u.ID)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "setJWT系统错误",
		})
		return
	}
	ctx.String(http.StatusOK, "校验成功")
}

func (user *UserHandle) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string
		Password string
	}
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		//不需要填写返回信息，因为 gin绑定错误会自动返回错误信息 400
		return
	}
	//验证 通过一层层下传数据验证
	//调用service层的方法
	usvc, err := user.svc.Login(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		ctx.String(http.StatusOK, "账号或者密码不对")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "login系统错误")
		return
	}

	//保持登录状态
	//用session
	sess := sessions.Default(ctx)
	sess.Set("userID", usvc.ID)
	sess.Options(sessions.Options{
		MaxAge: 30,
	})
	err = sess.Save()
	if err != nil {
		ctx.String(200, "session保存错误")
		return
	}

	//登录校验

	ctx.JSON(http.StatusOK, gin.H{
		"email":    usvc.Email,
		"password": usvc.Password,
	})
}

func (user *UserHandle) Edit(ctx *gin.Context) {

}
func (user *UserHandle) Profile(ctx *gin.Context) {
	type Profile struct {
		ID    int64
		Email string
	}
	var userProfile Profile
	err := ctx.Bind(&userProfile)
	if err != nil {
		return
	}
	c, ok := ctx.Get("claims")

	claims, ok := c.(*myjwt.UserClaims)
	if !ok {
		ctx.String(http.StatusOK, "Claims代码系统错误")
		return
	}
	userProfile.ID = claims.Uid
	//userProfile.Email = claims.Email
	ctx.JSON(http.StatusOK, gin.H{
		"userID": userProfile.ID,
		//"userEmail": userProfile.Email,
	})
}
