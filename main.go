package main

import (
	"getContractDeployment/configs"
	"getContractDeployment/helper"
	"getContractDeployment/routes"
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	config, err := configs.LoadConfig(".")
	helper.CheckError(err)

	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	routes.Route(config)

	// data := `[{
	// 	"tool_name": "mythril",
	// 	"no_error": 2,
	// 	"sum_up": [
	// 		{
	// 			"name": "Reentrancy",
	// 			"description": "A call to a user-supplied address is executed.\nAn external message call to an address specified by the caller is executed. Note that the callee account might contain arbitrary code and could re-enter any function within this contract. Reentering the contract in an intermediate state may lead to unexpected behaviour. Make sure that no state modifications are executed after this call and/or reentrancy guards are in place.",
	// 			"severity": "Low"
	// 		},
	// 		{
	// 			"name": "Reentrancy",
	// 			"description": "Write to persistent state following external call\nThe contract account state is accessed after an external call to a user defined address. To prevent reentrancy issues, consider accessing the state only before the call, especially if the callee is untrusted. Alternatively, a reentrancy lock can be used to prevent untrusted callees from re-entering the contract in an intermediate state.",
	// 			"severity": "Medium"
	// 		}
	// 	],
	// 	"detail": {
	// 		"error": null,
	// 		"issues": [
	// 			{
	// 				"address": 351,
	// 				"code": "msg.sender.call{value: bal}(\"\")",
	// 				"contract": "EtherStore",
	// 				"description": "A call to a user-supplied address is executed.\nAn external message call to an address specified by the caller is executed. Note that the callee account might contain arbitrary code and could re-enter any function within this contract. Reentering the contract in an intermediate state may lead to unexpected behaviour. Make sure that no state modifications are executed after this call and/or reentrancy guards are in place.",
	// 				"filename": "/mnt/result/contracts628/example.sol",
	// 				"function": "withdraw()",
	// 				"lineno": 15,
	// 				"max_gas_used": 62205,
	// 				"min_gas_used": 7261,
	// 				"severity": "Low",
	// 				"sourceMap": ":::-",
	// 				"swc-id": "107",
	// 				"title": "External Call To User-Supplied Address",
	// 				"txsequence": null
	// 			},
	// 			{
	// 				"address": 534,
	// 				"code": "balances[msg.sender] = 0",
	// 				"contract": "EtherStore",
	// 				"description": "Write to persistent state following external call\nThe contract account state is accessed after an external call to a user defined address. To prevent reentrancy issues, consider accessing the state only before the call, especially if the callee is untrusted. Alternatively, a reentrancy lock can be used to prevent untrusted callees from re-entering the contract in an intermediate state.",
	// 				"filename": "/mnt/result/contracts628/example.sol",
	// 				"function": "withdraw()",
	// 				"lineno": 18,
	// 				"max_gas_used": 62205,
	// 				"min_gas_used": 7261,
	// 				"severity": "Medium",
	// 				"sourceMap": ":24",
	// 				"swc-id": "107",
	// 				"title": "State access after external call",
	// 				"txsequence": null
	// 			}
	// 		],
	// 		"success": true
	// 	},
	// 	"time_elapsed": 55.8879681
	// },
	// {
	// 	"tool_name": "slither",
	// 	"no_error": 1,
	// 	"sum_up": [
	// 		{
	// 			"name": "Reentrancy vulnerabilities",
	// 			"description": "Reentrancy in EtherStore.withdraw() (example.sol#11-19):\n\tExternal calls:\n\t- (sent) = msg.sender.call{value: bal}() (example.sol#15)\n\tState variables written after the call(s):\n\t- balances[msg.sender] = 0 (example.sol#18)\n\tEtherStore.balances (example.sol#5) can be used in cross function reentrancies:\n\t- EtherStore.balances (example.sol#5)\n\t- EtherStore.deposit() (example.sol#7-9)\n\t- EtherStore.withdraw() (example.sol#11-19)\n",
	// 			"severity": "High"
	// 		}
	// 	],
	// 	"detail": {
	// 		"success": true,
	// 		"error": null,
	// 		"results": {
	// 			"detectors": [
	// 				{
	// 					"elements": [
	// 						{
	// 							"type": "function",
	// 							"name": "withdraw",
	// 							"source_mapping": {
	// 								"ending_column": 6,
	// 								"filename_absolute": "/share/result/contracts628/example.sol",
	// 								"filename_relative": "example.sol",
	// 								"filename_short": "example.sol",
	// 								"is_dependency": false,
	// 								"length": 243,
	// 								"lines": [
	// 									11,
	// 									12,
	// 									13,
	// 									14,
	// 									15,
	// 									16,
	// 									17,
	// 									18,
	// 									19
	// 								],
	// 								"start": 224,
	// 								"starting_column": 5
	// 							},
	// 							"type_specific_fields": {
	// 								"parent": {
	// 									"name": "EtherStore",
	// 									"source_mapping": {
	// 										"ending_column": 2,
	// 										"filename_absolute": "/share/result/contracts628/example.sol",
	// 										"filename_relative": "example.sol",
	// 										"filename_short": "example.sol",
	// 										"is_dependency": false,
	// 										"length": 575,
	// 										"lines": [
	// 											4,
	// 											5,
	// 											6,
	// 											7,
	// 											8,
	// 											9,
	// 											10,
	// 											11,
	// 											12,
	// 											13,
	// 											14,
	// 											15,
	// 											16,
	// 											17,
	// 											18,
	// 											19,
	// 											20,
	// 											21,
	// 											22,
	// 											23,
	// 											24,
	// 											25
	// 										],
	// 										"start": 58,
	// 										"starting_column": 1
	// 									},
	// 									"type": "contract"
	// 								},
	// 								"signature": "withdraw()"
	// 							}
	// 						},
	// 						{
	// 							"type": "node",
	// 							"name": "(sent) = msg.sender.call{value: bal}()",
	// 							"source_mapping": {
	// 								"ending_column": 55,
	// 								"filename_absolute": "/share/result/contracts628/example.sol",
	// 								"filename_relative": "example.sol",
	// 								"filename_short": "example.sol",
	// 								"is_dependency": false,
	// 								"length": 46,
	// 								"lines": [
	// 									15
	// 								],
	// 								"start": 332,
	// 								"starting_column": 9
	// 							},
	// 							"type_specific_fields": {
	// 								"parent": {
	// 									"name": "withdraw",
	// 									"source_mapping": {
	// 										"ending_column": 6,
	// 										"filename_absolute": "/share/result/contracts628/example.sol",
	// 										"filename_relative": "example.sol",
	// 										"filename_short": "example.sol",
	// 										"is_dependency": false,
	// 										"length": 243,
	// 										"lines": [
	// 											11,
	// 											12,
	// 											13,
	// 											14,
	// 											15,
	// 											16,
	// 											17,
	// 											18,
	// 											19
	// 										],
	// 										"start": 224,
	// 										"starting_column": 5
	// 									},
	// 									"type": "function",
	// 									"type_specific_fields": {
	// 										"parent": {
	// 											"name": "EtherStore",
	// 											"source_mapping": {
	// 												"ending_column": 2,
	// 												"filename_absolute": "/share/result/contracts628/example.sol",
	// 												"filename_relative": "example.sol",
	// 												"filename_short": "example.sol",
	// 												"is_dependency": false,
	// 												"length": 575,
	// 												"lines": [
	// 													4,
	// 													5,
	// 													6,
	// 													7,
	// 													8,
	// 													9,
	// 													10,
	// 													11,
	// 													12,
	// 													13,
	// 													14,
	// 													15,
	// 													16,
	// 													17,
	// 													18,
	// 													19,
	// 													20,
	// 													21,
	// 													22,
	// 													23,
	// 													24,
	// 													25
	// 												],
	// 												"start": 58,
	// 												"starting_column": 1
	// 											},
	// 											"type": "contract"
	// 										},
	// 										"signature": "withdraw()"
	// 									}
	// 								}
	// 							}
	// 						},
	// 						{
	// 							"type": "node",
	// 							"name": "balances[msg.sender] = 0",
	// 							"source_mapping": {
	// 								"ending_column": 33,
	// 								"filename_absolute": "/share/result/contracts628/example.sol",
	// 								"filename_relative": "example.sol",
	// 								"filename_short": "example.sol",
	// 								"is_dependency": false,
	// 								"length": 24,
	// 								"lines": [
	// 									18
	// 								],
	// 								"start": 436,
	// 								"starting_column": 9
	// 							},
	// 							"type_specific_fields": {
	// 								"parent": {
	// 									"name": "withdraw",
	// 									"source_mapping": {
	// 										"ending_column": 6,
	// 										"filename_absolute": "/share/result/contracts628/example.sol",
	// 										"filename_relative": "example.sol",
	// 										"filename_short": "example.sol",
	// 										"is_dependency": false,
	// 										"length": 243,
	// 										"lines": [
	// 											11,
	// 											12,
	// 											13,
	// 											14,
	// 											15,
	// 											16,
	// 											17,
	// 											18,
	// 											19
	// 										],
	// 										"start": 224,
	// 										"starting_column": 5
	// 									},
	// 									"type": "function",
	// 									"type_specific_fields": {
	// 										"parent": {
	// 											"name": "EtherStore",
	// 											"source_mapping": {
	// 												"ending_column": 2,
	// 												"filename_absolute": "/share/result/contracts628/example.sol",
	// 												"filename_relative": "example.sol",
	// 												"filename_short": "example.sol",
	// 												"is_dependency": false,
	// 												"length": 575,
	// 												"lines": [
	// 													4,
	// 													5,
	// 													6,
	// 													7,
	// 													8,
	// 													9,
	// 													10,
	// 													11,
	// 													12,
	// 													13,
	// 													14,
	// 													15,
	// 													16,
	// 													17,
	// 													18,
	// 													19,
	// 													20,
	// 													21,
	// 													22,
	// 													23,
	// 													24,
	// 													25
	// 												],
	// 												"start": 58,
	// 												"starting_column": 1
	// 											},
	// 											"type": "contract"
	// 										},
	// 										"signature": "withdraw()"
	// 									}
	// 								}
	// 							}
	// 						}
	// 					],
	// 					"description": "Reentrancy in EtherStore.withdraw() (example.sol#11-19):\n\tExternal calls:\n\t- (sent) = msg.sender.call{value: bal}() (example.sol#15)\n\tState variables written after the call(s):\n\t- balances[msg.sender] = 0 (example.sol#18)\n\tEtherStore.balances (example.sol#5) can be used in cross function reentrancies:\n\t- EtherStore.balances (example.sol#5)\n\t- EtherStore.deposit() (example.sol#7-9)\n\t- EtherStore.withdraw() (example.sol#11-19)\n",
	// 					"markdown": "Reentrancy in [EtherStore.withdraw()](example.sol#L11-L19):\n\tExternal calls:\n\t- [(sent) = msg.sender.call{value: bal}()](example.sol#L15)\n\tState variables written after the call(s):\n\t- [balances[msg.sender] = 0](example.sol#L18)\n\t[EtherStore.balances](example.sol#L5) can be used in cross function reentrancies:\n\t- [EtherStore.balances](example.sol#L5)\n\t- [EtherStore.deposit()](example.sol#L7-L9)\n\t- [EtherStore.withdraw()](example.sol#L11-L19)\n",
	// 					"first_markdown_element": "example.sol#L11-L19",
	// 					"id": "6cdea23aac6058d23f99876c7d927dbd5e6af229d1a1fa83c64007c70ec6b892",
	// 					"check": "reentrancy-eth",
	// 					"impact": "High",
	// 					"confidence": "Medium"
	// 				},
	// 				{
	// 					"elements": [
	// 						{
	// 							"type": "pragma",
	// 							"name": "^0.8.22",
	// 							"source_mapping": {
	// 								"ending_column": 25,
	// 								"filename_absolute": "/share/result/contracts628/example.sol",
	// 								"filename_relative": "example.sol",
	// 								"filename_short": "example.sol",
	// 								"is_dependency": false,
	// 								"length": 24,
	// 								"lines": [
	// 									2
	// 								],
	// 								"start": 32,
	// 								"starting_column": 1
	// 							},
	// 							"type_specific_fields": {
	// 								"directive": [
	// 									"solidity",
	// 									"^",
	// 									"0.8",
	// 									".22"
	// 								]
	// 							}
	// 						}
	// 					],
	// 					"description": "Pragma version^0.8.22 (example.sol#2) necessitates a version too recent to be trusted. Consider deploying with 0.8.18.\n",
	// 					"markdown": "Pragma version[^0.8.22](example.sol#L2) necessitates a version too recent to be trusted. Consider deploying with 0.8.18.\n",
	// 					"first_markdown_element": "example.sol#L2",
	// 					"id": "106b756b77ac15951aea326f44d3b136eef6ae66174c83cd5f92cec57f8ed7e2",
	// 					"check": "solc-version",
	// 					"impact": "Informational",
	// 					"confidence": "High"
	// 				},
	// 				{
	// 					"elements": [],
	// 					"description": "solc-0.8.22 is not recommended for deployment\n",
	// 					"markdown": "solc-0.8.22 is not recommended for deployment\n",
	// 					"first_markdown_element": "",
	// 					"id": "2584eed3f1a6c37118da98709458a611a90fc5486a349e6fad8b5af8d201ac98",
	// 					"check": "solc-version",
	// 					"impact": "Informational",
	// 					"confidence": "High"
	// 				},
	// 				{
	// 					"elements": [
	// 						{
	// 							"type": "function",
	// 							"name": "withdraw",
	// 							"source_mapping": {
	// 								"ending_column": 6,
	// 								"filename_absolute": "/share/result/contracts628/example.sol",
	// 								"filename_relative": "example.sol",
	// 								"filename_short": "example.sol",
	// 								"is_dependency": false,
	// 								"length": 243,
	// 								"lines": [
	// 									11,
	// 									12,
	// 									13,
	// 									14,
	// 									15,
	// 									16,
	// 									17,
	// 									18,
	// 									19
	// 								],
	// 								"start": 224,
	// 								"starting_column": 5
	// 							},
	// 							"type_specific_fields": {
	// 								"parent": {
	// 									"name": "EtherStore",
	// 									"source_mapping": {
	// 										"ending_column": 2,
	// 										"filename_absolute": "/share/result/contracts628/example.sol",
	// 										"filename_relative": "example.sol",
	// 										"filename_short": "example.sol",
	// 										"is_dependency": false,
	// 										"length": 575,
	// 										"lines": [
	// 											4,
	// 											5,
	// 											6,
	// 											7,
	// 											8,
	// 											9,
	// 											10,
	// 											11,
	// 											12,
	// 											13,
	// 											14,
	// 											15,
	// 											16,
	// 											17,
	// 											18,
	// 											19,
	// 											20,
	// 											21,
	// 											22,
	// 											23,
	// 											24,
	// 											25
	// 										],
	// 										"start": 58,
	// 										"starting_column": 1
	// 									},
	// 									"type": "contract"
	// 								},
	// 								"signature": "withdraw()"
	// 							}
	// 						},
	// 						{
	// 							"type": "node",
	// 							"name": "(sent) = msg.sender.call{value: bal}()",
	// 							"source_mapping": {
	// 								"ending_column": 55,
	// 								"filename_absolute": "/share/result/contracts628/example.sol",
	// 								"filename_relative": "example.sol",
	// 								"filename_short": "example.sol",
	// 								"is_dependency": false,
	// 								"length": 46,
	// 								"lines": [
	// 									15
	// 								],
	// 								"start": 332,
	// 								"starting_column": 9
	// 							},
	// 							"type_specific_fields": {
	// 								"parent": {
	// 									"name": "withdraw",
	// 									"source_mapping": {
	// 										"ending_column": 6,
	// 										"filename_absolute": "/share/result/contracts628/example.sol",
	// 										"filename_relative": "example.sol",
	// 										"filename_short": "example.sol",
	// 										"is_dependency": false,
	// 										"length": 243,
	// 										"lines": [
	// 											11,
	// 											12,
	// 											13,
	// 											14,
	// 											15,
	// 											16,
	// 											17,
	// 											18,
	// 											19
	// 										],
	// 										"start": 224,
	// 										"starting_column": 5
	// 									},
	// 									"type": "function",
	// 									"type_specific_fields": {
	// 										"parent": {
	// 											"name": "EtherStore",
	// 											"source_mapping": {
	// 												"ending_column": 2,
	// 												"filename_absolute": "/share/result/contracts628/example.sol",
	// 												"filename_relative": "example.sol",
	// 												"filename_short": "example.sol",
	// 												"is_dependency": false,
	// 												"length": 575,
	// 												"lines": [
	// 													4,
	// 													5,
	// 													6,
	// 													7,
	// 													8,
	// 													9,
	// 													10,
	// 													11,
	// 													12,
	// 													13,
	// 													14,
	// 													15,
	// 													16,
	// 													17,
	// 													18,
	// 													19,
	// 													20,
	// 													21,
	// 													22,
	// 													23,
	// 													24,
	// 													25
	// 												],
	// 												"start": 58,
	// 												"starting_column": 1
	// 											},
	// 											"type": "contract"
	// 										},
	// 										"signature": "withdraw()"
	// 									}
	// 								}
	// 							}
	// 						}
	// 					],
	// 					"description": "Low level call in EtherStore.withdraw() (example.sol#11-19):\n\t- (sent) = msg.sender.call{value: bal}() (example.sol#15)\n",
	// 					"markdown": "Low level call in [EtherStore.withdraw()](example.sol#L11-L19):\n\t- [(sent) = msg.sender.call{value: bal}()](example.sol#L15)\n",
	// 					"first_markdown_element": "example.sol#L11-L19",
	// 					"id": "5f7f5945205c9c9ab7f2f02e104cf1465660b5832861cb5b528f09dcfa540d49",
	// 					"check": "low-level-calls",
	// 					"impact": "Informational",
	// 					"confidence": "High"
	// 				}
	// 			]
	// 		}
	// 	},
	// 	"time_elapsed": 3.8052021
	// }]`

	
	// var toolResult []models.ToolResult
	// err := json.Unmarshal([]byte(data), &toolResult)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(toolResult[0].ToolName)
	// standardizeResult := controller.StandardizeResult(toolResult)
	// fmt.Println(standardizeResult)

	// file := "typeoverflow_honey.sol"
	// contractFolder := "contracts"
	// remappingJSON := false
	// var toolsResult []models.ToolResult

	// honeybadgerStart := time.Now()
	// honeybadgerDetail, err := docker.RunHoneyBadgerAnalysisWithTimeOut(file, contractFolder, remappingJSON)
	// if err != nil{
	// 	helper.WriteFileExtra(err.Error(), "log.txt")
	// 	// return models.Result{}, err
	// }
	// // fmt.Print(detail)
	// honeybadgerSumUp := docker.GetHoneyBadgerSumUp(honeybadgerDetail)
	// honeybadgerEnd := time.Since(honeybadgerStart)

	// var honeybadger models.ToolResult
	// honeybadger.ToolName = "honeybadger"
	// honeybadger.SumUps = honeybadgerSumUp
	// honeybadger.NoError = len(honeybadgerSumUp)
	// honeybadger.Detail = honeybadgerDetail
	// honeybadger.TimeElapsed = honeybadgerEnd.Seconds()
	// toolsResult = append(toolsResult, honeybadger)
	

	// standardize := controller.StandardizeResult(toolsResult)
	// fmt.Print(standardize)

	// docker.Test(file, contractFolder, remappingJSON)

}