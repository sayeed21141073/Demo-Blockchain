package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Block struct {
	Timestamp     int64
	Transaction   Transaction
	Nonce         int
	Difficulty    int
	Hash          string
	PreviousBlock string
}

type Transaction struct {
	Sender   string
	Receiver string
	Amount   int64
}

var blockchain []Block

func main() {
	fmt.Println("Program started")
	router := mux.NewRouter()
	router.HandleFunc("/block/{id}", handleBlock).Methods("GET")
	router.HandleFunc("/new", handleNewBlock).Methods("POST")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./templates")))
	router.HandleFunc("/all", handleAll).Methods("GET")

	http.ListenAndServe(":8080", router)
}
func handleAll(w http.ResponseWriter, r *http.Request) {
	blockData, _ := json.Marshal(blockchain)
	w.Write(blockData)
}
func handleBlock(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	blockId, _ := strconv.Atoi(vars["id"])
	if blockId >= len(blockchain) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Block not found"))
		return
	}

	blockData, _ := json.Marshal(blockchain[blockId])
	w.Write(blockData)
}

func handleNewBlock(w http.ResponseWriter, r *http.Request) {
	var newTransaction Transaction
	json.NewDecoder(r.Body).Decode(&newTransaction)
	t := time.Now().Unix()
	difficulty := 2
	nonce := 0
	var previousHash string
	if len(blockchain) == 0 {
		previousHash = ""
	} else {
		previousHash = blockchain[len(blockchain)-1].Hash
	}

	newBlock := Block{t, newTransaction, nonce, difficulty, "", previousHash}

	for newBlock.Hash != calculateHash(newBlock.Timestamp, newBlock.Transaction, newBlock.Nonce, newBlock.Difficulty, previousHash) {
		nonce++
		difficulty++
		newBlock.Nonce = nonce
		newBlock.Hash = calculateHash(newBlock.Timestamp, newBlock.Transaction, newBlock.Nonce, newBlock.Difficulty, previousHash)
	}

	blockchain = append(blockchain, newBlock)
	blockData, _ := json.Marshal(newBlock)
	w.Write(blockData)
}

func calculateHash(t int64, transaction Transaction, nonce int, difficulty int, previousHash string) string {
	blockData := strconv.FormatInt(t, 10) + transaction.Sender + transaction.Receiver + strconv.FormatInt(transaction.Amount, 10) + strconv.Itoa(nonce) + strconv.Itoa(difficulty) + previousHash
	hash := sha256.Sum256([]byte(blockData))
	return hex.EncodeToString(hash[:])
}
