package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	//"encoding/json"
	//"fmt"
	"io"
	"log"
	// "math/rand"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

// Block represents each 'item' in the blockchain
type Block struct {
	Index     int
	Timestamp string
	Testimony string
	Hash      string
	PrevHash  string
	Validator string
}

// Blockchain is a series of validated Blocks
var Blockchain []Block
var tempBlocks []Block

// candidateBlocks handles incoming blocks for validation
var candidateBlocks = make(chan Block)

// announcements broadcasts winning validator to all nodes
var announcements = make(chan string)
var icomoney uint64 = 1000000000

var mutex = &sync.Mutex{}

// validators keeps track of open validators and balances
var validators = make(map[string]float32)

var validators_percent = make(map[string]float32)

var reward_percent = make(map[string]float32)

var reward_map = make(map[string]float32)
var reward_final = make(map[string]float32)

var transactionReward int = 100

var delegateCounter int = 1
var delegates [5]string
var allDelegates [51]string
var delegateBalance = make(map[string]float32)
var data [10] string
var count int

var testimonyMap = make(map[string]string)
var completeData string = ""
var wv string = ""
var delinit int=0
var countdelegate int = 0

var votes[10] int
var p int=0
var delvotes int=0
var currentblockdata string=""

//var length_block int=1


var address string

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}


	generateAllDelegates()
	chooseDelegates()
	//We are initializing the stake of all the delegates in the crypto-system to 1000
	for i := 0; i<=4; i++ {

		delegateBalance[delegates[i]] = 1000
	}

	// create genesis block
	
	log.Println("Creating genesis block\n")
	time.Sleep(2*time.Second)
	t := time.Now()
	genesisBlock := Block{}
	genesisBlock = Block{0, t.String(), "0", calculateBlockHash(genesisBlock), "", ""}
	log.Println("Genesis Block Created!!")
	time.Sleep(2*time.Second)
	spew.Dump(genesisBlock)
	time.Sleep(2*time.Second)
	Blockchain = append(Blockchain, genesisBlock)
	log.Println("\t\t\t\t\t\t\n****************************************************************************\n")
	log.Println("Listening to ports")

	httpPort := os.Getenv("PORT")

	// start TCP and serve TCP server
	server, err := net.Listen("tcp", ":"+httpPort)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("HTTP Server Listening on port :", httpPort)
	// log.Println("This is it ",validators)

	defer server.Close()

	go func() {
		for {
			pickWinner()
		}
	}()
	// log.Println("After pickWinner:",validators)
	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn)
	}
}

// pickWinner creates a lottery pool of validators and chooses the validator who gets to forge a block to the blockchain
// by random selecting from the pool, weighted by amount of tokens staked
func pickWinner() {
	//log.Println("Inside Pick winner")
	var totalstake float32 =0.00
	time.Sleep(20 * time.Second)


	// lotteryPool := []string{}
	if count > 0 {

		// Getting highest voted delegate
		var max int = 0
		for i := range votes {
			if votes[i]>max {
				max = votes[i]
				wv=delegates[i]

			}

		}
		log.Println("\n****************************************************************************")
		log.Println("Chosen validator is: ", wv)
		time.Sleep(2*time.Second)
				mutex.Lock()
				oldLastIndex := Blockchain[len(Blockchain)-1]
				mutex.Unlock()

				// create newBlock for consideration to be forged
				newBlock, err := generateBlock(oldLastIndex, completeData, address)
				//log.Println("block generated")
				//log.Println("iNCREASING LENGTH", length_block)
				//length_block=length_block + 1
				//log.Println("Length increased", length_block)
				if err != nil {
					log.Println(err)
					
				}
				validatebyconsensus(newBlock)
				Blockchain = append(Blockchain, newBlock)
				if !isBlockValid(newBlock, oldLastIndex) {
					log.Println("block invalid not added")
				}
				log.Println("The block has been appended in the chain.\n")
				time.Sleep(2*time.Second)
				log.Println("\n****************************************************************************")

		for _, bal := range validators {
			totalstake = totalstake + bal
		}

		var limit_totalstake float32
		limit_totalstake=(0.3)*float32(totalstake)
		var percent_stake float32
		//reward_map := validators
		for address,bal := range validators{
			percent_stake = (float32(bal)/float32(totalstake))*100
			reward_map[address] = bal
			if(percent_stake > 30.00){
				reward_map[address]=limit_totalstake

			}
			// validators_percent[address] = percent_stake  
		}
		//log.Println("The reward map is ",reward_map)
		//log.Println("Total Stake is :",totalstake)
		var totalstake2 float32=0.00
		for _, bal1 := range reward_map {
			totalstake2 = totalstake2 + bal1
		}
		//log.Println("Total final Stake is :",totalstake2)
		var result float32
		for address1,bal3 := range reward_map{
			//log.Println("The bal3 is......... ",bal3)
			//log.Println("validators",validators)
			//log.Println(validators[address1])
			validators[address1] = validators[address1] + ((float32(bal3)/float32(totalstake2))*75.00)
			result = validators[address1] + ((float32(bal3)/float32(totalstake2))*75.00)
			//log.Println("The result is: ",result)
			reward_final[address1] = result
		}

		delegateBalance[wv] = delegateBalance[wv] + 25

		log.Println(" The new token balances of participants: \n", validators)
		time.Sleep(2*time.Second)
		//log.Println(" \n\n\n\n The new token balances of the delegates: \n ", delegateBalance)
		countdelegate=countdelegate+1
		if(countdelegate%3 == 0){
			chooseDelegates()
		}
		//log.Println("Delegates are  ",delegates)
		//Reward time,Transaction fees is 100
		// reward_map := validators
		// log.Println("The percentage map is ",validators_percent)
		// for address,bal := range reward_map{
		// 	if ( validators_percent[address] > 30.00 ){
		// 		bal = int(limit_totalstake)
		// 		// reward_map[address]=bal
		// 		log.Println(bal)
		// 	}
			  
		// }

		//log.Println("The reward final is ",reward_final)
		
	}
	
	

	// mutex.Lock()
	// tempBlocks = []Block{}
	// mutex.Unlock()
}

func handleConn(conn net.Conn) {

	var sd int = 0
	err := godotenv.Load()
	defer conn.Close()

	go func() {
		for {
			msg := <-announcements
			io.WriteString(conn, msg)
		}
	}()
	log.Println("A new participant is connected.")
	// validator address
	// var address string

	// allow user to allocate number of tokens to stake
	// the greater the number of tokens, the greater chance to forging a new block
	//var sd int
	var d int
	d=count
	//count=count+1
	io.WriteString(conn, "\n\nWelcome to the implementation of Proof-of-Pots Protocol.\nThe user is a participant who votes for a delegator who is responsible to forge the block.The user also bets his stake to get a proportion of the reward fees.\n\n")
	
	
	//Input stake of participant
	io.WriteString(conn, "Enter stake:")
	scanBalance := bufio.NewScanner(conn)
	for scanBalance.Scan() {
		balance, err := strconv.ParseFloat(scanBalance.Text(), 32)
		//strconv.Atoi(scanBalance.Text())
		if err != nil {
			log.Printf("%v not a number: %v", scanBalance.Text(), err)
			return
		}
		t := time.Now()
		address = calculateHash(t.String())
		
		validators[address] = float32(balance)
		//fmt.Println(validators)
		break
	}

	//Input data from participant


	//Asking user to enter delegate number
	io.WriteString(conn, "Enter your delegate[0-4]:")
	scanDelegate := bufio.NewScanner(conn)
	for scanDelegate.Scan() {
		sd, err = strconv.Atoi(scanDelegate.Text())
		if err != nil {
			//log.Printf("%v not a number: %v", scanDelegate.Text(), err)
			return
		}
		break
	}

	votes[sd]=votes[sd]+1
	//log.Println("Votes", votes[sd])
	
	

	
io.WriteString(conn, "\nEnter  new data: ")
	scanBPM := bufio.NewScanner(conn)
	go func() {

			// take in BPM from stdin and add it to blockchain after conducting necessary validation
			for scanBPM.Scan() {
				//bpm, err := strconv.Atoi(scanBPM.Text())
				// bpm := scanBPM.Text()
				data[d] = scanBPM.Text()
				completeData = completeData + address + " : " + scanBPM.Text() + " ; "
				// log.Println("complete data currently: ", completeData)
				//testimonyMap[validators[address]] = scanBPM.Text()
				// if malicious party tries to mutate the chain with a bad input, delete them as a validator and they lose their staked tokens
				// if err != nil {
				// 	log.Printf("%v not a number: %v", scanBPM.Text(), err)
				// 	delete(validators, address)
				// 	conn.Close()
				// }

				// mutex.Lock()
				// oldLastIndex := Blockchain[len(Blockchain)-1]
				// mutex.Unlock()

				// // create newBlock for consideration to be forged
				// newBlock, err := generateBlock(oldLastIndex, bpm, address)
				// if err != nil {
				// 	log.Println(err)
				// 	continue
				// }
				// if isBlockValid(newBlock, oldLastIndex) {
				// 	candidateBlocks <- newBlock
				// }
				break
			}
			// io.WriteString(conn, "\nEnter a new data: ")
			count=count+1
		
	}()

	

	// simulate receiving broadcast
	for {
		//io.WriteString(conn,"Your data is being forged")
		time.Sleep(30*time.Second)
		mutex.Lock()
		//log.Println("The length is: ",len(Blockchain),"and the length is ",length_block)
		//output, err := json.Marshal(Blockchain)
		//if length_block != len(Blockchain) {
		for i := range Blockchain{
			io.WriteString(conn, "\nBLOCK INDEX: "+ strconv.Itoa(Blockchain[i].Index)+"\n")
			io.WriteString(conn, "TIMESTAMP: " + Blockchain[i].Timestamp+"\n")
			io.WriteString(conn, "TESTIMONY: " + Blockchain[i].Testimony+"\n")
			io.WriteString(conn, "HASH: "+ Blockchain[i].Hash+"\n")
			io.WriteString(conn, "PREVIOUS HASH: " + Blockchain[i].PrevHash+"\n")
			io.WriteString(conn, "VALIDATOR: " + Blockchain[i].Validator+"\n")

			io.WriteString(conn, "\n\n\n")
		}
		//log.Println("iNCREASING LENGTH", length_block)
		//length_block=length_block + 1
		//log.Println("Length increased", length_block)

		mutex.Unlock()
		if err != nil {
			log.Fatal(err)
		}
		// io.WriteString(conn, string(output)+"\n")
	}


}

// isBlockValid makes sure block is valid by checking index
// and comparing the hash of the previous block
func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if calculateBlockHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

// SHA256 hasing
// calculateHash is a simple SHA256 hashing function
func calculateHash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

//calculateBlockHash returns the hash of all block information
func calculateBlockHash(block Block) string {
	record := string(block.Index) + block.Timestamp + string(block.Testimony) + block.PrevHash
	return calculateHash(record)
}

// generateBlock creates a new block using previous block's hash
func generateBlock(oldBlock Block, BPM string, address string) (Block, error) {

	var newBlock Block
	currentblockdata=completeData
	log.Println("\n GENERATING BLOCK")
	time.Sleep(2*time.Second)
	

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	//newBlock.BPM = BPM
	//TODO: 
	newBlock.Testimony = completeData
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateBlockHash(newBlock)
	newBlock.Validator = wv

	count = 0
	completeData = ""

	for r :=0; r<=9; r++{
		votes[r] = 0
	}



	return newBlock, nil
}

// We are intiliazing the total delegates in the system 
func generateAllDelegates() {


	for i := 0; i<=50; i++ {

		allDelegates[i] = calculateHash(time.Now().String())
	}

}
//We assign the new delegates in the round robin fashion
func chooseDelegates() {

	
	for i:=0;i<5;i++{
		if p>4 {
			
			delegates[p%5] = allDelegates[p]
			p=p+1
		} else{
			
			delegates[i] = allDelegates[i]
			p=p+1
		}
		
		
	}
	//log.Println("\nThe chosen deligates are\n",delegates)

	//delinit=delinit+5
}
func validatebyconsensus(newblock Block){
	log.Println("Block is getting validated.\n")
	time.Sleep(2*time.Second)
	delvotes=0
	newblock.Testimony=currentblockdata
	var hashstr string
	for i := range delegates{
		if delegates[i] != wv{
			
			hashstr=calculateBlockHash(newblock)
			
			
			if hashstr == newblock.Hash{
			delvotes=delvotes+1
		}

		}
	}
	if delvotes >= 3{
	log.Println("Consensus reached : - Block will be appended.")
	time.Sleep(2*time.Second)
	}

}