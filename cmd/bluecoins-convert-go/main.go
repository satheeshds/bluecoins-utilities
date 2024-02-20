package main

import (
	"bluecoins-to-splitwise-go/pkg/bank"
	"bluecoins-to-splitwise-go/pkg/cui"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jroimartin/gocui"
)

func main() {
	var filePath string

	// Define the command-line options
	flag.StringVar(&filePath, "file", "", "The path to the file")

	// Parse the command-line options
	flag.Parse()

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Fatalf("The file does not exist at the path: %s", filePath)
	} else if err != nil {
		log.Fatalf("An error occurred: %v", err)
	}

	transactionService, err := bank.NewTransactionService()
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

	mainView := &cui.MainView{
		Name:         "main",
		Transactions: transactions,
	}
	g.SetManagerFunc(mainView.Layout)
	if err := mainView.Create(g); err != nil {
		log.Panicln(err)
	}

	log.Println("in main loop...")
	err = g.MainLoop()
	if err != nil && err != gocui.ErrQuit {
		log.Println("Error in main loop:", err)
		log.Panicln(err)
	}

	// // Open output file
	// f, err := os.OpenFile("output.txt", os.O_RDWR|os.O_CREATE, 0755)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer f.Close()

	mainView.GetSelectedTransactions()

	fmt.Println("Press Enter to exit...")
	fmt.Scanln()
	// var descriptions []string
	// for _, transaction := range transactions {
	// 	descriptions = append(descriptions, transaction.Description)
	// }
	// list := cui.NewSelectableList("list", descriptions, 0, 0)

	// g.SetManagerFunc(list.Layout)

	// if err := g.SetKeybinding("list", gocui.KeyArrowDown, gocui.ModNone, list.Down); err != nil {
	// 	log.Fatalf("Error setting keybinding: %v", err)
	// }

	// if err := g.SetKeybinding("list", gocui.KeyArrowUp, gocui.ModNone, list.Up); err != nil {
	// 	log.Fatalf("Error setting keybinding: %v", err)
	// }

	// if err := g.SetKeybinding("list", gocui.KeyEnter, gocui.ModNone, list.Enter); err != nil {
	// 	log.Fatalf("Error setting keybinding: %v", err)
	// }

	// log.Println("in main loop...")
	// if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
	// 	log.Fatalf("Error in main loop: %v", err)
	// }

	// log.Printf("Selected item: %s", list.GetSelected())
	// log.Println("Exiting...")
}
