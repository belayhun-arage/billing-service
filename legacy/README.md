# Legacy: Console Billing App

This directory contains the original console-based billing application — the starting point of this project, written roughly 4 years ago as a Go learning exercise.

## What it does

A simple CLI app for a restaurant (@ቀሃስ-Restaurant) that lets you:
- Create a named bill
- Add items with prices
- Add a tip
- Save the formatted bill to a `.txt` file in the `bills/` directory

## Files

| File | Description |
|------|-------------|
| `bill.go` | `bill` struct with `addItem`, `addTip`, `formatBill`, `saveBill` |
| `main.go` | CLI loop using `bufio.Reader` for user input |
| `bills/` | Sample saved bill files |

## Why it's still here

Kept for git history — it represents the very first commit of this repository and shows the evolution from a simple CLI script to the production-grade REST API in the parent directory.
