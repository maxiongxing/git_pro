package BLC

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"math/big"
	"os"
)

const dbName  = "blockchain.db"
const blockTableName  = "blocks"

//创建区块链结构体
type Blockchain struct {
	Tip []byte    //存储最新区块的hash值
	DB  *bolt.DB    //数据库
	//Blocks []*Block    //存储有序的区块
}

//迭代器
func (blockchain *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{blockchain.Tip,blockchain.DB}
}
//判断是否存在
func dbExists() bool {
	if _,err := os.Stat(dbName); os.IsNotExist(err) {
		return false
	}
	return true
}


//遍历输出所有区块的信息

func (blc *Blockchain) PrintChain() {
	blockchainIterator := blc.Iterator()

	for  {
		block := blockchainIterator.Next()

		fmt.Printf("Height: %d\n",block.Height)
		fmt.Printf("Hash: %d\n",block.Hash)
		fmt.Printf("Data: %s\n",block.Data)
		//fmt.Printf("Timestamp: %s\n",time.Unix(block.Timestamp,0).Format("2006-01-02 03:04:25 PM"))
		fmt.Printf("Nonce: %d\n",block.Nonce)

		fmt.Println()

		//判断结束
		var hashInt big.Int
		hashInt.SetBytes(block.PrevBlockHash)

		if big.NewInt(0).Cmp(&hashInt) == 0{
			break;
		}
	}
}

//增加区块到区块链里面
func (blc *Blockchain) AddBlockToBlockchain(data string)  {

	err := blc.DB.Update(func(tx *bolt.Tx) error {
		//获取表
		b := tx.Bucket([]byte(blockTableName))

		if b != nil {
			//获取最新区块并反序列化
			blockBytes := b.Get(blc.Tip)
			block := DeserializeBlock(blockBytes)

			//创建新区块(高度加1
			newBlock := NewBlock(data, block.Height+1, block.Hash)

			//将区块序列化并存储到数据库
			err := b.Put(newBlock.Hash,newBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}
			//更新数据库里面的“l“对应的hash
			err = b.Put([]byte("l"),newBlock.Hash)
			if err != nil {
				log.Panic(err)
			}
			//更新blockchain中的tip
			blc.Tip = newBlock.Hash
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

//创建带有创世区块的区块链
func CreateBlockchainWithGenesisBlock() *Blockchain {

	//判断区块链是否存在，避免每次运行main都创建创世区块
	if dbExists() {
		fmt.Println("创世区块已经存在")

		db,err := bolt.Open(dbName,0600,nil)
		if err != nil {
			log.Fatal(err)
		}

		var blockchain *Blockchain
		err = db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(blockTableName))
			hash := b.Get([]byte("l"))
			blockchain = &Blockchain{hash,db}

			return nil
		})

		if err != nil {
			log.Panic(err)
		}
		return blockchain
	}

	db,err := bolt.Open(dbName,0600,nil)
	if err != nil {
		log.Fatal(err)
	}
	//defer db.Close()

	var blockHash []byte

	err = db.Update(func(tx *bolt.Tx) error {
		//获取表（创建
		b := tx.Bucket([]byte(blockTableName))
		if b == nil {
			b ,err = tx.CreateBucket([]byte(blockTableName))
			if err != nil {
				log.Panic(err)
			}
		}

		if b != nil {
			genesisBlock := CreateGenesisBlock("Genesis data")
			//将创世区块存储到库里
			err := b.Put(genesisBlock.Hash,genesisBlock.Serialize())
			if err != nil {
				log.Panic(err)
			}

			err = b.Put([]byte("l"),genesisBlock.Hash) // 存储最新区块的hash
			if err != nil {
				log.Panic(err)
			}

			blockHash = genesisBlock.Hash
		}
		return nil
	})

	//返回区块链对象
	return &Blockchain{blockHash,db}
}



