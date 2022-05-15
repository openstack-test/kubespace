package alert

import (
	utlAlert "kubespace/server/utils/alert"
	"fmt"
)

//send alert if rule is triggered
func Sender(SendClass map[string][]*utlAlert.Ready2Send, now string) {
	for k, v := range SendClass {
		//fmt.Println("v===>", v)
		switch k {
		case utlAlert.AlertMethodMail:
			go Mail(v, now, 2) //发送邮件
		case utlAlert.AlertMethodWorkChat:
			go WorkChat(v, 2)
		default:
			fmt.Println("default")
		}
	}
}

/*
 send recovery message if alert recovered.
*/
func RecoverSender(SendClass map[string]map[[2]int64]*utlAlert.Ready2Send, now string) {
	for k := range SendClass {
		var hook []*utlAlert.Ready2Send
		for _, u := range SendClass[k] {
			hook = append(hook, u)
		}
		//fmt.Println("methods ==>", k)
		//fmt.Println("hook==>", hook)
		//t := strings.Split(k, " ")
		switch k {
		case utlAlert.AlertMethodWorkChat:
			go WorkChat(hook, 0)
		case utlAlert.AlertMethodMail:
			go Mail(hook, now, 0)
		default:
			fmt.Println("not match !!")
		}
	}
}
