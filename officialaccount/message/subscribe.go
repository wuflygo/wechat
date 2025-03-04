package message

import (
	"fmt"

	"github.com/wuflygo/wechat/v2/officialaccount/context"
	"github.com/wuflygo/wechat/v2/util"
)

const (
	subscribeSendURL         = "https://api.weixin.qq.com/cgi-bin/message/subscribe/bizsend"
	subscribeTemplateListURL = "https://api.weixin.qq.com/wxaapi/newtmpl/gettemplate"
	subscribeTemplateAddURL  = "https://api.weixin.qq.com/wxaapi/newtmpl/addtemplate"
	subscribeTemplateDelURL  = "https://api.weixin.qq.com/wxaapi/newtmpl/deltemplate"
)

// Subscribe 订阅消息
type Subscribe struct {
	*context.Context
}

// NewSubscribe 实例化
func NewSubscribe(context *context.Context) *Subscribe {
	tpl := new(Subscribe)
	tpl.Context = context
	return tpl
}

// SubscribeMessage 发送的订阅消息内容
type SubscribeMessage struct {
	ToUser      string                        `json:"touser"`         // 必须, 接受者OpenID
	TemplateID  string                        `json:"template_id"`    // 必须, 模版ID
	Page        string                        `json:"page,omitempty"` // 可选, 跳转网页时填写
	Data        map[string]*SubscribeDataItem `json:"data"`           // 必须, 模板数据
	MiniProgram struct {
		AppID    string `json:"appid"`    // 所需跳转到的小程序appid（该小程序appid必须与发模板消息的公众号是绑定关联关系）
		PagePath string `json:"pagepath"` // 所需跳转到小程序的具体页面路径，支持带参数,（示例index?foo=bar）
	} `json:"miniprogram"` // 可选,跳转至小程序地址
}

// SubscribeDataItem 模版内某个 .DATA 的值
type SubscribeDataItem struct {
	Value string `json:"value"`
}

// Send 发送订阅消息
func (tpl *Subscribe) Send(msg *SubscribeMessage) (err error) {
	var accessToken string
	accessToken, err = tpl.GetAccessToken()
	if err != nil {
		return
	}
	uri := fmt.Sprintf("%s?access_token=%s", subscribeSendURL, accessToken)
	response, err := util.PostJSON(uri, msg)
	if err != nil {
		return
	}
	return util.DecodeWithCommonError(response, "SendSubscribeMessage")
}

// PrivateSubscribeItem 私有订阅消息模板
type PrivateSubscribeItem struct {
	PriTmplID string `json:"priTmplId"` //	添加至帐号下的模板 id，发送订阅通知时所需
	Title     string `json:"title"`     // 模版标题
	Content   string `json:"content"`   // 模版内容
	Example   string `json:"example"`   // 模板内容示例
	SubType   int    `json:"type"`      // 模版类型，2 为一次性订阅，3 为长期订阅
}

type resPrivateSubscribeList struct {
	util.CommonError
	SubscriptionList []*PrivateSubscribeItem `json:"data"`
}

// List 获取私有订阅消息模板列表
func (tpl *Subscribe) List() (templateList []*PrivateSubscribeItem, err error) {
	var accessToken string
	accessToken, err = tpl.GetAccessToken()
	if err != nil {
		return
	}
	uri := fmt.Sprintf("%s?access_token=%s", subscribeTemplateListURL, accessToken)
	var response []byte
	response, err = util.HTTPGet(uri)
	if err != nil {
		return
	}
	var res resPrivateSubscribeList
	err = util.DecodeWithError(response, &res, "ListSubscribe")
	if err != nil {
		return
	}
	templateList = res.SubscriptionList
	return
}

type resSubscribeAdd struct {
	util.CommonError

	TemplateID string `json:"priTmplId"`
}

// Add 添加订阅消息模板
func (tpl *Subscribe) Add(ShortID string, kidList []int, sceneDesc string) (templateID string, err error) {
	var accessToken string
	accessToken, err = tpl.GetAccessToken()
	if err != nil {
		return
	}
	var msg = struct {
		TemplateIDShort string `json:"tid"`
		SceneDesc       string `json:"sceneDesc"`
		KidList         []int  `json:"kidList"`
	}{TemplateIDShort: ShortID, SceneDesc: sceneDesc, KidList: kidList}
	uri := fmt.Sprintf("%s?access_token=%s", subscribeTemplateAddURL, accessToken)
	var response []byte
	response, err = util.PostJSON(uri, msg)
	if err != nil {
		return
	}
	var result resSubscribeAdd
	err = util.DecodeWithError(response, &result, "AddSubscribe")
	if err != nil {
		return
	}
	templateID = result.TemplateID
	return
}

// Delete 删除私有模板
func (tpl *Subscribe) Delete(templateID string) (err error) {
	var accessToken string
	accessToken, err = tpl.GetAccessToken()
	if err != nil {
		return
	}
	var msg = struct {
		TemplateID string `json:"priTmplId"`
	}{TemplateID: templateID}
	uri := fmt.Sprintf("%s?access_token=%s", subscribeTemplateDelURL, accessToken)
	var response []byte
	response, err = util.PostJSON(uri, msg)
	if err != nil {
		return
	}
	return util.DecodeWithCommonError(response, "DeleteSubscribe")
}
