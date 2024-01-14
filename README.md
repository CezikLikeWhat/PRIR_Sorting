# Programowanie równoległe i rozproszone - Go-Sort

## Overview
Go-Sort is a CLI application that provides a suite of commands for handling binary files and sorting operations. The application is built in Go and features three main commands:

1. **Generate:** Creates a sample binary file in a specific format (the first number as `int64` indicating the number of subsequent `int32` numbers).
2. **Sort:** Sorts a given binary file using the sample sort algorithm with concurrent sorting techniques and a k-ways merge algorithm implemented using a min-heap.
3. **Verify:** Sequentially verifies the sorted file for correctness.

## Requirements
- Go version 1.21.4
- Dependencies:
    - `github.com/fatih/color v1.16.0`
    - `github.com/spf13/cobra v1.8.0`

## Installation
To install Go-Sort, you need to have Go installed on your system (Follow the instructions: https://go.dev/dl/). 

If you have Go set up, follow these steps:

1. Clone the Go-Sort repository:
   ```bash
   git clone https://github.com/CezikLikeWhat/PRIR_Sorting.git
   ```
2. Change directory:
   ```bash
   cd PRIR_Sorting
   ```
3. Build a project for your OS:
   ```bash
   Make build-[mac, linux, windows]
   ```
   or build for all operating systems:
   ```bash
   Make build
   ```
4. Run application:
   ```bash
   ./bin/<OS>-<arch>/go-sort help
   ```

## License
MIT