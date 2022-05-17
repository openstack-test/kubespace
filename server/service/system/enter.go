package system

type ServiceGroup struct {
	JwtService
	ApiService
	MenuService
	UserService
	CasbinService
	InitDBService
	AutoCodeService
	BaseMenuService
	AuthorityService
	SystemConfigService
	AutoCodeHistoryService
	OperationRecordService
	AuthorityBtnService
}
