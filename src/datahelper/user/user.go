package user

import (
	"fmt"
	"utils/function"
	"datahelper/usercache"
	"errors"
)

var ERROR_PASSWORD_WRONG=errors.New("密码不正确！")

type LoginSuccessData struct {
	Auth    string `json:"auth"`
	Avatar   string`json:"avatar"`
}

type LoginFailData struct {
	FailCode int    `json:"failCode"`
	Msg      string `json:"msg"`
}
//验证用户
func UserValid(uid uint32, passwordm string) (valid bool, e error) {
	//验证
	valid=true
	return
}
func GetUidbyName(name string) (uid uint32, e error) {
	uid=0
	if(name=="wind"){
		uid=331805370
	}
	return
}
func CheckUserForbid(uid uint32) (forbid bool, e error) {
	//验证
	forbid=false
	return
}
func CheckUserState(uid uint32) (state bool, e error) {
	//验证
	state=true
	return
}
func CheckAuth(uid uint32, password string) (ud *usercache.UserDetail, e error) {
	//fmt.Print(utils.Md5String(fmt.Sprintf("%s_%d",password,148360)))
	ud,e=usercache.GetUserDetail(uid)
	if e!=nil{
		return
	}
	if passwordm :=function.Md5String(fmt.Sprintf("%s_%d",password,ud.Salt));ud.Password==passwordm{
		return ud,nil
	}else {
		return nil,ERROR_PASSWORD_WRONG
	}
}

func CreateSuccessResp(ud *usercache.UserDetail) (res map[string]interface{}, e error) {
	res = make(map[string]interface{}, 0)
	res["loginstatus"] = 0
	var sdata LoginSuccessData
	sdata.Auth = fmt.Sprintf("%d_%s", ud.Uid, function.Md5String(fmt.Sprintf("%s|%s", ud.Uid, ud.Password)))
	sdata.Avatar=ud.Avatar
	res["userdata"] = sdata
	return
}

func CreateFailResp(code int, msg string) (res map[string]interface{}) {
	res = make(map[string]interface{}, 0)
	var fdata LoginFailData
	res["loginstatus"] = 1
	fdata.FailCode = code
	fdata.Msg = msg
	res["faildata"] = fdata
	return
}