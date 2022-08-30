package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/33cn/chain33-sdk-go/client"
	evm "github.com/33cn/chain33-sdk-go/dapp/evm"
	crypto "github.com/33cn/chain33-sdk-go/crypto"
	"github.com/33cn/chain33-sdk-go/types"
)

var (
	// BTC类型地址
	// 合约部署人的地址和私钥
	deployAddress    = "14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
	deployPrivateKey = "CC38546E9E659D15E6B4893F0AB32A06D103931A8230B0BDE71459D2B27D6944"
	// 代扣手续费的地址和私钥
	withholdAddress    = "14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
	withholdPrivateKey = "CC38546E9E659D15E6B4893F0AB32A06D103931A8230B0BDE71459D2B27D6944"
	// 用户地址
	useraAddress    = "17RH6oiMbUjat3AAyQeifNiACPFefvz3Au"
	useraPrivateKey = "56d1272fcf806c3c5105f3536e39c8b33f88cb8971011dfe5886159201884763"

	url       = "http://172.22.16.179:8901"
	paraName  = "user.p.mbaas."
	// 普通地址类型（BTC形式）
	addressID = NormalAddressID
	chainID = 0

	// ETH类型地址
	// 合约部署人的地址和私钥
	//deployAddress    = "0x4cb94044427edb06ae7aeb8e8dd6eba078c8bc0a"
	//deployPrivateKey = "7dfe80684f7007b2829a28c85be681304f7f4cf6081303dbace925826e2891d1"
	//// 代扣手续费的地址和私钥
	//withholdAddress    = "0xfd89c32962f19bcea69b76093a64a03618cb33be"
	//withholdPrivateKey = "56d1272fcf806c3c5105f3536e39c8b33f88cb8971011dfe5886159201884763"
	//// 用户地址
	//useraAddress    = "0x6856f610b40e7321cace9e1f8752315110862573"
	//useraPrivateKey = "0x3967abcafaea83fee72766ca6dae578f4f156b5d1dae1ddf119e4564d5e2658c"
	//
	//url       = "http://122.224.77.188:8991"
	//paraName  = "user.p.para_pressuretest_2."
	//// 以太坊类型地址
	//addressID = ETHAddressID
	//chainID   = 999
)

func main() {
	jsonclient, err := client.NewJSONClient("", url)
	if err != nil {
		fmt.Println("The connection of jsonrpc is failed!")
	}
	// 部署合约
	code, err := types.FromHex(CODES_1155)

	tx, err := evm.CreateEvmContract(code, "", "evm-sdk-test", paraName, int32(addressID), int32(chainID))
	unsignTx := types.ToHex(types.Encode(tx))
	gas, err := evm.QueryEvmGas(url, unsignTx, deployAddress)

	fmt.Println("gas fee = ", gas)
	fee := evm.GetProperFee(url)
	fmt.Println("proper fee = ", fee)
	evm.UpdateTxFee(tx, gas, fee)
	err = crypto.SignTx(tx, deployPrivateKey, int32(addressID))

	signTx := types.ToHexPrefix(types.Encode(tx))
	txhash, err := jsonclient.SendTransaction(signTx)

	fmt.Print("部署合约交易hash = ", txhash)
	time.Sleep(8 * time.Second)
	detail, err := jsonclient.QueryTransaction(txhash)

	fmt.Println("; 部署合约交易执行结果（结果码=2代表成功） = ", detail.Receipt.Ty)

	// 计算合约地址
	contractAddress := evm.GetContractAddr(deployAddress, strings.TrimPrefix(txhash, "0x"), url)
	fmt.Println("部署好的合约地址 = " + contractAddress)

	length := 2
	// tokenId数组
	ids := make([]int, length)
	// 同一个tokenid发行多少份
	amounts := make([]int, length)
	// 每一个tokenid对应的URI信息（一般对应存放图片的描述信息，图片内容的一个url）
	uris := make([]string, length)
	for i := 0; i < length; i++ {
		ids[i] = 10000 + i
		amounts[i] = 100
		// 图片的属性，存储路径等描述信息
		uris[i] = "{\"图片描述\":\"由xxx创作\";\"创作时间\":\"2022/12/25\";\"图片存放路径\":\"http://www.baidu.com\"}";
	}
	idStr, _ := json.Marshal(ids)
	amountStr, _ := json.Marshal(amounts)
	uriStr, _ := json.Marshal(uris)

	// 调用合约
	param := fmt.Sprintf("mint(%s,%s,%s,%s)", useraAddress, idStr, amountStr, uriStr)
	initNFT, err := evm.EncodeParameter(ABI_1155, param)

	tx, err = evm.CallEvmContract(initNFT, "", 0, contractAddress, paraName, int32(addressID), int32(chainID))
	unsignTx = types.ToHex(types.Encode(tx))
	gas, err = evm.QueryEvmGas(url, unsignTx, deployAddress)

	fmt.Println("gas fee = ", gas)
	fee = evm.GetProperFee(url)
	fmt.Println("proper fee = ", fee)
	evm.UpdateTxFee(tx, gas, fee)
	// 构造交易组, deployPrivateKey:用于签名部署合约的交易， withholdPrivateKey：用于签名代扣交易
	group, err := evm.CreateNobalance(tx, deployPrivateKey, withholdPrivateKey, paraName, int32(addressID), int32(chainID))

	signTx = types.ToHexPrefix(types.Encode(group.Tx()))
	txhash, err = jsonclient.SendTransaction(signTx)

	fmt.Print("代扣交易hash = ", txhash)
	time.Sleep(8 * time.Second)
	detail, err = jsonclient.QueryTransaction(txhash)
	// 从交易组中取出EVM交易的hash值
	nextHash := detail.Tx.Next
	fmt.Println("; EVM交易hash = ", nextHash)

	detail, err = jsonclient.QueryTransaction(nextHash)
	fmt.Println("; EVM交易执行结果（结果码=2代表成功） = ", detail.Receipt.Ty)

	// 合约查询
	param = fmt.Sprintf("balanceOf(%s,%d)", useraAddress, ids[0])
	balance, err := evm.QueryContract(url, contractAddress, ABI_1155, param, contractAddress)

	fmt.Println(param, " = ", balance[0])
	param = fmt.Sprintf("uri(%d)", ids[1])
	uri, err := evm.QueryContract(url, contractAddress, ABI_1155, param, contractAddress)

	fmt.Println(param, " = ", uri[0])
}
