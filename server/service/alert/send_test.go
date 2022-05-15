package alert

import (
	"fmt"
	"testing"
	"time"
)

type Options struct {
	Status      int64             `json:"status"`
	Summary     string            `json:"summary"`
	PromName    string            `json:"promName"`
	Count       int64             `json:"count"`
	Value       float64           `json:"value"`
	RuleId      int64             `json:"ruleId"`
	Labels      map[string]string `json:"labels"`
	Description string            `json:"description"`
	FirstAt     time.Time         `json:"firstAt"`
	ResolvedAt  time.Time         `json:"resolvedAt"`
	ConformUrl  string            `json:"conformUrl"`
	Title       string            `json:"title"`
}

var st = `
>**环境:** 测试环境
>**标签**
<font color=\"comment\"> collection: job_config_master |</font><font color=\"comment\"> database: mk_master |</font><font color=\"comment\"> endpoint: metrics |</font><font color=\"comment\"> instance: 192.168.0.207-7210 |</font><font color=\"comment\"> job: mongo-rs0-monitor-192.168.0.207-7210 |</font><font color=\"comment\"> namespace: monitoring |</font><font color=\"comment\"> pod: mongo-rs0-monitor-prometheus-mongodb-exporter-5c5cc8468b-76p6j |</font><font color=\"comment\"> service: mongo-rs0-monitor |</font><font color=\"comment\"> type: Total |</font> 
>**说明**
<font color=\"#dd0000\"> 每分钟读写时间超过 5000ms 的热点表, 可能造成对应业务功能体验变差 </font>
<font color=\"#dd0000\"> '192.168.0.207-7210 节点的 mk_master.job_config_master 数据表 每分钟读写时间 25266.192704600082ms' </font>
>**持续时间**
<font color=\"comment\"> 3分钟</font>
>**规则ID**
<font color=\"comment\"> 107 </font>
>**时间**
<font color=\"comment\">触发时间: 2021-08-17 11:40:34 </font>

>**告警确认:** 
[http://localhost:8888/confirm?id=1705]( http://localhost:8888/confirm?id=1705 )
`

func TestWorkChat(t *testing.T) {
	url := "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=f5a2dbfd-986a-461f-85f9-4f3f52bb7de9"
	fmt.Println(len(st))
	err := CallWorkWeChat(url, st)
	fmt.Println("err ==>", err)
	//fmt.Println(5 % 15)

}
