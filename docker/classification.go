package docker

import (
	"fmt"
	"getContractDeployment/models"
)

var OverallVulna map[int]string = map[int]string{
	1:  "Reentrancy", //Reentrancy vulnerabilities
	2:  "Arithmetic Overflow and Underflow",
	3:  "Delegatecall", //Controlled Delegatecall   Payable functions using delegatecall inside a loop
	4:  "Signature Replay",
	5:  "Random numbers generation", //Weak PRNG    Block timestamp
	6:  "Private Data",
	7:  "Phishing with tx.origin", //Dangerous usage of tx.origin
	8:  "Hiding Malicious Code with External Contract", // arbitrary-send-eth
	9:  "Honeypots",
	10: "Denial of Service",
}

var MythrilVulnaClass map[string]string = map[string]string{
	"100": "Function Default Visibility",
	"101": "Integer Overflow and Underflow",
	"102": "Outdated Compiler Version",
	"103": "Floating Pragma",
	"104": "Unchecked Call Return Value",
	"105": "Unprotected Ether Withdrawal",
	"106": "Unprotected SELFDESTRUCT Instruction",
	"107": "Reentrancy",
	"108": "State Variable Default Visibility",
	"109": "Uninitialized Storage Pointer",
	"110": "Assert Violation",
	"111": "Use of Deprecated Solidity Functions",
	"112": "Delegatecall to Untrusted Callee",
	"113": "DoS with Failed Call",
	"114": "Transaction Order Dependence",
	"115": "Authorization through tx.origin",
	"116": "Block values as a proxy for time",
	"117": "Signature Malleability",
	"118": "Incorrect Constructor Name",
	"119": "Shadowing State Variables",
	"120": "Weak Sources of Randomness from Chain Attributes",
	"121": "Missing Protection against Signature Replay Attacks",
	"122": "Lack of Proper Signature Verification",
	"123": "Requirement Violation",
	"124": "Write to Arbitrary Storage Location",
	"125": "Incorrect Inheritance Order",
	"126": "Insufficient Gas Griefing",
	"127": "Arbitrary Jump with Function Type Variable",
	"128": "DoS With Block Gas Limit",
	"129": "Typographical Error",
	"130": "Right-To-Left-Override control character (U+202E)",
	"131": "Presence of unused variables",
	"132": "Unexpected Ether balance",
	"133": "Hash Collisions With Multiple Variable Length Arguments",
	"134": "Message call with hardcoded gas amount",
	"135": "Code With No Effects",
	"136": "Unencrypted Private Data On-Chain",
}

func MythrilStandardize(sumup models.SumUp) (string, string) {
	name := sumup.Name
	vulnaDefine := ""
	switch name {
	case "Integer Overflow and Underflow":
		vulnaDefine = OverallVulna[2]
		break
	case "Delegatecall to Untrusted Callee":
		vulnaDefine = OverallVulna[3]
		break
	case "DoS with Failed Call":
	case "DoS With Block Gas Limit":
		vulnaDefine = OverallVulna[10]
		break
	case "Authorization through tx.origin":
		vulnaDefine = OverallVulna[7]
		break
	case "Signature Malleability":
	case "Missing Protection against Signature Replay Attacks":
		vulnaDefine = OverallVulna[4]
		break
	case "Weak Sources of Randomness from Chain Attributes":
		vulnaDefine = OverallVulna[5]
		break
	case "Unencrypted Private Data On-Chain":
		vulnaDefine = OverallVulna[6]
		break
	case "Reentrancy":
		vulnaDefine = OverallVulna[1]
		break
	case "MYTHRIL ERROR":
		vulnaDefine = ""
		break
	default:
		vulnaDefine = name
		break
	}
	fmt.Println("Mythril said: ", vulnaDefine, " with ", sumup.Severity)
	return vulnaDefine, sumup.Severity
}


var SlitherVulnaClass map[string]string = map[string]string{
	"abiencoderv2-array":          "Storage abiencoderv2 array",
	"arbitrary-send-erc20":        "Arbitrary from in transferFrom",
	"array-by-reference":          "Modifying storage array by value",
	"encode-packed-collision":     "ABI encodePacked Collision",
	"incorrect-shift":             "Incorrect shift in assembly.",
	"multiple-constructors":       "Multiple constructor schemes",
	"name-reused":                 "Name reused",
	"protected-vars":              "Protected Variables",
	"public-mappings-nested":      "Public mappings with nested variables",
	"rtlo":                        "Right-to-Left-Override character",
	"shadowing-state":             "State variable shadowing",
	"suicidal":                    "Suicidal",
	"uninitialized-state":         "Uninitialized state variables",
	"uninitialized-storage":       "Uninitialized storage variables",
	"unprotected-upgrade":         "Unprotected upgradeable contract",
	"arbitrary-send-erc20-permit": "Arbitrary from in transferFrom used with permit",
	"arbitrary-send-eth":          "Functions that send Ether to arbitrary destinations",
	"controlled-array-length":     "Array Length Assignment",
	"controlled-delegatecall":     "Controlled Delegatecall",
	"delegatecall-loop":           "Payable functions using delegatecall inside a loop",
	"incorrect-exp":               "Incorrect exponentiation",
	"incorrect-return":            "Incorrect return in assembly",
	"msg-value-loop":              "msg.value inside a loop",
	"reentrancy-eth":              "Reentrancy vulnerabilities",
	"return-leave":                "Return instead of leave in assembly",
	"storage-array":               "Storage Signed Integer Array",
	"unchecked-transfer":          "Unchecked transfer",
	"weak-prng":                   "Weak PRNG",
	"codex":                       "Codex",
	"domain-separator-collision":  "Domain separator collision",
	"enum-conversion":             "Dangerous enum conversion",
	"erc20-interface":             "Incorrect erc20 interface",
	"erc721-interface":            "Incorrect erc721 interface",
	"incorrect-equality":          "Dangerous strict equalities",
	"locked-ether":                "Contracts that lock Ether",
	"mapping-deletion":            "Deletion on mapping containing a structure",
	"shadowing-abstract":          "State variable shadowing from abstract contracts",
	"tautological-compare":        "Tautological compare",
	"tautology":                   "Tautology or contradiction",
	"write-after-write":           "Write after write",
	"boolean-cst":                 "Misuse of a Boolean constant",
	"constant-function-asm":       "Constant functions using assembly code",
	"constant-function-state":     "Constant functions changing the state",
	"divide-before-multiply":      "Divide before multiply",
	"out-of-order-retryable":      "Out-of-order retryable transactions",
	"reentrancy-no-eth":           "Reentrancy vulnerabilities",
	"reused-constructor":          "Reused base constructors",
	"tx-origin":                   "Dangerous usage of tx.origin",
	"unchecked-lowlevel":          "Unchecked low-level calls",
	"unchecked-send":              "Unchecked Send",
	"uninitialized-local":         "Uninitialized local variables",
	"unused-return":               "Unused return",
	"incorrect-modifier":          "Incorrect modifier",
	"shadowing-builtin":           "Builtin Symbol Shadowing",
	"shadowing-local":             "Local variable shadowing",
	"uninitialized-fptr-cst":      "Uninitialized function pointers in constructors",
	"variable-scope":              "Pre-declaration usage of local variables",
	"void-cst":                    "Void constructor",
	"calls-loop":                  "Calls inside a loop",
	"events-access":               "Missing events access control",
	"events-maths":                "Missing events arithmetic",
	"incorrect-unary":             "Dangerous unary expressions",
	"missing-zero-check":          "Missing zero address validation",
	"reentrancy-benign":           "Reentrancy vulnerabilities",
	"reentrancy-events":           "Reentrancy vulnerabilities",
	"return-bomb":                 "Return Bomb",
	"timestamp":                   "Block timestamp",
	"assembly":                    "Assembly usage",
	"assert-state-change":         "Assert state change",
	"boolean-equal":               "Boolean equality",
	"cyclomatic-complexity":       "Cyclomatic complexity",
	"deprecated-standards":        "Deprecated standards",
	"erc20-indexed":               "Unindexed ERC20 event parameters",
	"function-init-state":         "Function Initializing State",
	"incorrect-using-for":         "Incorrect usage of using-for statement",
	"low-level-calls":             "Low-level calls",
	"missing-inheritance":         "Missing inheritance",
	"naming-convention":           "Conformance to Solidity naming conventions",
	"pragma":                      "Different pragma directives are used",
	"redundant-statements":        "Redundant Statements",
	"solc-version":                "Incorrect versions of Solidity",
	"unimplemented-functions":     "Unimplemented functions",
	"unused-import":               "Unused Imports",
	"unused-state":                "Unused state variable",
	"costly-loop":                 "Costly operations inside a loop",
	"dead-code":                   "Dead-code",
	"reentrancy-unlimited-gas":    "Reentrancy vulnerabilities",
	"similar-names":               "Variable names too similar",
	"too-many-digits":             "Too many digits",
	"cache-array-length":          "Cache array length",
	"constable-states":            "State variables that could be declared constant",
	"external-function":           "Public function that could be declared external",
	"immutable-states":            "State variables that could be declared immutable",
	"var-read-using-this":         "Public variable read in external context",
}



func SlitherStandardize(sumup models.SumUp) (string, string) {
	name := sumup.Name
	vulnaDefine := ""
	switch name {
	case "Storage Signed Integer Array":
		vulnaDefine = OverallVulna[2]
		break
	case "Controlled Delegatecall":
	case "Payable functions using delegatecall inside a loop":
		vulnaDefine = OverallVulna[3]
		break
	case "Functions that send Ether to arbitrary destinations":
		vulnaDefine = OverallVulna[8]
		break
	case "Dangerous usage of tx.origin":
		vulnaDefine = OverallVulna[7]
		break
	case "Weak PRNG":
	case "Block timestamp":
		vulnaDefine = OverallVulna[5]
		break
	case "Reentrancy vulnerabilities":
		vulnaDefine = OverallVulna[1]
		break
	case "SLITHER ERROR":
		vulnaDefine = ""
		break
	default:
		vulnaDefine = name
		break
	}
	fmt.Println("Slither said: ", vulnaDefine, " with ", sumup.Severity)
	return vulnaDefine, sumup.Severity
}



func SolhintStandardize(sumup models.SumUp) (string, string) {
	name := sumup.Name
	severity := sumup.Severity
	if severity == "None"{
		return "",""
	}
	vulnaDefine := ""
	
	switch name {
	case "avoid-tx-origin":
		vulnaDefine = OverallVulna[7]
		break
	case "not-rely-on-block-hash":
	case "not-rely-on-time":
		vulnaDefine = OverallVulna[5]
		break
	case "reentrancy":
		vulnaDefine = OverallVulna[1]
		break
	case "SOLHINT ERROR":
		vulnaDefine = ""
		break
	default:
		vulnaDefine = sumup.Description
		break
	}

	// fmt.Println("Solhint said: ", vulnaDefine, " with ", severity)
	return vulnaDefine, severity
}

func HoneyBadgerStandardize(sumup models.SumUp) (string, string) {
	if sumup.Name == "HONEYBADGER ERROR"{
		return "",""
	}
	name := sumup.Name
	severity := sumup.Severity
	return name, severity
}
