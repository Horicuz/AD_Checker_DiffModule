package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// Store formatted output for side-by-side comparison
type FormattedOutput struct {
	reference []string
	output    []string
}

// Store differences for each file
type FileCompareResult struct {
	matched         bool
	diffs           []diffmatchpatch.Diff
	formattedOutput FormattedOutput
}

func readFile(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed reading file: %w", err)
	}
	return string(data), nil
}

func generateFormattedOutput(diffs []diffmatchpatch.Diff) FormattedOutput {
    var refText, outText string
    for _, diff := range diffs {
        switch diff.Type {
        case diffmatchpatch.DiffInsert:
            outText += fmt.Sprintf("\033[42m%s\033[0m", diff.Text) // Green background
            //refText += fmt.Sprintf("\033[41m%s\033[0m", strings.Repeat(" ", len(diff.Text))) // Red background for spaces
        case diffmatchpatch.DiffDelete:
            refText += fmt.Sprintf("\033[41m%s\033[0m", diff.Text) // Red background
            //outText += fmt.Sprintf("\033[42m%s\033[0m", strings.Repeat(" ", len(diff.Text))) // Green background for spaces
        case diffmatchpatch.DiffEqual:
            refText += diff.Text
            outText += diff.Text
        }
    }

    return FormattedOutput{
        reference: strings.Split(refText, "\n"),
        output:    strings.Split(outText, "\n"),
    }
}

func compareFilesInFolders(folder1, folder2 string, numFiles int) (int, map[string]FileCompareResult, error) {
	matchedCount := 0
	results := make(map[string]FileCompareResult)

	for i := 1; i <= numFiles; i++ {
		fileName := fmt.Sprintf("data%d.out", i)
		file1 := fmt.Sprintf("%s/%s", folder1, fileName)
		file2 := fmt.Sprintf("%s/%s", folder2, fileName)

		text1, err := readFile(file1)
		if err != nil {
			return 0, nil, fmt.Errorf("error reading file1: %w", err)
		}

		text2, err := readFile(file2)
		if err != nil {
			return 0, nil, fmt.Errorf("error reading file2: %w", err)
		}

		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(text1, text2, false)

		matched := len(diffs) == 1 && diffs[0].Type == diffmatchpatch.DiffEqual
		if matched {
			matchedCount++
		}

		results[fileName] = FileCompareResult{
			matched:         matched,
			diffs:           diffs,
			formattedOutput: generateFormattedOutput(diffs),
		}
	}

	return matchedCount, results, nil
}

func showDifferences(fileName string, result FileCompareResult, displayType int) {
	if result.matched {
		fmt.Printf("\033[32mFile %s: Files are identical\033[0m\n", fileName)
		return
	}

	fmt.Printf("\033[31mFile %s: Files are different\033[0m\n", fileName)

	if displayType == 1 {
		// Original inline display
		for _, diff := range result.diffs {
			switch diff.Type {
			case diffmatchpatch.DiffInsert:
				fmt.Printf("\033[32m%s\033[0m", diff.Text)
			case diffmatchpatch.DiffDelete:
				fmt.Printf("\033[31m%s\033[0m", diff.Text)
			case diffmatchpatch.DiffEqual:
				fmt.Printf("%s", diff.Text)
			}
		}
	} else if displayType == 2 {
		fmt.Println("\nReference:")
		fmt.Println("----------")
		for _, line := range result.formattedOutput.reference {
			if line != "" {
				fmt.Println(line)
			}
		}
		
		fmt.Println("\nOutput:")
		fmt.Println("-------")
		for _, line := range result.formattedOutput.output {
			if line != "" {
				fmt.Println(line)
			}
		}
	}
	fmt.Println()
}

func countFilesInFolder(folder string) (int, error) {
    pattern := filepath.Join(folder, "data*.out")
    matches, err := filepath.Glob(pattern)
    if err != nil {
        return 0, fmt.Errorf("error counting files: %w", err)
    }
    return len(matches), nil
}

func showIncorrectFiles(results map[string]FileCompareResult) {
    fmt.Println("\nIncorrect files:")
    fmt.Println("---------------")
    hasIncorrect := false
    
    // Sort the files numerically
    incorrectFiles := make([]int, 0)
    for fileName, result := range results {
        if (!result.matched) {
            fileNum := 0
            fmt.Sscanf(fileName, "data%d.out", &fileNum)
            incorrectFiles = append(incorrectFiles, fileNum)
            hasIncorrect = true
        }
    }
    
    if !hasIncorrect {
        fmt.Println("None - all files are correct!")
        return
    }

    // Sort numbers for better readability
    sort.Ints(incorrectFiles)
    
    // Print in a clean format
    for _, num := range incorrectFiles {
        fmt.Printf("- data%d.out\n", num)
    }
    fmt.Println()
}

func main() {

	//config = utils.Config.UserConfig

	folder1 := "./LastRef" //Folder for reference // config.ref_path
	folder2 := "./OutputData" //User output folder // config.output_path
	
	numFiles, err := countFilesInFolder(folder1)
	if (err != nil) {
		fmt.Println("Error counting files:", err)
		return
	}

	if numFiles == 0 {
		fmt.Println("No files found in reference folder")
		return
	}

	matchedCount, results, err := compareFilesInFolders(folder1, folder2, numFiles)
	if err != nil {
		fmt.Println("Error comparing files:", err)
		return
	}

	fmt.Printf("Matched files: %d/%d\n", matchedCount, numFiles)

	if matchedCount == numFiles {
		fmt.Println("\033[32mSUCCESS! - All tasks tested!\033[0m")
	} else {
		fmt.Println("\033[31mSome files have differences!\033[0m")
		showIncorrectFiles(results)
	}

	var continueViewing = true
	for continueViewing {
		// Ask user if they want to see specific differences
		fmt.Print("\nWould you like to see differences for a file? (y/n): ")
		var response string
		fmt.Scanln(&response)

		if response != "y" && response != "Y" {
			continueViewing = false
			continue
		}

		fmt.Print("Enter file number (1-34): ")
		var fileNum int
		fmt.Scanln(&fileNum)

		if fileNum < 1 || fileNum > numFiles {
			fmt.Printf("Invalid file number. Please enter a number between 1 and %d\n", numFiles)
			continue
		}

		fmt.Print("Enter display type (1 for inline, 2 for side by side): ")
		var displayType int
		fmt.Scanln(&displayType)

		if displayType != 1 && displayType != 2 {
			fmt.Println("Invalid display type")
			continue
		}

		fileName := fmt.Sprintf("data%d.out", fileNum)
		if result, exists := results[fileName]; exists {
			showDifferences(fileName, result, displayType)
		}
	}
	
	fmt.Println("Program finished. Goodbye!")
}
