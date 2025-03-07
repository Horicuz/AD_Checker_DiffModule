# Go Output Diff Checker

A streamlined command-line tool for comparing output files, featuring both inline and side-by-side difference views. Perfect for testing and validation workflows. 

The project uses the algorithm for detecting differences between two blocks of text made by @sergi (https://github.com/sergi/go-diff) 

## Overview

This tool compares files between a reference directory (`LastRef`) and an output directory (`OutputData`), highlighting differences with color-coded output. Great for:
- Automated testing validation
- Output comparison
- File difference analysis

## Quick Setup

```bash
# Clone the repository
git clone https://github.com/horicuz/go-diff-checker.git
cd go-diff-checker

# Install required dependency
go get github.com/sergi/go-diff/diffmatchpatch
```
## Directory Structure

```
your-project/
├── LastRef/      # Reference files
│   └── dataN.out
└── OutputData/   # Files to compare
    └── dataN.out
```
## How to use

  1. Place reference files in LastRef directory
  2. Place files to check in OutputData directory
  3. Run the program:

```bash
go run DiffChecker.go
```

## Features

### Two View Modes:
    
  - Inline comparison
  - Side-by-side view
    
### Color Coding:
  - Green: Added content
  - Red: Removed content
    
### Interactive Interface:
  - File selection
  - View mode selection
    
### Smart Analysis:
  - Success rate calculation
  - Sorted incorrect file listing

## Credit
Diff functionality powered by go-diff by Sergi Mansilla.

## License
MIT License

