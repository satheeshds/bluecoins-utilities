package main

import (
	"bluecoins-to-splitwise-go/pkg/bank"
	"bluecoins-to-splitwise-go/pkg/bluecoins"
	"bluecoins-to-splitwise-go/pkg/cui"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jroimartin/gocui"
)

func main() {
	var filePath string
	var verbose bool

	// Define the command-line options
	flag.StringVar(&filePath, "file", "", "The path to the file")
	flag.BoolVar(&verbose, "v", false, "Enable verbose")

	// Parse the command-line options
	flag.Parse()

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Fatalf("The file does not exist at the path: %s", filePath)
	} else if err != nil {
		log.Fatalf("An error occurred: %v", err)
	}
	bluecoinsService, err := bluecoins.NewBluecoinsService()
	if err != nil {
		log.Fatalf("Error creating Bluecoins service: %v", err)
	}

	transactionService, err := bank.NewTransactionService(filePath)
	if err != nil {
		log.Fatalf("Error creating transaction service: %v", err)
	}
	transactions, err := transactionService.GetBankTransactions(filePath)
	if err != nil {
		log.Fatalf("Error getting bank transactions: %v", err)
	}
	// log.Printf("Transactions: %v \n", transactions)

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicf("Error creating GUI: %v", err)
	}
	defer g.Close()

	g.Cursor = true

	// Open output file
	f, err := os.OpenFile("output.log", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	mainView := &cui.MainView{
		Name:             "main",
		Transactions:     transactions,
		Logfile:          f,
		Verbose:          verbose,
		BluecoinsService: bluecoinsService,
	}
	g.SetManagerFunc(mainView.Layout)

	log.Println("in main loop...")
	err = g.MainLoop()
	if err != nil && err != gocui.ErrQuit {
		log.Println("Error in main loop:", err)
		log.Panicln(err)
	}

	g.Close()
	importTxns := mainView.GetSelectedTransactions()
	transactionService.WriteTransactionRecords(importTxns, "output.csv")

	fmt.Println("Press Enter to exit...")
	fmt.Scanln()
}
