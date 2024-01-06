package main

import (
	"getContractDeployment/configs"
	"getContractDeployment/helper"
	"getContractDeployment/routes"
)

func main() {
    config, err := configs.LoadConfig(".")
    helper.CheckError(err)

    routes.Route(config)


    
    // contractPath := "C:\\DEV\\Backend\\getContractDeployment\\result"

    // result, err := docker.RunMythrilAnalysis(contractPath, "OssifiableProxy.sol", remappingJSON)


    // str := "// SPDX-License-Identifier: MIT\\n// OpenZeppelin Contracts (last updated v4.5.0) (token/ERC20/extensions/ERC20Burnable.sol)\\n\\npragma solidity ^0.8.0;\\n\\nimport \\\"./ERC20.sol\\\";\\nimport \\\"./Context.sol\\\";\\n\\n/**\\n * @dev Extension of {ERC20} that allows token holders to destroy both their own\\n * tokens and those that they have an allowance for, in a way that can be\\n * recognized off-chain (via event analysis).\\n */\\nabstract contract ERC20Burnable is Context, ERC20 {\\n    /**\\n     * @dev Destroys `amount` tokens from the caller.\\n     *\\n     * See {ERC20-_burn}.\\n     */\\n    function burn(uint256 amount) public virtual {\\n        _burn(_msgSender(), amount);\\n    }\\n\\n    /**\\n     * @dev Destroys `amount` tokens from `account`, deducting from the caller's\\n     * allowance.\\n     *\\n     * See {ERC20-_burn} and {ERC20-allowance}.\\n     *\\n     * Requirements:\\n     *\\n     * - the caller must have allowance for ``accounts``'s tokens of at least\\n     * `amount`.\\n     */\\n    function burnFrom(address account, uint256 amount) public virtual {\\n        _spendAllowance(account, _msgSender(), amount);\\n        _burn(account, amount);\\n    }\\n}\\n"
    // fmt.Println(str)
    // str, err := helper.UnescapeJSON(str)
    // helper.CheckError(err)
    // fmt.Println("--------------------------------------------------------------------------")
    // fmt.Println(helper.FindImport(str))


    // docker.TestDockerSDK()


    // config, err := configs.LoadConfig(".")
    // helper.CheckError(err)

    // // 0x889edC2eDab5f40e902b864aD4d7AdE8E412F9B1
    // // 0xb4039240E71535100BE947116c778d5D98bd9f62
    // contractAddr := controller.ContractAddress{
    //     Address: "0xb4039240E71535100BE947116c778d5D98bd9f62",
    //     ChainID: 1,
    // }

    // data := controller.GetContractSourceCode1(contractAddr, config)
    // dataReturn := controller.ContractCodeDataHandler(data, true)

    // err = docker.CreateMythrilMappingJson(dataReturn)
    // helper.CheckError(err)

    // fmt.Println("Analysis begin...")
    // contractPath := "C:\\DEV\\Backend\\getContractDeployment\\result"
    // analysisResult, err := docker.RunMythrilAnalysis(contractPath, dataReturn["Main Contract"].(string), true)
    // helper.CheckError(err)

    // fmt.Println(analysisResult)


    // // for contract, _ := range dataReturn{
    // //     fmt.Println(contract)
    // // }

    // resultPrint := make(map[string]map[string]string) 

    // for contract, content := range dataReturn {
    //     sourcePath := fmt.Sprintf("result/contracts/%s", contract)
    //     helper.WriteFile(content.(string), sourcePath)

    //     allImportPath := helper.FindImportPath(content.(string))
    //     importReplacement := make(map[string]string)
    //     for _, eachImportPath := range allImportPath{
    //         // if strings.HasPrefix(contract, "openzeppelin"){
    //         //     replacePath := helper.GetLastFilePath(eachImportPath)
    //         //     importReplacement[eachImportPath] = replacePath
                
    //         // }
    //         if strings.HasPrefix(eachImportPath, "@openzeppelin"){
    //             replacePath := "openzeppelin/" + helper.GetLastFilePath(eachImportPath)
    //             importReplacement[eachImportPath] = replacePath
    //         } else {
    //             replacePath := helper.GetLastFilePath(eachImportPath)
    //             importReplacement[eachImportPath] = replacePath
    //         }
    //     }
    //     resultPrint[contract] = importReplacement
    // }

    // for contract, importReplacement := range resultPrint{
    //     fmt.Println("=======================================================================================")
    //     fmt.Println("Contract: ", contract)
    //     for importPath, replacement := range importReplacement {
    //         fmt.Println(importPath, " = ", replacement)
    //     }
    //     fmt.Println("=======================================================================================")
    //     fmt.Println("")
    // }


    // fmt.Println(helper.GetPathToFile("@openzeppelin/contracts-v4.4/utils/StorageSlot.sol"))


//     url  := fmt.Sprintf("https://api.etherscan.io/api?module=%s&action=%s&address=%s&apikey=%s", "contract", "getsourcecode", "0xdAC17F958D2ee523a2206206994597C13D831ec7", config.ETHER_SCAN_API)
// // 0xdAC17F958D2ee523a2206206994597C13D831ec7
// // 0xca9b78435Be8267922E7Ac5cDE70401e7502c9cc
//     bodyChan := make(chan []byte)
//     var wg sync.WaitGroup
//     wg.Add(1)
//     go api.FetchAPIData(&wg, url ,bodyChan)
//     body := <- bodyChan
//     wg.Wait()

//     var data interface{}
//     err = json.Unmarshal(body, &data)
//     helper.CheckError(err)

//     resultsMap := data.(map[string]interface{})["result"]

//     resultInterface := resultsMap.([]interface{})[0]

//     result := resultInterface.(map[string]interface{})

//     sourceStr := result["SourceCode"].(string)

//     fmt.Println(sourceStr)

    // sourceStr = sourceStr[1:len(sourceStr)-1]

    // var sourceCode map[string]interface{}

    // err = json.Unmarshal([]byte(sourceStr), &sourceCode)
    // helper.CheckError(err)

    // source := sourceCode["sources"].(map[string]interface{})

    // for contract, contentInterface := range source {
    //     // if !strings.HasPrefix(contract, "@openzeppelin"){
    //     content := contentInterface.(map[string]interface{})["content"].(string)
    //     // fmt.Println("Contract name: ", contract)
    //     // fmt.Println("Source code: ", content)
    //     // fmt.Println("")
        
    //     if !strings.HasPrefix(contract, "@openzeppelin"){
    //         sourcePath := fmt.Sprintf("result/%s", contract)
    //         helper.WriteFile(content, sourcePath)
    //     }else {
    //         sourcePath := fmt.Sprintf("result/contracts/%s", path.Base(contract))
    //         helper.WriteFile(content, sourcePath)
    //     }
    // }

}