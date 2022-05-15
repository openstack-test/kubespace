package alert

import (
	"kubespace/server/global"
	"errors"
	"strconv"
	"sync"
	"time"
)

var ErrHttpRequest = errors.New("create HTTP request failed")
var Maintain map[string]bool
var RuleCount map[[2]int64]int64
var Recover2Send = map[string]map[[2]int64]*Ready2Send{
	//"HOOK":   map[[2]int64]*Ready2Send{},
}

var Lock sync.Mutex
var Rw sync.RWMutex

func UpdateRecovery2Send(ug UserGroup, alert Alert, users []string, call []string, alertId int64, alertCount int, hostname string) {
	ruleId, _ := strconv.ParseInt(alert.Annotations.RuleId, 10, 64)
	Lock.Lock()
	defer Lock.Unlock()
	//查询标题
	title := (*RuleScanField)(nil)
	name := (*Names)(nil)
	/*global.GDB.Model("monitor_rule").Where("id = ?", ruleId).Scan(&title)
	global.GDB.Model("monitor_prom").Where("id = ?", global.GDB.Model("monitor_rule").Fields("prom_id").Where("id = ?", ruleId)).Scan(&name)*/
	global.GVA_DB.Model("monitor_rule").Where("id = ?", ruleId).Scan(&title)
	global.GVA_DB.Model("monitor_prom").Where("id = ?", global.GVA_DB.Model("monitor_rule").Select("prom_id").Where("id = ?", ruleId)).Scan(&name)
	//fmt.Println("name===>", name.Name)
	if _, ok := Recover2Send[ug.Method]; !ok {
		Recover2Send[ug.Method] = map[[2]int64]*Ready2Send{{ruleId, ug.Id}: {
			RuleId:   ruleId,
			Start:    ug.Id,
			CallUrl:  call,
			User:     users,
			PromName: name.Name,
			Alerts: []*SingleAlert{{
				Id:          alertId,
				Count:       alertCount,
				Value:       alert.Value,
				Summary:     alert.Annotations.Summary,
				Description: alert.Annotations.Description,
				Hostname:    hostname,
				Title:       title.Title,
				FirstAt:     alert.FiredAt,
				ResolvedAt:  alert.ResolvedAt,
				AlarmLevel:  alert.Annotations.AlarmLevel,
			}},
		}}
	} else {
		if _, ok := Recover2Send[ug.Method][[2]int64{ruleId, ug.Id}]; !ok {
			Recover2Send[ug.Method][[2]int64{ruleId, ug.Id}] = &Ready2Send{
				RuleId:   ruleId,
				Start:    ug.Id,
				User:     users,
				PromName: name.Name,
				CallUrl:  call,
				Alerts: []*SingleAlert{{
					Id:          alertId,
					Count:       alertCount,
					Value:       alert.Value,
					Summary:     alert.Annotations.Summary,
					Description: alert.Annotations.Description,
					Hostname:    hostname,
					Title:       title.Title,
					Labels:      alert.Labels,
					FirstAt:     alert.FiredAt,
					ResolvedAt:  alert.ResolvedAt,
					AlarmLevel:  alert.Annotations.AlarmLevel,
				}},
			}
		} else {
			Recover2Send[ug.Method][[2]int64{ruleId, ug.Id}].Alerts = append(Recover2Send[ug.Method][[2]int64{ruleId, ug.Id}].Alerts, &SingleAlert{
				Id:          alertId,
				Count:       alertCount,
				Value:       alert.Value,
				Summary:     alert.Annotations.Summary,
				Hostname:    hostname,
				Description: alert.Annotations.Description,
				Labels:      alert.Labels,
				Title:       title.Title,
				FirstAt:     alert.FiredAt,
				ResolvedAt:  alert.ResolvedAt,
				AlarmLevel:  alert.Annotations.AlarmLevel,
			})
		}
	}
}

// AuthModel holds information used to authenticate.
type AuthModel struct {
	Username string
	Password string
}

type Res struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

type BrokenList struct {
	Hosts []struct {
		Hostname string `json:"hostname"`
	} `json:"hosts"`
	Error interface{} `json:"error"`
}

type Msg struct {
	Content string   `json:"content"`
	From    string   `json:"from"`
	Title   string   `json:"title"`
	To      []string `json:"to"`
}

type SingleAlert struct {
	Id          int64             `json:"id"`
	Count       int               `json:"count"`
	Value       float64           `json:"value"`
	Summary     string            `json:"summary"`
	Title       string            `json:"title"`
	Hostname    string            `json:"hostname"`
	Status      int               `json:"status"`
	Description string            `json:"description"`
	Labels      map[string]string `json:"labels"`
	FirstAt     time.Time         //记录开始时间
	ResolvedAt  time.Time         //恢复时间
	SendCount   int               `json:"send_count"`
	RuleValue   int               `json:"rule_value"`
	AlarmLevel  string            `json:"alarm_level"`
}

type Ready2Send struct {
	RuleId   int64
	Start    int64
	PromName string
	CallUrl  []string
	User     []string
	Alerts   []*SingleAlert
}

type UserGroup struct {
	Id                    int64
	StartTime             string
	EndTime               string
	Start                 int
	Period                int
	ReversePolishNotation string
	Group                 string
	Method                string
	CallUrl               string `json:"call_url"`
}

/*
 Check if UserGroup is valid.
*/
func (u UserGroup) IsValid() bool {
	//fmt.Println("u.user==>", u.User)
	return u.Group != "" || u.CallUrl != ""
}

/*
 IsOnDuty return if current UserGroup is on duty or not by StartTime & EndTime.
 If the UserGroup is not on duty, alerts should not be sent to them.
*/
func (u UserGroup) IsOnDuty() bool {
	now := time.Now().Format("15:04")
	a := (u.StartTime <= u.EndTime && u.StartTime <= now && u.EndTime >= now)  // 不跨 00:00
	b := (u.StartTime > u.EndTime && (u.StartTime <= now || now <= u.EndTime)) // // 跨 00:00
	return a || b
}

type Alerts []Alert

type Alert struct {
	ActiveAt    time.Time `json:"active_at"`
	Annotations struct {
		Description string `json:"description"`
		Summary     string `json:"summary"`
		RuleId      string `json:"rule_id"`
		AlarmLevel  string `json:"alarm_level"`
	} `json:"annotations"`
	FiredAt    time.Time         `json:"fired_at"`
	Labels     map[string]string `json:"labels"`
	LastSentAt time.Time         `json:"last_sent_at"`
	ResolvedAt time.Time         `json:"resolved_at"`
	State      int               `json:"state"`
	ValidUntil time.Time         `json:"valid_until"`
	Value      float64           `json:"value"`
}

type AlertForShow struct {
	Id              int64             `json:"id,omitempty"`
	RuleId          int64             `json:"rule_id"`
	Labels          map[string]string `json:"labels"`
	Value           float64           `json:"value"`
	Count           int               `json:"count"`
	Status          int8              `json:"status"`
	Summary         string            `json:"summary"`
	Description     string            `json:"description"`
	ConfirmedBy     string            `json:"confirmed_by"`
	FiredAt         *time.Time        `json:"fired_at"`
	ConfirmedAt     *time.Time        `json:"confirmed_at"`
	ConfirmedBefore *time.Time        `json:"confirmed_before"`
	ResolvedAt      *time.Time        `json:"resolved_at"`
	AlarmLevel      string            `json:"alarm_level"`
}

type Confirm struct {
	Duration int
	User     string
	Ids      []int
}

type ValidUserGroup struct {
	CallUrl string
	Group   string
}

