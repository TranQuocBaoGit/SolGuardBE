package models

// MYTHRIL
type MythrilResultDetail struct {
	Error   interface{}           `json:"error"`
	Issues  []mythrilIssuesDetail `json:"issues"`
	Success bool                  `json:"success"`
}

type mythrilIssuesDetail struct {
	Address     int         `json:"address"`
	Code        string      `json:"code"`
	Contract    string      `json:"contract"`
	Description string      `json:"description"`
	FileName    string      `json:"filename"`
	Function    string      `json:"function"`
	LineNo      int         `json:"lineno"`
	MaxGasUsed  int         `json:"max_gas_used"`
	MinGasUsed  int         `json:"min_gas_used"`
	Severity    string      `json:"severity"`
	SourceMap   string      `json:"sourceMap"`
	SwcID       string      `json:"swc-id"`
	Title       string      `json:"title"`
	TxSequence  interface{} `json:"txsequence"`
}

// SLITHER
type SlitherResultDetail struct {
	Success bool                  `json:"success"`
	Error   interface{}           `json:"error"`
	Results SlitherDetectorDetail `json:"results"`
}

type SlitherDetectorDetail struct {
	Detectors []EachSiltherDetector `json:"detectors"`
}

type EachSiltherDetector struct {
	Elements             []SlitherElement `json:"elements"`
	Description          string           `json:"description"`
	Markdown             string           `json:"markdown"`
	FirstMarkdownElement string           `json:"first_markdown_element"`
	ID                   string           `json:"id"`
	Check                string           `json:"check"`
	Impact               string           `json:"impact"`
	Confidence           string           `json:"confidence"`
}

type SlitherElement struct {
	Type               string      `json:"type"`
	Name               string      `json:"name"`
	SourceMapping      interface{} `json:"source_mapping"`
	TypeSpecificFields interface{} `json:"type_specific_fields"`
}

// SOLHINT
type SolhintResultDetail struct {
	Issues      []SolhintIssue `json:"issues"`
	Conclusion  string         `json:"conclusion,omitempty"`
}

type SolhintIssue struct {
	Line     int    `json:"line"`
	Column   int    `json:"column"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
	RuleID   string `json:"ruleId"`
	Fix      string `json:"fix,omitempty"`
	FilePath string `json:"filePath"`
}
