package alert

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"kubespace/server/global"
	"kubespace/server/model/common/response"
	utlAlert "kubespace/server/utils/alert"
	"kubespace/server/utils/monitoring"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultStep   = 10 * time.Minute
	DefaultFilter = ".*"
	DefaultOrder  = "desc"
	DefaultPage   = 1
	DefaultLimit  = 5

	ComponentEtcd      = "etcd"
	ComponentAPIServer = "apiserver"
	ComponentScheduler = "scheduler"

	ErrNoHit           = "'end' must be after the namespace creation time."
	ErrParamConflict   = "'time' and the combination of 'start' and 'end' are mutually exclusive."
	ErrInvalidStartEnd = "'start' must be before 'end'."
	ErrInvalidPage     = "Invalid parameter 'page'."
	ErrInvalidLimit    = "Invalid parameter 'limit'."
)

type reqParams struct {
	time             string
	start            string
	end              string
	step             string
	target           string
	order            string
	page             string
	limit            string
	metricFilter     string
	resourceFilter   string
	nodeName         string
	workspaceName    string
	namespaceName    string
	workloadKind     string
	workloadName     string
	podName          string
	containerName    string
	pvcName          string
	storageClassName string
	componentType    string
	expression       string
	metric           string
	timers           string
}

type queryOptions struct {
	metricFilter string
	namedMetrics []string

	start  time.Time
	end    time.Time
	time   time.Time
	step   time.Duration
	timers string

	target     string
	identifier string
	order      string
	page       int
	limit      int

	option monitoring.QueryOption
}

func (q queryOptions) isRangeQuery() bool {
	return q.time.IsZero()
}

func (q queryOptions) shouldSort() bool {
	return q.target != "" && q.identifier != ""
}

func parseRequestParams(g *gin.Context) reqParams {
	var r reqParams
	r.time = g.DefaultQuery("time", "")
	r.start = g.DefaultQuery("start", "")
	r.end = g.DefaultQuery("end", "")
	r.step = g.DefaultQuery("step", "")
	r.target = g.DefaultQuery("sort_metric", "")
	r.order = g.DefaultQuery("sort_type", "")
	r.page = g.DefaultQuery("page", "")
	r.limit = g.DefaultQuery("limit", "")
	r.metricFilter = g.DefaultQuery("metrics_filter", "")
	r.resourceFilter = g.DefaultQuery("resources_filter", "")
	r.nodeName = g.Query("node")
	r.workspaceName = g.Param("workspace")
	r.namespaceName = g.Param("namespace")
	r.workloadKind = g.Param("kind")
	r.workloadName = g.Param("workload")
	r.timers = g.DefaultQuery("times", "")
	r.podName = g.Param("pod")
	r.containerName = g.Param("container")
	r.pvcName = g.Param("pvc")
	r.storageClassName = g.Param("storageclass")
	r.componentType = g.Param("component")
	r.expression = g.DefaultQuery("expr", "")
	r.metric = g.DefaultQuery("metric", "")
	return r
}

func makeQueryOptions(r reqParams, lvl monitoring.Level) (q queryOptions, err error) {
	if r.resourceFilter == "" {
		r.resourceFilter = DefaultFilter
	}

	q.metricFilter = r.metricFilter
	if r.metricFilter == "" {
		q.metricFilter = DefaultFilter
	}

	switch lvl {
	case monitoring.LevelNode:
		q.identifier = monitoring.IdentifierNode
		q.namedMetrics = monitoring.NodeMetrics
		q.timers = r.timers
		q.option = monitoring.NodeOption{
			ResourceFilter: r.resourceFilter,
			Instance:       r.nodeName,
		}
	}

	// Parse time params
	if r.start != "" && r.end != "" {
		startInt, err := strconv.ParseInt(r.start, 10, 64)
		if err != nil {
			return q, err
		}
		q.start = time.Unix(startInt, 0)

		endInt, err := strconv.ParseInt(r.end, 10, 64)
		if err != nil {
			return q, err
		}
		q.end = time.Unix(endInt, 0)

		if r.step == "" {
			q.step = DefaultStep
		} else {
			q.step, err = time.ParseDuration(r.step)
			if err != nil {
				return q, err
			}
		}

		if q.start.After(q.end) {
			return q, errors.New(ErrInvalidStartEnd)
		}
	} else if r.start == "" && r.end == "" {
		if r.time == "" {
			q.time = time.Now()
		} else {
			timeInt, err := strconv.ParseInt(r.time, 10, 64)
			if err != nil {
				return q, err
			}
			q.time = time.Unix(timeInt, 0)
		}
	} else {
		return q, errors.Errorf(ErrParamConflict)
	}

	// Parse sorting and paging params
	if r.target != "" {
		q.target = r.target
		q.page = DefaultPage
		q.limit = DefaultLimit
		q.order = r.order
		if r.order != monitoring.OrderAscending {
			q.order = DefaultOrder
		}
		if r.page != "" {
			q.page, err = strconv.Atoi(r.page)
			if err != nil || q.page <= 0 {
				return q, errors.New(ErrInvalidPage)
			}
		}
		if r.limit != "" {
			q.limit, err = strconv.Atoi(r.limit)
			if err != nil || q.limit <= 0 {
				return q, errors.New(ErrInvalidLimit)
			}
		}
	}

	return q, nil
}

func handleNameMetricsQuery(c *gin.Context, q queryOptions) {
	cli := global.Monitor
	var res monitoring.Metrics
	var metrics []string
	for _, metric := range q.namedMetrics {
		ok, _ := regexp.MatchString(q.metricFilter, metric)
		if ok {
			metrics = append(metrics, metric)
		}
	}
	if len(metrics) == 0 {
		response.MonitData(true, "OK", res, 0, c)
	}

	if q.isRangeQuery() {
		//fmt.Println("metric===>", metrics)
		res.Results = cli.GetNamedMetricsOverTime(metrics, q.start, q.end, q.step, q.option)
	} else {
		res.Results = cli.GetNamedMetrics(metrics, q.time, q.option)
		if q.shouldSort() {
			res = *res.Sort(q.target, q.order, q.identifier).Page(q.page, q.limit)
		}
	}
	for _, v := range res.Results {
		if strings.Contains(v.MetricName, "processTop_useCpuRate") {
			fmt.Println("metrics==>", v.MetricName)
			ds := parseTime(q.start, q.end, q.step, q.timers)
			for _, vl := range v.MetricData.MetricValues {
				for k, v2 := range ds {
					d := vl.Series[k].Timestamp()
					if ds[k] == d {
						continue
					} else {
						if len(vl.Series) < k {
							//小于k,证明需要补全数据
							vl.Series[k] = monitoring.Point{v2, 0}
						} else {
							//不小于k, 证明时间轴错开了; 重新整理时间轴
							vl.Series[k] = monitoring.Point{v2, vl.Series[k].Value()}
						}
					}
				}
				fmt.Println("value==>", len(vl.Series))
				//	fmt.Println(string(d))
			}
		}
	}
	response.MonitData(true, "OK", res, len(res.Results), c)
}

// 查询自定义数据
func handleNameCustMetricsQuery(c *gin.Context, q queryOptions) {
	cli := global.Monitor
	var res monitoring.Metrics
	var metrics []string
	for _, metric := range q.namedMetrics {
		ok, _ := regexp.MatchString(q.metricFilter, metric)
		if ok {
			metrics = append(metrics, metric)
		}
	}
	if len(metrics) == 0 {
		response.MonitData(true, "OK", res, 0, c)
	}

	if q.isRangeQuery() {
		//fmt.Println("metric===>", metrics)
		res.Results = cli.GetNamedMetricsOverTime(metrics, q.start, q.end, q.step, q.option)
	} else {
		res.Results = cli.GetNamedMetrics(metrics, q.time, q.option)
		if q.shouldSort() {
			res = *res.Sort(q.target, q.order, q.identifier).Page(q.page, q.limit)
		}
	}
	response.MonitData(true, "OK", res, len(res.Results), c)
}

// 生成时间戳序列
func parseTime(start, end time.Time, step time.Duration, timer string) []float64 {
	var dates []float64
	count, _ := strconv.Atoi(timer)
	for i := 0; i <= count; i++ {
		dates = append(dates, float64(start.Unix())+float64(i)*step.Seconds())
	}
	return dates
}

// monitor center helper
type alertForQuery struct {
	*utlAlert.Alert
	label    string
	hostname string
	ruleId   int64
	firedAt  time.Time
}

type users struct {
	User string `json:"user"`
}

type Call struct {
	CallUrl string `json:"call_url"`
}

type mail struct {
	Email string `json:"email"`
}

/*
 set value for fields in alertForQuery
*/
func (a *alertForQuery) setFields() {
	var orderKey []string
	var labels []string

	// set ruleId
	a.ruleId, _ = strconv.ParseInt(a.Annotations.RuleId, 10, 64)
	for key := range a.Labels {
		orderKey = append(orderKey, key)
	}
	sort.Strings(orderKey)
	for _, i := range orderKey {
		labels = append(labels, i+"\a"+a.Labels[i])
	}
	// set label
	a.label = strings.Join(labels, "\v")
	// set firedAt
	a.firedAt = a.FiredAt.Truncate(time.Second)
	// set hostname
	a.setHostname()
}

/*
 set hostname by instance label on data
*/
func (a *alertForQuery) setHostname() {
	h := ""
	if _, ok := a.Labels["instance"]; ok {
		h = a.Labels["instance"]
		boundary := strings.LastIndex(h, ":")
		if boundary != -1 {
			h = h[:boundary]
		}
	}
	a.hostname = h
}

type RecoverInfo struct {
	Id       int64
	Count    int
	Hostname string
}

type planId struct {
	PlanId  int64
	Summary string
}

func RecoverAlert(a alertForQuery, cache map[int64][]utlAlert.UserGroup) {
	// 开启事务
	//err := global.GDB.Transaction(context.TODO(), func(ctx context.Context, tx *gdb.TX) error {
	err := global.GVA_DB.Transaction(func(ctx *gorm.DB) error {
		var receinfo *RecoverInfo
		//err := tx.Ctx(ctx).Model("monitor_alert").Where("rule_id = ? and labels = ? and fired_at = ?", a.ruleId, a.label, a.firedAt).Fields("id,count,hostname").LockUpdate().Scan(&receinfo)
		err := global.GVA_DB.Model("monitor_alert").Where("rule_id = ? and labels = ? and fired_at = ?", a.ruleId, a.label, a.firedAt).Select("id,count,hostname").Scan(&receinfo).Error
		if err == nil {
			if receinfo.Id != 0 {
				// update alert state
				/*_, err := tx.Ctx(ctx).Model("monitor_alert").Data(gdb.Map{
					"status":      a.State,
					"summary":     a.Annotations.Summary,
					"description": a.Annotations.Description,
					"value":       a.Value,
					"resolved_at": a.ResolvedAt,
					"send_count":  0,
				}).Where("id = ?", receinfo.Id).Update()*/
				err := global.GVA_DB.Exec("UPDATE monitor_alert SET status=?, summary=?,description=?,value=?,resolved_at=?,send_count=0 WHERE id=?",a.State, a.Annotations.Summary, a.Annotations.Description,a.Value,a.ResolvedAt,receinfo.Id).Error
				if err == nil {
					//lock for reading map Maintain
					utlAlert.Rw.RLock()
					if _, ok := utlAlert.Maintain[a.hostname]; !ok {
						userGroupList := ([]utlAlert.UserGroup)(nil)
						var plan *planId
						//根据告警信息,查配置告警的rule信息
						//err := tx.Ctx(ctx).Model("monitor_rule").Where("id =?", a.ruleId).Fields("plan_id,summary").Scan(&plan)
						err := global.GVA_DB.Model("monitor_rule").Where("id =?", a.ruleId).Select("plan_id,summary").Scan(&plan).Error
						if err != nil {
							global.GVA_LOG.Error("query alert for rule error", zap.Any("err", err))
						}
						if _, ok := cache[plan.PlanId]; !ok {
							//如果通知渠道中没有找到plan_id,需要重新查询plan_receiver然后缓存起来
							//err := tx.Ctx(ctx).Model("monitor_plan_receiver").Where("plan_id=?", plan.PlanId).Fields("id,start_time,end_time,start,period,reverse_polish_notation,`group`,method,call_url").Scan(&userGroupList)
							err := global.GVA_DB.Model("monitor_plan_receiver").Where("plan_id=?", plan.PlanId).Select("id,start_time,end_time,start,period,reverse_polish_notation,`group`,method,call_url").Scan(&userGroupList).Error
							if err != nil {
								global.GVA_LOG.Error("get plan_receiver for plan_id error", zap.Any("err", err))
							}
							cache[plan.PlanId] = userGroupList
						}
						// 遍历该告警事件对于的所有告警策略
						for _, element := range cache[plan.PlanId] {
							if !(element.IsValid()) && element.IsOnDuty() {
								continue
							}
							// 如果当前告警的延迟时间不大于告警计划设置的延迟时间,跳过
							if !(receinfo.Count >= element.Start) {
								continue
							}
							if ok := shouldSend(receinfo.Id, a.ruleId, receinfo.Count, element); !ok {
								continue
							}
							//log.Println("element==>", element)
							if element.ReversePolishNotation != "label=ALL" {
								if element.ReversePolishNotation != "" && !utlAlert.CalculateReversePolishNotation(a.Labels, element.ReversePolishNotation) {
									continue
								}
							}
							// merge users
							users := SendAlertsForMail(&utlAlert.ValidUserGroup{
								Group: element.Group,
							})
							call := SendAlertsForWeChat(&utlAlert.ValidUserGroup{
								CallUrl: element.CallUrl,
							})
							fmt.Println("Recover info=======>")
							// update Recover2Send, other goroutines in timer.go will handle it
							utlAlert.UpdateRecovery2Send(element, *a.Alert, users, call, receinfo.Id, receinfo.Count, receinfo.Hostname)
						}
					}
					utlAlert.Rw.RUnlock()
				} else {
					fmt.Println("update alert state error", err)
					return err
				}
			}
			return nil
		} else {
			//if exceed the max waiting time for getting the lock
			/*_, err := global.GDB.Model("monitor_alert").Data(gdb.Map{
				"status":      a.State,
				"summary":     a.Annotations.Summary,
				"description": a.Annotations.Description,
				"value":       a.Value,
				"resolved_at": a.ResolvedAt,
				"send_count":  0,
			}).Where("id =?", receinfo.Id).Update()*/
			err := global.GVA_DB.Exec("UPDATE monitor_alert SET status=?,summary=?,description=?,value=?,resolved_at=?,send_count=0 WHERE id=?",a.State,a.Annotations.Summary,a.Annotations.Description,a.Value,a.ResolvedAt,receinfo.Id).Error
			return err
		}
	})
	fmt.Println("recover error ==>", err)
}

// whether recovery should be send
func shouldSend(alertId, ruleId int64, alertCount int, ug utlAlert.UserGroup) (sendFlag bool) {
	if alertCount-ug.Start > ug.Period { //满足告警条件
		sendFlag = true
	} else {
		if _, ok := utlAlert.RuleCount[[2]int64{ruleId, int64(ug.Start)}]; ok {
			global.GVA_LOG.Debug(fmt.Sprintf("id:%d,ruleCount:%d,count:%d,start:%d,period:%d", alertId, utlAlert.RuleCount[[2]int64{ruleId, int64(ug.Start)}], alertCount, ug.Start, ug.Period))
			if utlAlert.RuleCount[[2]int64{ruleId, int64(ug.Start)}] >= int64(alertCount-ug.Start) {
				global.GVA_LOG.Debug(fmt.Sprintf("[%d] id:%d %v", alertId, (utlAlert.RuleCount[[2]int64{ruleId, int64(ug.Start)}]-int64(alertCount)+int64(ug.Start))%int64(ug.Period), utlAlert.RuleCount[[2]int64{ruleId, int64(ug.Start)}]-((utlAlert.RuleCount[[2]int64{ruleId, int64(ug.Start)}]-int64(alertCount)+int64(ug.Start))/int64(ug.Period))*int64(ug.Period) >= int64(ug.Period)))
				if (utlAlert.RuleCount[[2]int64{ruleId, int64(ug.Start)}]-int64(alertCount)+int64(ug.Start))%int64(ug.Period) == 0 || utlAlert.RuleCount[[2]int64{ruleId, int64(ug.Start)}]-((utlAlert.RuleCount[[2]int64{ruleId, int64(ug.Start)}]-int64(alertCount)+int64(ug.Start))/int64(ug.Period))*int64(ug.Period) >= int64(ug.Period) {
					sendFlag = true
				}
			}
		}
	}
	return
}

// 整理告警发生的用户和用户组
func SendAlertsForMail(group *utlAlert.ValidUserGroup) []string {
	var userList []string
	if group.Group != "" {
		user := (*users)(nil)
		//err := global.GDB.Model("monitor_group").Where("name IN(?)", g.SliceStr{group.Group}).Scan(&user)
		err := global.GVA_DB.Model("monitor_group").Where("name IN(?)", []string{group.Group}).Scan(&user).Error
		if err != nil {
			global.GVA_LOG.Error("get group for users error", zap.Any("err", err))
		}
		email := ([]*mail)(nil)
		//err = global.GDB.Model("user_list").Where("username IN(?)", strings.Split(user.User, ",")).Scan(&email)
		err = global.GVA_DB.Model("user_list").Where("username IN(?)", strings.Split(user.User, ",")).Scan(&email).Error
		if err != nil {
			log.Println("get user_list error", err)
		}
		for _, v := range email {
			userList = append(userList, v.Email)
		}
	}
	return userList
}

// 绑定method
func SendAlertsForWeChat(method *utlAlert.ValidUserGroup) []string {
	var methods []string
	call := (*Call)(nil)
	if method.CallUrl != "" {
		//global.GDB.Model("monitor_method").Where("name IN(?)", g.SliceStr{method.CallUrl}).Scan(&call)
		global.GVA_DB.Model("monitor_method").Where("name IN(?)", []string{method.CallUrl}).Scan(&call)
		//fmt.Println("call all ==>", call.CallUrl)
		methods = append(methods, call.CallUrl)
		//fmt.Println("all call url==>", methods)
	}
	return methods
}

// 时间单位处理
// 如果For单位为s返回2;如果单位为m直接获取值返回;如果是h返回10;
func timeStrconvInt(t string) int {
	tmp := string(t[len(t)-1])
	if tmp == "m" {
		tmpM, _ := strconv.Atoi(t[:len(t)-1])
		return tmpM
	} else if tmp == "h" {
		return 20
	}
	return 2
}
