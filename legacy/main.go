package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func acceptUsrInput(text string, reader *bufio.Reader) (string, error) {
	fmt.Print(text)
	input, err := reader.ReadString('\n')
	return strings.TrimSpace(input), err
}

func createNewBill() bill {
	reader := bufio.NewReader(os.Stdin)
	name, _ := acceptUsrInput("Creating new bill :\n ....Please enter the name:", reader)
	myBill := newBill(name)

	return myBill
}

func handleUserOptions(b bill) {
	reader := bufio.NewReader(os.Stdin)

	//options
	options, _ := acceptUsrInput("Please Enter (a---to add items t---to add tips s---to save bill):  ", reader)
	opt := strings.TrimSpace(options)
	switch opt {
	case "a":
		item, _ := acceptUsrInput("item name here: ", reader)
		price, _ := acceptUsrInput("item price here: ", reader)
		p, err := strconv.ParseFloat(price, 64)
		if err != nil {
			fmt.Println("price must be a number please enter again")
			handleUserOptions(b)
		}
		b.addItem(item, p)
		handleUserOptions(b)
	case "t":
		tip, _ := acceptUsrInput("tip amount here: ", reader)
		t, err := strconv.ParseFloat(tip, 64)
		if err != nil {
			fmt.Println("price must be a number please enter again")
			handleUserOptions(b)
		}
		b.addTip(t)
		handleUserOptions(b)
	case "s":
		b.saveBill()
		fmt.Println("You have successfully saved the bill.")
	default:
		fmt.Println("Please fill again")
		handleUserOptions(b)
	}

}

func main() {
	mybill := createNewBill()
	handleUserOptions(mybill)
}
