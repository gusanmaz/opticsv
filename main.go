package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var multipleAnswerMap = map[byte]string{
	'F': "AB", 'G': "AC", 'H': "AD", 'I': "AE", 'J': "BC",
	'K': "BD", 'L': "BE", 'M': "CD", 'N': "CE", 'O': "DE", 'X': "ABCDE",
}

func compare(teacher, student string) int {
	correctCount := 0
	for question, answer := range teacher {
		studentAnswer := student[question]
		if studentAnswer == ' ' {
			continue
		}

		multipleAnswers, ok := multipleAnswerMap[byte(answer)]
		if ok {
			if strings.Contains(multipleAnswers, string(studentAnswer)) {
				correctCount++
			}
		} else {
			if byte(answer) == studentAnswer {
				correctCount++
			}
		}
	}
	return correctCount
}

func getPoint(answerKey string, studentAnswers string, ppq float64) int {
	correctCnt := compare(answerKey, studentAnswers)
	return int(float64(correctCnt) * ppq)
}

func processExam(datFilePath, keyFilePath string) {
	dataFile, err := os.Open(datFilePath)
	if err != nil {
		fmt.Printf("Error opening data file: %v\n", err)
		return
	}
	defer dataFile.Close()

	csvFile, err := os.Create(strings.TrimSuffix(datFilePath, ".dat") + ".csv")
	if err != nil {
		fmt.Printf("Error creating CSV file: %v\n", err)
		return
	}
	defer csvFile.Close()

	keyFile, err := os.Open(keyFilePath)
	if err != nil {
		fmt.Printf("Error opening key file: %v\n", err)
		return
	}
	defer keyFile.Close()

	scanner := bufio.NewScanner(keyFile)
	var answerKeys []string
	for scanner.Scan() {
		answerKeys = append(answerKeys, scanner.Text())
	}

	totalQuestions := len(answerKeys[0])
	ppq := 100.0 / float64(totalQuestions)

	dataScanner := bufio.NewScanner(dataFile)
	regExpStr := fmt.Sprintf(`([0-9]{1,%d})([A-D])((?:[A-E ]){%d})`, 10, totalQuestions)
	matchExp := regexp.MustCompile(regExpStr)

	studentCnt := 0
	gradeSum := 0

	for dataScanner.Scan() {
		line := dataScanner.Text()
		matches := matchExp.FindStringSubmatch(line)
		if matches == nil {
			fmt.Println("Invalid record:", line)
			continue
		}

		studentNo, sheetType, studentAnswers := matches[1], matches[2], matches[3]
		if len(studentNo) < 10 {
			fmt.Printf("Incomplete student number: %s. Record will not be included in CSV.\n", studentNo)
		}

		sheetIndex := strings.Index("ABCD", sheetType)
		if sheetIndex == -1 || sheetIndex >= len(answerKeys) {
			fmt.Println("Invalid sheet type for student:", studentNo)
			continue
		}

		point := getPoint(answerKeys[sheetIndex], studentAnswers, ppq)
		if len(studentNo) == 10 {
			fmt.Fprintf(csvFile, "%s;%d\n", studentNo, point)
		} else {
			fmt.Printf("Error for student with incomplete number %s: Calculated points %d\n", studentNo, point)
		}

		studentCnt++
		gradeSum += point
	}

	avg := float64(gradeSum) / float64(studentCnt)
	fmt.Printf("Processed %d students. Average score: %.2f\n", studentCnt, avg)
}

func main() {
	filename := flag.String("filename", "", "Base filename for the .dat and .key files (without extension)")
	flag.Parse()

	if *filename == "" {
		fmt.Println("Filename is required. Use -filename flag to specify it.")
		os.Exit(1)
	}

	datFilePath := *filename + ".dat"
	keyFilePath := *filename + ".key"

	processExam(datFilePath, keyFilePath)
}

//![Sekonic OMR Software SS](/images/sekonic_ss.png)
