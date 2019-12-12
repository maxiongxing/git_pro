package main

import (
	"../part24_persistence_cli/BLC"
)

//创建创世区块并存到数据库

func main()  {
	//创世区块
	blockchain := BLC.CreateBlockchainWithGenesisBlock()
	//defer blockchain.DB.Close()

	cli := BLC.CLI{blockchain}
	cli.Run()
}


