package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
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

var mutex sync.Mutex

func main(){
	err := godotenv.Load()
	if err != nil{
		log.Fatal(err)
	}
	
	go func(){
		t := time.Now()	
		genesisBlock := Block{}
		genesisBlock = Block{0,t.String(),0,calculateHash(genesisBlock),"",difficulty,""}
		spew.Dump(genesisBlock)
		mutex.Lock()
		BlockChain = append(BlockChain, genesisBlock)
		mutex.Unlock()
	}()
	log.Fatal(run())
}

func run() error{
	mux := makeMuxRouter()
	httpPort := os.Getenv("PORT")
	log.Println("HTTP server is running and lsitening on port:",httpPort)
	s := &http.Server{
		Addr: ":"+httpPort,
		Handler: mux,
		ReadTimeout: 10*time.Second,
		WriteTimeout: 10*time.Second,
		MaxHeaderBytes: 1>>20,
	}

	if err := s.ListenAndServe();err != nil{
		return err
	}
	return nil
}

func makeMuxRouter() http.Handler{
	muxRouter := mux.NewRouter()
	muxRouter.HandlerFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandlerFunc("/", handleWriteBlock).Methods("POST")
	return muxRouter
}

func handleGetBlockchain(w http.ResponseWriter,h *http.Request){
	bytes,err := json.MarshalIndent(BlockChain,""," ")
	if err !=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return
	}
	io.WriteString(w ,string(bytes))
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
