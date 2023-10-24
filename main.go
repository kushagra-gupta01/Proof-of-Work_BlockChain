package main
import (
	"net/http"
)

const difficulty = 1

type Block struct{
	Index int
	TimeStamp string
	Data int
	Hash string
	PrevHash string
	Difficulty int
	Nonce string
}

var BlockChain []Block

func main(){

}

func run() error{

}

func makeMuxRouter() http.Handler{
	HandlerFunc("/", handleGetBlockchain).Methods("GET")
	HandlerFunc("/", handleWriteBlock).Methods("POST")
}

func handleGetBlockchain(){

}

func handleWriteBlock(){

}

func responseWithJSON(){

}

func isBlockValid() bool{

}

func calculateHash() string{

}

func generateBlock(){ 

}

func isHashValid() bool{

}
