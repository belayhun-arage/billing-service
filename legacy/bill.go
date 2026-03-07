package main

import (
	"fmt"
	"os"
)

type bill struct {
	name  string
	items map[string]float64
	tip   float64
}

func newBill(name string) bill {
	myBill := bill{
		name:  name,
		items: map[string]float64{},
		tip:   0,
	}
	return myBill
}

func (b *bill) formatBill() string {
	formatedString := "@ቀሃስ-Restaurant:\n ...\n"
	price := 0.0
	for key, value := range b.items {
		formatedString += fmt.Sprintf("%-25v  $  :  %-25v \n", key, value)
		price += value
	}

	formatedString += fmt.Sprintf("total price $  :%-35v ", price+b.tip)

	return formatedString
}

func (b *bill) addTip(tip float64) {
	b.tip = tip
}

func (bill *bill) addItem(item string, price float64) {
	bill.items[item] = price
}

func (b *bill) saveBill() {
	fb := b.formatBill()
	billl := []byte(fb)

	err := os.WriteFile("bills/"+b.name+".txt", billl, 0644)
	if err != nil {
		fmt.Println("failed to save the bill")
	}
	fmt.Println("Successfully saved the bill")
}
