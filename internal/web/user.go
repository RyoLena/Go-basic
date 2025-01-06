package web

import (
	"Project/webBook_git/internal/domain"
	"Project/webBook_git/internal/service"
	"errors"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

type UserHandle struct {
	svc         *service.UserService
	emilExp     *regexp.Regexp
	passwordExp *regexp.Regexp
}

func NewUserHandle(svc *service.UserService) *UserHandle {
	const (
		//正则表达式 简单的邮箱验证以及至少需要八位且含有一个特殊字符的密码验证
		emailRegexPattern    = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		passwordRegexPattern = `^(?=.*[!@#$%^&*()_+\-=$$$${};':"\\|,.<>/?])(?=.*[a-zA-Z0-9]).{9,}$`
	)
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	return &UserHandle{
		svc:         svc,
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
	err = user.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
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

	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			//设置过期时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * 30)),
		},
		Uid: usvc.ID,
		//Email:     usvc.Email,
		UserAgent: ctx.Request.UserAgent(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString([]byte("sUvca2dpn7veAV4odb4xQNwYFV0EescZ"))
	if err != nil {
		ctx.String(http.StatusOK, "jwt加密系统错误")
		return
	}
	ctx.Header("x-jwt-token", tokenStr)

	ctx.JSON(http.StatusOK, gin.H{
		"email":    usvc.Email,
		"password": usvc.Password,
	})
}

type UserClaims struct {
	jwt.RegisteredClaims
	//想要获取什么从这里添加
	Uid int64
	//Email     string
	UserAgent string
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

	claims, ok := c.(*UserClaims)
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
