package check
type State string
type NodeType string
const(
	PASS State = "PASS"
	FAIL State = "FAIL"
	WARN State = "WARN"
	INFO State = "INFO"
	MASTER NodeType = "master"
	NODE NodeType = "node"
	CONTROLPLANE NodeType = "controlplane"
	ETCD NodeType = "etcd"
	POLICIES NodeType = "polices"
)
type Check struct {
	ID            string  `yaml:"id"`
	Text          string   `json:"test_desc"`
	// 审计命令
	Audit         string   `json:"audit"`
	AuditEnv      string  `yaml:"audit_env"`
	AuditConfig   string   `yaml:"audit_config"`
	Type          string   `json:"type"`
	Tests         *tests   `json:"-"`
	Set           bool     `json:"-"`
	Remediation   string    `json:"remediation"`
	TestInfo      []string  `json:"test_info"`
	State         `json:"status"`
	ActualValue    string `json:"actual_value"`
	Scored         bool    `json:"scored"`
	IsMultiple     bool    `yaml:"use_multiple_values"`
	ExpectedResult string   `json:"expected_result"`
	Reason         string   `json:"reason,omitempty"`
	// 审计输出
	AuditOutput    string
	AuditEnvOutput string
	AuditConfigOutput string
	DisableEnvTesting bool
}

type tests struct{
	TestItems []*testItem `yaml:"test_items"`
	BinOp     string        `yaml:"bin_op"`
}
type testItem struct{
	Flag    string
	Env     string
	Path    string
	Output  string
	Value   string
	Set     string
	Compare compare
	isMultipleOutput  bool
	auditUsed string
}
type compare struct{
	Op string
	Value string
}
type testOutput struct{
	ID string
	Audit string
    State State
	ActualResult string
	Remediation string
}

// Controls 读取文件生成的结构体
type Controls struct{
	ID               string     `yaml:"id" json:"id"`
	Version          string     `json:"version"`
	DetectedVersion  string     `json:"detected_version,omitempty"`
	Text             string     `json:"text"`
	Type             NodeType   `json:"node_type"`
	Groups           []*Group   `json:"tests"`
}
type Group struct{
	ID           string `yaml:"id" json:"section"`
	Type         string `yaml:"type" json:"type"`
	Pass         int    `json:"pass"`
	Fail         int    `json:"fail"`
	Warn         int    `json:"warn"`
	Info         int     `json:"info"`
	Text         string   `json:"desc"`
	Checks       []*Check `json:"results"`

}