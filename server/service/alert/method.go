package alert

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-gomail/gomail"
	"html/template"
	"io/ioutil"
	"kubespace/server/global"
	utlAlert "kubespace/server/utils/alert"
	"log"
	"net/http"
	"time"
)

// 邮件通知
func Mail(content []*utlAlert.Ready2Send, sendTime string, status int64) { //status: 0->恢复;2->告警
	if len(content) == 0 {
		return
	}
	for _, i := range content {
		opt := &MailOptions{}
		var tempByte bytes.Buffer
		host := global.GVA_CONFIG.Domain.Host
		var t *template.Template
		var title string
		var err error
		for _, v := range i.Alerts {
			opt.Status = status
			opt.Title = v.Title
			opt.Count = v.Count
			opt.Value = v.Value
			opt.Summary = v.Summary
			opt.PromName = i.PromName
			opt.Labels = v.Labels
			opt.Description = v.Description
			opt.FirstAt = v.FirstAt
			opt.ResolvedAt = v.ResolvedAt
			opt.RuleId = i.RuleId
			opt.ConformUrl = fmt.Sprintf("%s/confirm?id=", host)
			if status == 0 {
				title = "告警恢复"
				t, err = template.ParseFiles("pkg/alert/template/recover.html")
			} else {
				title = "告警触发"
				t, err = template.ParseFiles("pkg/alert/template/alert.html")
			}
			if err != nil {
				log.Fatalf("load email html template error:%v\n", err)
				return
			}
			err = t.Execute(&tempByte, opt)
			if err != nil {
				fmt.Println("email params error!", err)
				return
			}
		}
		SendEmail(tempByte.String(), i.User, title)
	}
}

// SendEmail
func SendEmail(EmailBody string, Emails []string, title string) string {
	serverHost := global.GVA_CONFIG.Email.Host
	serverPort := global.GVA_CONFIG.Email.Port
	fromEmail := global.GVA_CONFIG.Email.From
	Passwd := global.GVA_CONFIG.Email.Secret

	//var SendToEmails []string
	m := gomail.NewMessage()
	if len(Emails) == 0 {
		return "收件人不能为空"
	}
	//log.Println("userList ==>", Emails)
	//收件人,...代表打散列表填充不定参数
	m.SetHeader("To", Emails...)
	// 发件人

	m.SetAddressHeader("From", fromEmail, title)
	// 主题
	m.SetHeader("Subject", title)
	// 正文
	m.SetBody("text/html", EmailBody)
	d := gomail.NewDialer(serverHost, serverPort, fromEmail, Passwd)
	//发送
	err := d.DialAndSend(m)
	if err != nil {
		log.Fatal("[email]", err.Error())
	}
	log.Println("email send ok to ", Emails)
	return ""
}

// 企微通知
func WorkChat(content []*utlAlert.Ready2Send, status int64) { //status: 0->恢复;2->告警
	if len(content) == 0 {
		return
	}
	host := global.GVA_CONFIG.Domain.Host
	for _, i := range content {
		var err error
		var res []byte
		opt := &WeChatOptions{}
		opt.PromName = i.PromName
		opt.Title = i.Alerts[0].Title
		opt.ConformUrl = fmt.Sprintf("%s/confirm?id=", host)
		opt.Alerts = i.Alerts
		if len(i.Alerts) > 3 {
			opt.Alerts = i.Alerts[:3]   // 等同于 i.Alerts[0:3],即获取索引为0 1 2的数据
		}
		opt.Status = status
		opt.RuleId = i.RuleId

		res, err = ParseString(WorkChartCallTemplate, opt)
		if err != nil {
			fmt.Println(err)
			return
		}
		// 发送告警
		//log.Println("Send Message ====>", string(res))
		for _, k := range i.CallUrl {
			err := CallWorkWeChat(k, string(res))
			if err != nil {
				log.Println("send wechat err", err)
				return
			}
		}
	}
}

// send WorkWeChat
func CallWorkWeChat(url, contents string) error {
	fmt.Println("len content==>", len(contents))
	data := fmt.Sprintf(`{"msgtype": "markdown","markdown": {"content": "%s"}}`, contents)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		fmt.Println("http request error", err)
		return err
	}
	//req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 5 * time.Second}
	res, err := client.Do(req)
	defer res.Body.Close()
	if err != nil {
		fmt.Println("post work WeChat error")
		return err
	}
	if res.StatusCode != 200 {
		return errors.New("call wechat response.code != 200")
	}
	d, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(d))
	return err
}
