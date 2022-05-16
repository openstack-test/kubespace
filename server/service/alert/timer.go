package alert

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"kubespace/server/global"
	utlAlert "kubespace/server/utils/alert"
	"log"
	"math"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var Rw sync.RWMutex
var Lock sync.Mutex

type Record struct {
	Id              int64     `json:"id"`
	RuleId          int64     `json:"rule_id"`
	Value           float64   `json:"value"`
	Status          int       `json:"status"`
	Count           int       `json:"count"`
	Summary         string    `json:"summary"`
	Description     string    `json:"description"`
	Title           string    `json:"title"`
	Hostname        string    `json:"hostname"`
	ConfirmedBefore time.Time `json:"confirmed_before"`
	FiredAt         time.Time `json:"fired_at"`
	Labels          string    `json:"labels"`
	SendCount       int       `json:"send_count"`
	RuleValue       int       `json:"rule_value"`
	AlarmLevel      string    `json:"alarm_level"`
}

type RecoverRecord struct {
	Id       int64
	RuleId   int64
	Value    float64
	Count    int
	Summary  string
	Hostname string
}

type mainIds struct {
	Id int64 `json:"id"`
}

type host struct {
	HostName string `json:"hostname"`
}

type PlanId struct {
	PlanId  int64
	Summary string
}

func (r Record) getLabelMap() map[string]string {
	labelMap := make(map[string]string)
	if r.Labels != "" {
		for _, j := range strings.Split(r.Labels, "\v") {
			kv := strings.Split(j, "\a")
			labelMap[kv[0]] = kv[1]
		}
	}
	return labelMap
}

func (r Record) getLabelBoolMap() map[string]bool {
	labelMap := make(map[string]bool)
	if r.Labels != "" {
		for _, j := range strings.Split(r.Labels, "\v") {
			kv := strings.Split(j, "\a")
			k := fmt.Sprintf("%s=%s", kv[0], kv[1])
			labelMap[k] = true
		}
	}
	return labelMap
}

// 告警屏蔽
func UpdateMaintainlist() {
	defer func() {
		if e := recover(); e != nil {
			buf := make([]byte, 16384)
			buf = buf[:runtime.Stack(buf, false)]
			log.Fatalf("Panic in UpdateMaintainlist:%v\n%s", e, buf)
		}
	}()
	delte, _ := time.ParseDuration("30s")
	datetime := time.Now().Add(delte)
	now := datetime.Format("15:04")
	maintainIds := ([]*mainIds)(nil)
	//err := global.GDB.Model("monitor_maintain").Where("valid>=? AND day_start<=? AND day_end>=? AND (flag=true AND (time_start<=? OR time_end>=?) OR flag=false AND time_start<=? AND time_end>=?) AND month&"+strconv.Itoa(int(math.Pow(2, float64(time.Now().Month()))))+">0", datetime.Format("2006-01-02 15:04:05"), datetime.Day(), datetime.Day(), now, now, now, now).Scan(&maintainIds)
	err := global.GVA_DB.Raw("SELECT * FROM monitor_maintain WHERE valid>=? AND day_start<=? AND day_end>=? AND (flag=true AND (time_start<=? OR time_end>=?) OR flag=false AND time_start<=? AND time_end>=?) AND month&"+strconv.Itoa(int(math.Pow(2, float64(time.Now().Month()))))+">0", datetime.Format("2006-01-02 15:04:05"), datetime.Day(), datetime.Day(), now, now, now, now).Scan(&maintainIds).Error
	if err != nil {
		log.Println("get table maintains id list error", err)
	}
	m := make(map[string]bool)
	for _, mid := range maintainIds {
		hosts := ([]*host)(nil)
		//err := global.GDB.Model("monitor_host").Where("mid = ?", mid.Id).Scan(&hosts)
		err := global.GVA_DB.Model("monitor_host").Where("mid = ?", mid.Id).Scan(&hosts)
		if err != nil {
			log.Println("get table host hostname list error", err)
		}
		for _, name := range hosts {
			m[name.HostName] = true
		}
	}
	Rw.Lock()
	utlAlert.Maintain = m
	Rw.Unlock()
}

func Filter(tx *gorm.DB, alerts map[int64][]Record, maxCount map[int64]int) map[string][]*utlAlert.Ready2Send {
	SendClass := make(map[string][]*utlAlert.Ready2Send)
	Cache := make(map[int64][]utlAlert.UserGroup)
	NewResultCount := make(map[[2]int64]int64)
	//fmt.Println("alerts ====>", alerts)
	for key := range alerts {
		//fmt.Println("alerts ====>", key)
		var userGroupList []utlAlert.UserGroup
		planId := (*PlanId)(nil)
		AlertsMap := make(map[int][]*utlAlert.SingleAlert)
		err := tx.Model("monitor_rule").Where("id = ?", key).Scan(&planId)
		if err != nil {
			log.Println("get monitor_rule ids error", err)
			return nil
		}
		if _, ok := Cache[planId.PlanId]; !ok {
			err := tx.Model("monitor_plan_receiver").Where("plan_id = ?", planId.PlanId).Scan(&userGroupList)
			if err != nil {
				log.Fatalf("get monitor_plan_receiver plan_id list error:%v", err)
				return nil
			}
			Cache[planId.PlanId] = userGroupList
		}
		//fmt.Println("len===>", len(Cache))
		for _, element := range Cache[planId.PlanId] {
			if !element.IsValid() || !element.IsOnDuty() {
				break
			}
			if maxCount[key] < element.Start {
				break
			}
			k := [2]int64{key, int64(element.Start)}
			if _, ok := utlAlert.RuleCount[k]; !ok {
				NewResultCount[k] = 0
			} else {
				NewResultCount[k] = 1 + utlAlert.RuleCount[k]
			}
			if NewResultCount[k]%int64(element.Period) != 0 {
				break
			}
			// add alerts to AlertsMap
			if _, ok := AlertsMap[element.Start]; !ok {
				putToAlertMap(AlertsMap, element, alerts[key])
			}
			// forward alerts in AlertsMap to SendClass
			if len(AlertsMap[element.Start]) > 0 {
				var filteredAlerts []*utlAlert.SingleAlert
				if element.ReversePolishNotation == "" {
					filteredAlerts = AlertsMap[element.Start]
				} else {
					for _, alert := range AlertsMap[element.Start] {
						// 先判断优先级; 如果次数为0,优先发送.
						var flag bool = false
						if (alert.SendCount == 0) || (alert.Count%(alert.RuleValue*alert.SendCount)) == 0 {
							// 更新发生次数
							/*_, err = tx.Model("monitor_alert").Data(gdb.Map{
								"send_count": gdb.Counter{
									Field: "send_count",
									Value: 1,
								},
							}).Where("status = 2 and rule_id = ?", key).Update()
							if err != nil {
								fmt.Println("update send_count error", err)
							}*/
							tx.Exec("UPDATE monitor_alert SET send_count=send_count+1 WHERE status = 2 and rule_id = ?", key)
							flag = true
						}
						if flag && element.ReversePolishNotation == "label=ALL" {
							filteredAlerts = append(filteredAlerts, alert)
						} else if flag && utlAlert.CalculateReversePolishNotation(alert.Labels, element.ReversePolishNotation) {
							filteredAlerts = append(filteredAlerts, alert)
						}
					}
				}
				//fmt.Println("filteredAlerts=====>", filteredAlerts)
				putToSendClass(SendClass, key, element, filteredAlerts)
			}
		}
	}
	utlAlert.RuleCount = NewResultCount
	return SendClass
}

func putToAlertMap(alertMap map[int][]*utlAlert.SingleAlert, ug utlAlert.UserGroup, alerts []Record) {
	alertMap[ug.Start] = []*utlAlert.SingleAlert{}
	//fmt.Println("ug ===<>", ug)
	for _, alert := range alerts {
		//fmt.Println("alerts===>", alert)
		if alert.Count >= ug.Start {
			if _, ok := utlAlert.Maintain[alert.Hostname]; !ok {
				alertMap[ug.Start] = append(alertMap[ug.Start], &utlAlert.SingleAlert{
					Id:          alert.Id,
					Count:       alert.Count,
					Value:       alert.Value,
					Summary:     alert.Summary,
					Title:       alert.Title,
					Description: alert.Description,
					Hostname:    alert.Hostname,
					FirstAt:     alert.FiredAt,
					Status:      alert.Status,
					SendCount:   alert.SendCount,
					RuleValue:   alert.RuleValue,
					Labels:      alert.getLabelMap(),
					AlarmLevel:  alert.AlarmLevel,
				})
			}
		}
	}
}

// 绑定通知媒介的相关人和组
func putToSendClass(sendClass map[string][]*utlAlert.Ready2Send, ruleId int64, ug utlAlert.UserGroup, alerts []*utlAlert.SingleAlert) {
	if len(alerts) <= 0 {
		return
	}
	name := (*utlAlert.Names)(nil)
	//global.GDB.Model("monitor_prom").Where("id = ?", global.GDB.Model("monitor_rule").Fields("prom_id").Where("id = ?", ruleId)).Scan(&name)
	global.GVA_DB.Model("monitor_prom").Where("id = ?", global.GVA_DB.Model("monitor_rule").Select("prom_id").Where("id = ?", ruleId)).Scan(&name)
	sendClass[ug.Method] = append(sendClass[ug.Method], &utlAlert.Ready2Send{
		RuleId:   ruleId,
		Start:    ug.Id,
		PromName: name.Name,
		User:     SendAlertsForMail(&utlAlert.ValidUserGroup{Group: ug.Group}),
		CallUrl:  SendAlertsForWeChat(&utlAlert.ValidUserGroup{CallUrl: ug.CallUrl}),
		Alerts:   alerts,
	})
}

func Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				break
			default:
				current := time.Now()
				time.Sleep(time.Duration(90-current.Second()) * time.Second)
				UpdateMaintainlist()
			}
		}
	}()
	// 首次启动将所有正在告警的记录设置为已恢复
	// 如果启动前告警已经恢复,那么通过这个机制可以使已恢复的告警正常
	// 如果启动前告警未恢复,那启动之后告警将作为新的告警重新触发告警.
	/*_, err := global.GDB.Model("monitor_alert").Data(gdb.Map{
		"status":     3,
		"send_count": 0, // 重启系统之后,将所有之前已经是告警状态的告警次数重置为0; 保证重启之后所有的入库的告警都是新告警
	}).Where("status = 2").Update()
	if err != nil {
		log.Println("update alert to recover error", err)
	}*/
	global.GVA_DB.Exec("UPDATE monitor_alert SET status=3, send_count=0 WHERE status=2")
	go func() {
		ticker1 := time.NewTicker(30 * time.Second)
		defer ticker1.Stop()
		for {
			select {
			case <-ctx.Done():
				break
			case <-ticker1.C:
				//current := time.Now()
				//time.Sleep(time.Duration(60-current.Second()) * time.Second)
				now := time.Now().Format("2006-01-02 15:04:05")
				go func() {
					defer func() {
						if e := recover(); e != nil {
							buf := make([]byte, 16384)
							buf = buf[:runtime.Stack(buf, false)]
							log.Fatalf("Panic in timer:%v\n%s", e, buf)
						}
					}()
					info := ([]Record)(nil)
					/*_, err := global.GDB.Model("monitor_alert").Data(gdb.Map{
						"status": 2,
					}).Where("status=1 AND confirmed_before<?", now).Update() //更新确认时间已经过期的告警记录
					if err != nil {
						log.Fatalf("update alert status")
					}*/

					global.GVA_DB.Exec("UPDATE monitor_alert SET status=2 WHERE status=1 AND confirmed_before<?", now)

					/*if tx, err := global.GDB.Begin(); err == nil {
						_, err := tx.Model("monitor_alert").Data(gdb.Map{
							"count": gdb.Counter{
								Field: "count",
								Value: 1,
							},
						}).Where("status IN(?)", g.Slice{1, 2}).Update() //凡是非已恢复的告警记录时间全部增加1
						if err != nil {
							log.Printf("update alert status!=0 error:%v", err)
						}*/

					   global.GVA_DB.Exec("UPDATE monitor_alert SET count=count+1 WHERE status IN(1, 2)")

						/*err = tx.Model("monitor_alert").Where("status = ?", 2).Scan(&info) //查询正在告警的记录
						if err != nil {
							log.Printf("select alert status=2 data error:%v", err)
						}*/

						global.GVA_DB.Raw("SELECT * FROM monitor_alert WHERE status = ?", 2).Scan(&info)

						aggregation := make(map[int64][]Record)
						maxCount := make(map[int64]int)
						for _, i := range info {
							fmt.Printf("info alarm_level[%s], status[%d], ruleId[%d]\n", i.AlarmLevel, i.Status, i.RuleId)
							aggregation[i.RuleId] = append(aggregation[i.RuleId], i)
							if _, ok := maxCount[i.RuleId]; !ok {
								maxCount[i.RuleId] = i.Count
							} else {
								if i.Count > maxCount[i.RuleId] {
									maxCount[i.RuleId] = i.Count
								}
							}
						}
						Rw.RLock()
						tx := global.GVA_DB.Begin()
						ready2send := Filter(tx, aggregation, maxCount)
						Rw.RUnlock()
						tx.Commit()
						log.Println("Alert to send :", ready2send)
						// 发送告警
						Sender(ready2send, now)
						Lock.Lock()
						recover2send := utlAlert.Recover2Send
						utlAlert.Recover2Send = map[string]map[[2]int64]*utlAlert.Ready2Send{
							utlAlert.AlertMethodMail:     {},
							utlAlert.AlertMethodWorkChat: {},
						}
						Lock.Unlock()
						log.Println("Recoveries to send:", recover2send)
						RecoverSender(recover2send, now)

				}()
			}
		}
	}()
}

