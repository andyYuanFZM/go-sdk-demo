package main

import (
	crypto "github.com/33cn/chain33-sdk-go/crypto"
	c "github.com/33cn/chain33-sdk-go/client"
	"testing"
	"fmt"
)

// 生成BTC类型地址
func TestCreateAccountNormal(t *testing.T) {
	acc1, _ := crypto.NewAccount(NormalAddressID)
	fmt.Println("私钥", acc1.PrivateKey)
	fmt.Println("地址", acc1.Address)
}

// 生成以太坊类型地址
func TestCreateAccountEth(t *testing.T) {
	acc1, _ := crypto.NewAccount(ETHAddressID)
	fmt.Println("私钥", acc1.PrivateKey)
	fmt.Println("地址", acc1.Address)
}

// 获取当前最大区块高度
func TestGetLastHeader(t *testing.T) {
	url := "http://122.224.77.188:8991"
	client, _ := c.NewJSONClient("", url)
	header, _ := client.GetLastHeader()
	fmt.Println("当前最大区块高度为:", header.Height)
}

// 根据区块高度获取区块信息
func TestGetBlocks(t *testing.T) {
	url := "http://122.224.77.188:8991"
	client, _ := c.NewJSONClient("", url)
	blocks, _ := client.GetBlocks(10, 12, true)
	for i := 0; i < len(blocks.Items); i++ {
		block := blocks.Items[i].Block
		fmt.Println("前一个区块hash:", block.ParentHash)
		fmt.Println("默克尔根hash:", block.TxHash)
		fmt.Println("区块高度:", block.Height)
		fmt.Println("区块时间戳:", block.BlockTime)
		fmt.Println("区块中交易数:", len(block.Txs))
		for j := 0; j < len(block.Txs); j++ {
			fmt.Printf("第%d笔交易hash值:%v\n", j+1, block.Txs[j].Hash)
		}
	}
}

