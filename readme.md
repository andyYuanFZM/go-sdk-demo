# 结合chain33-go-sdk部署调用EVM合约

## 文档修改记
| 版本号 | 版本描述                              | 修改日期   | 备注 |
| ------ | ------------------------------------- | ---------- | ---- |
| V1.0   | 通过 GO-SDK 在平行链上发行 NFT| 2022/07/18 |

## 通过GO-SDK实现合约部署调用
### 1. 子目录说明
demo目录  
- commont_test.go：包含一些常用的方法（区块链地址和私钥的生成; 获取最大区块高度; 根据区块高度获取区块详情信息）  
- main.go: 通过go-sdk部署，运行和查询evm合约（ERC1155）的方法  

solidity目录  
- ERC1155ByManager.sol: ERC1155发行NFT,转移NFT的合约样例。  
- ERC721ByManager.sol: ERC721发行NFT,转移NFT的合约样例。  

### 2. 运行Demo程序
2.1 生成公私钥（YCC链用ETH类型地址，其余的(BTY或联盟链)用BTC类型地址）  
生成BTC类型地址：go test -v ./ -test.run TestCreateAccountNormal  
生成ETH类型地址：go test -v ./ -test.run TestCreateAccountEth  

2.2 运行部署，发行，调用合约程序  
- 修改main.go，将上一步生成的内容，分别填充到以下几个参数中，注意私钥即资产，要隐私存放，而地址是可以公开的  
```  
// 合约部署人的地址和私钥
deployAddress    = "14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
deployPrivateKey = "CC38546E9E659D15E6B4893F0AB32A06D103931A8230B0BDE71459D2B27D6944"
// 代扣手续费的地址和私钥
withholdAddress    = "14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
withholdPrivateKey = "CC38546E9E659D15E6B4893F0AB32A06D103931A8230B0BDE71459D2B27D6944"
// 用户地址
useraAddress    = "17RH6oiMbUjat3AAyQeifNiACPFefvz3Au"
useraPrivateKey = "56d1272fcf806c3c5105f3536e39c8b33f88cb8971011dfe5886159201884763"
```

- 给上述合约部署人和代扣地址下充燃料（BTY和YCC需要燃料，联盟链不需要）, 所有的用户都走代扣地址来扣除燃料  
- 修改main.go以下两个参数  
```  
// 改成自己平行链所在服务器IP地址
url  = "http://172.22.16.179:8901"
// 改成和实际部署的平行链配置文件中的title值一致
paraName  = "user.p.mbaas."
// 地址类型（区分是BTC类型地址还是ETH类型地址）
addressID = NormalAddressID
```
- 运行测试程序  
```
windows下： go run ./
linux下： go run *.go
```
