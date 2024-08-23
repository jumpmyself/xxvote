package logic

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"regexp"
	"time"
	"xxvote/app/model"
	"xxvote/app/tools"
)

type User struct {
	Name         string `json:"name" gorm:"name" form:"name"`
	Password     string `json:"password" gorm:"password" form:"password"`
	CaptchaId    string `json:"captcha_id" form:"captcha_id"`
	CaptchaValue string `json:"captcha_value" form:"captcha_value"`
}

func GetLogin(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login.tmpl", nil)
}

// DoLogin godoc
// @Summary         执行用户登录
// @Description     执行用户登录
// @Tags            login
// @Accept          json
// @Produce         json
// @Param           name   body   User  true  "login User"
// @Success         200   {object}    tools.ECode
// @Router          /login   [post]
func DoLogin(ctx *gin.Context) {
	var user User
	if err := ctx.ShouldBind(&user); err != nil {
		ctx.JSON(http.StatusOK, tools.ECode{
			Data:    0,
			Message: "err.Error()", //有风险

		})
		fmt.Println("user:", user)
		return
	}

	fmt.Printf("user:%+v\n", user)
	if !tools.CaptchaVerify(tools.CaptchaData{
		CaptchaId: user.CaptchaId,
		Data:      user.CaptchaValue,
	}) {
		ctx.JSON(http.StatusOK, tools.ECode{
			Code:    10010,
			Message: "验证码校验失败",
		})
		return
	}

	ret := model.GetUser(user.Name)
	fmt.Println(ret)
	if ret.ID < 1 || ret.Password != user.Password {
		ctx.JSON(http.StatusOK, tools.UserErr)
		return
	}

	ctx.SetCookie("name", user.Name, 3600, "/", "", true, false)
	ctx.SetCookie("Id", fmt.Sprint(ret.ID), 3600, "/", "", true, false)

	_ = model.SetSession(ctx, user.Name, ret.ID)

	ctx.JSON(http.StatusOK, tools.ECode{
		Message: "登录成功",
		Data:    ret,
	})
	return
}

// Logout godoc
// @Summary         执行用户退出
// @Description     执行用户退出
// @Tags            login
// @Accept          json
// @Produce         json
// @Success         200   {object}    tools.ECode
// @Router          /logout   [get]
func Logout(ctx *gin.Context) {
	//ctx.SetCookie("name", "", 3600, "/", "", true, false)
	//ctx.SetCookie("Id", "", 3600, "/", "", true, false)
	_ = model.FlushSession(ctx)
	ctx.Redirect(http.StatusFound, "/login")
}

// CUser 新创建一个结构体
type CUser struct {
	Name      string `json:"name"`
	Password  string `json:"password"`
	Password2 string `json:"password_2"`
}

func CreateUser(ctx *gin.Context) {
	var user CUser
	if err := ctx.ShouldBind(&user); err != nil {
		ctx.JSON(http.StatusOK, tools.ECode{
			Code:    10001,
			Message: err.Error(),
		})
		return
	}
	fmt.Printf("user:%+v", user)

	//encrypt(user.Password)
	//encryptV1(user.Password)
	//encryptV2(user.Password)
	//return

	if user.Name == "" || user.Password == "" || user.Password2 == "" {
		ctx.JSON(http.StatusOK, tools.ParamErr)
		return
	}

	//校验密码
	if user.Password != user.Password2 {
		ctx.JSON(http.StatusOK, tools.ECode{
			Code:    10003,
			Message: "两次密码不同",
		})
		return
	}

	//校验账号是否已经存在  这里有风险，并发
	if olduser := model.GetUser(user.Name); olduser.ID > 0 {
		ctx.JSON(http.StatusOK, tools.ECode{
			Code:    10004,
			Message: "用户名已存在", //这里有风险
		})
		return
	}

	nameLen := len(user.Name)
	passwordLen := len(user.Password)
	if nameLen > 16 || nameLen < 8 || passwordLen > 16 || passwordLen < 8 {
		ctx.JSON(http.StatusOK, tools.ECode{
			Code:    10005,
			Message: "账号或密码大于8小于16",
		})
		return
	}

	//引入正则表达式判断密码不能是纯数字
	regex := regexp.MustCompile(`^[0-9]+$`)
	if regex.MatchString(user.Password) {
		ctx.JSON(http.StatusOK, tools.ECode{
			Code:    10006,
			Message: "密码不能为纯数字",
		})
		return
	}

	newUser := model.User{
		Name:        user.Name,
		Password:    encryptV1(user.Password),
		CreatedTime: time.Now(),
		UpdatedTime: time.Now(),
		Uuid:        tools.GetUUID(),
	}
	if err := model.CreateUser(&newUser); err != nil {
		ctx.JSON(http.StatusOK, tools.ECode{
			Code:    10007,
			Message: "新用户创建失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, tools.OK)
	return
}

// 最基础的版本
func encrypt(pwd string) string {
	hash := md5.New()
	hash.Write([]byte(pwd))
	hashBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)
	fmt.Printf("加密后的密码: %s\n", hashString)

	return hashString
}

func encryptV1(pwd string) string {
	newPwd := pwd + "香香编程" //不能随便起，且不能暴露
	hash := md5.New()
	hash.Write([]byte(newPwd))
	hashBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)
	fmt.Printf("加密后的密码: %s\n", hashString)

	return hashString
}

func encryptV2(pwd string) string {
	//基于blowdish 实现加密 ，简单快速，但有安全风险
	//golang。org/x/crypto/ 中有大量的加密算法
	newPwd, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("密码加密失败:", err)
		return ""
	}
	newPwdStr := string(newPwd)
	fmt.Printf("加密后的密码：%s\n", newPwdStr)
	return newPwdStr
}
