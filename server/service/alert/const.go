package alert

import (
	"bytes"
	utlAlert "kubespace/server/utils/alert"

	"github.com/Masterminds/sprig"
	"github.com/pkg/errors"
	"text/template"
	"time"
)

type MailOptions struct {
	Title       string            `json:"title"`
	Status      int64             `json:"status"`
	Summary     string            `json:"summary"`
	RuleId      int64             `json:"rule_id"`
	Count       int               `json:"count"`
	Value       float64           `json:"value"`
	PromName    string            `json:"prom_name"`
	Labels      map[string]string `json:"labels"`
	Description string            `json:"description"`
	FirstAt     time.Time         `json:"first_at"`
	ResolvedAt  time.Time         `json:"resolved_at"`
	ConformUrl  string            `json:"conform_url"`
}

type WeChatOptions struct {
	RuleId     int64                   `json:"rule_id"`
	Title      string                  `json:"title"`
	Alerts     []*utlAlert.SingleAlert `json:"alerts"`
	Status     int64                   `json:"status"`
	PromName   string                  `json:"prom_name"`
	ConformUrl string                  `json:"conform_url"`
}

type Request struct {
	MsgType  string `json:"msgType"`
	Markdown struct {
		Content string `json:"content"`
	} `json:"markdown"`
}

// 企业微信告警模板
var WorkChartCallTemplate = `
## {{ if eq $.Status 2 }}<font color=\"warning\">{{else}}<font color=\"info\">{{end}}[{{ if eq $.Status 0 }}告警恢复{{ else }}告警触发{{ end }}] {{ $.Title }}</font>
{{ range $key, $value := $.Alerts }}
>**环境:** {{ $.PromName }}
>**标签**
{{ range $k, $v := $value.Labels }}<font color=\"comment\"> {{ $k }}: {{ $v }} |</font>{{ end }} 
>**说明**
<font color=\"#dd0000\"> {{ $value.Summary }} </font>
<font color=\"#dd0000\"> {{ $value.Description }} </font>
>**持续时间**
<font color=\"comment\"> {{ $value.Count }}分钟</font>
>**告警级别**
<font color=\"comment\"> {{ $value.AlarmLevel }} </font>
>**时间**
<font color=\"comment\">{{ if eq $.Status 2 }}触发时间: {{ $value.FirstAt.Format "2006-01-02 15:04:05" }} {{ else }} 触发时间: {{ $value.FirstAt.Format "2006-01-02 15:04:05" }}; 恢复时间: {{ $value.ResolvedAt.Format "2006-01-02 15:04:05" }} {{ end }}</font>
{{ if eq $.Status 2 }}
>**告警确认:** 
[http://{{ $.ConformUrl }}{{ $value.Id }}]( http://{{ $.ConformUrl }}{{ $value.Id }} )
{{ end }}
{{ end }}
`

func ParseString(strtmpl string, obj interface{}) ([]byte, error) {
	var buf bytes.Buffer
	tmpl, err := template.New("template").Funcs(sprig.TxtFuncMap()).Parse(strtmpl)
	if err != nil {
		return nil, errors.Wrap(err, "error when parsing template")
	}
	err = tmpl.Execute(&buf, obj)
	if err != nil {
		return nil, errors.Wrap(err, "error when executing template")
	}
	return buf.Bytes(), nil
}
