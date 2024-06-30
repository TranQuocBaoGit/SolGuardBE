package models

import "encoding/json"

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
	SourceMap   interface{} `json:"sourceMap"`
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
	Type               string      			`json:"type"`
	Name               string      			`json:"name"`
	SourceMapping      SourceMapping 		`json:"source_mapping,omitempty"`
	TypeSpecificFields TypeSpecificFields 	`json:"type_specific_fields,omitempty"`
}

type SourceMapping struct {
	Start 				int 	`json:"start"`
	Length  			int 	`json:"length"`
	FileNameRelative 	string 	`json:"filename_relative"`
	FileNameAbsolute 	string 	`json:"filename_absolute"`
	FileNameShort 		string 	`json:"filename_short"`
	IsDependancy 		bool 	`json:"is_dependency"`
	Lines 				[]int 	`json:"lines"`
	StartingColumn 		int 	`json:"starting_column"`
	EndingColumn 		int 	`json:"ending_column"`
}

type TypeSpecificFields struct {
	Parent 		interface{} `json:"parent"`
	Signature 	string 		`json:"signature"`
}

// SOLHINT
type SolhintResultDetail []SolhintDetail

type SolhintDetail struct {
	Issues      *SolhintIssue 		`json:"issues"`
	Conclusion  *SolhintConclusion  `json:"conclusion,omitempty"`
}

type SolhintConclusion struct {
	Conclusion string `json:"conclusion,omitempty"`
}

type SolhintIssue struct {
	Line     int    	`json:"line"`
	Column   int    	`json:"column"`
	Severity string 	`json:"severity"`
	Message  string 	`json:"message"`
	RuleID   string 	`json:"ruleId"`
	// Fix      *string 	`json:"fix,omitempty"`
	FilePath string 	`json:"filePath"`
}

func (i *SolhintDetail) UnmarshalJSON(data []byte) error {
	var temp map[string]interface{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	if _, exists := temp["conclusion"]; exists {
		var c SolhintConclusion
		if err := json.Unmarshal(data, &c); err != nil {
			return err
		}
		i.Conclusion = &c
	} else {
		var w SolhintIssue
		if err := json.Unmarshal(data, &w); err != nil {
			return err
		}
		i.Issues = &w
	}
	return nil
}

// HONEYBADGER
type HoneyBadgerResultDetail struct {
	ExecutionPaths 			string 		`json:"execution_paths"`
	MoneyFlow 				bool 		`json:"money_flow"`
	HiddenStateUpdate 		interface{} `json:"hidden_state_update"`
	BalanceDisorder 		interface{} `json:"balance_disorder"`
	ExecutionTime 			string 		`json:"execution_time"`
	HiddenTransfer 			interface{} `json:"hidden_transfer"`
	AttackMethods 			[]string 	`json:"attack_methods"`
	CashoutMethods 			[]string 	`json:"cashout_methods"`
	EVMCodeCoverage 		string 		`json:"evm_code_coverage"`
	StrawManContract 		interface{} `json:"straw_man_contract"`
	SkipEmptyStringLiteral 	interface{} `json:"skip_empty_string_literal"`
	InheritanceDisorder 	interface{} `json:"inheritance_disorder"`
	UninitialisedStruct 	interface{} `json:"uninitialised_struct"`
	TimeOut					interface{} `json:"timeout"`
	DeadCode 				[]string 	`json:"dead_code"`
	TypeDeductionOverflow 	interface{} `json:"type_deduction_overflow"`
}