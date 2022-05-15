package alert

const (
	AlertMethodSms      = "SMS"
	AlertMethodLanxin   = "LANXIN"
	AlertMethodCall     = "CALL"
	AlertMethodHook     = "HOOK"
	AlertMethodDingTalk = "DINGTALK"
	AlertMethodMail     = "Email"
	AlertMethodWorkChat = "WorkChat"
)

const (
	AlertStatusInactive = 0
	AlertStatusPending  = 1
	AlertStatusFiring   = 2
)

type UpdateAlerts struct {
	Id     int64 `json:"id"`
	Status uint8 `json:"status"`
}

type RuleScanField struct {
	Title string `json:"title"`
	For   string `json:"for"`
	AlarmLevel string `json:"alarm_level"`
}

type Names struct {
	Name string `json:"name"`
}
