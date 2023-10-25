package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
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

type Message struct{
	data int
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

func handleWriteBlock(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","application/json")
	var m Message
	err := json.NewDecoder(r.Body).Decode(&m);if err !=nil{
		responseWithJSON(w,r,http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	mutex.Lock()
	newBlock := generateBlock(BlockChain[len(BlockChain)-1],m.data)
	mutex.Unlock()

	if isBlockValid(newBlock,BlockChain[len(BlockChain)-1]){
		BlockChain = append(BlockChain, newBlock)
		spew.Dump(BlockChain)
	}

	responseWithJSON(w,r,http.StatusCreated,newBlock)
}

func responseWithJSON(w http.ResponseWriter,r *http.Request, code int, payload interface{}){
	w.Header().Set("Content-Type","application/json")
	response,err :=json.MarshalIndent(payload,""," ")
	if err !=nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500:Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}

func isBlockValid(newBlock,oldBlock Block) bool{
	if oldBlock.Index +1 != newBlock.Index{
		return false
	}
	
	if oldBlock.Hash != newBlock.PrevHash{
		return false
	}

	if calculateHash(newBlock) !=newBlock.Hash {
		return false
	}
	return true
}

func calculateHash(block Block) string{
	record := strconv.Itoa(block.Index) + block.Nonce + strconv.Itoa(block.Data) + block.TimeStamp + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed) 
}

func generateBlock(oldBlock Block,Data int)(Block){ 
	var newBlock Block
	t := time.Now()

	newBlock.Index  = oldBlock.Index +1
	newBlock.Data = Data
	newBlock.TimeStamp = t.String()
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Difficulty = difficulty
	
	for i := 0 ; ; i++ {
		hex := fmt.Sprintf("%x", i)
		newBlock.Nonce = hex
		if !isHashValid(calculateHash(newBlock), newBlock.Difficulty){
			fmt.Println(calculateHash(newBlock), "do more work")
			time.Sleep(time.Second)
			continue
		}
		else{
			fmt.Println(calculateHash(newBlock), "do more work")
			newBlock.hash = calculateHash(newBlock)
			break
		}
	}
	return newBlock
}

func isHashValid(hash string,difficulty int) bool{
	prefix := strings.Repeat("0",1)
	return strings.HasPrefix(hash,prefix)
}
