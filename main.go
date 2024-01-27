package main

import (
	"bufio"
	"flag"
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type AnswerKey struct {
	Answer string
	Points float64
}

var multipleAnswerMap = map[rune]string{
	'F': "AB", 'G': "AC", 'H': "AD", 'I': "AE", 'J': "BC",
	'K': "BD", 'L': "BE", 'M': "CD", 'N': "CE", 'O': "DE", 'X': "ABCDE",
}

func main() {
	filename := flag.String("filename", "", "Base filename for the .dat and .key files (without extension)")
	totalPoints := flag.Float64("totalPoints", 100, "Total points for the exam")
	flag.Parse()

	answerKeys, err := readAnswerKeys(*filename+".key", *totalPoints)
	if err != nil {
		fmt.Printf("Error reading answer keys: %v\n", err)
		return
	}

	processedStudents, averageScore, err := processExam(*filename+".dat", *filename+".csv",
		*filename+"_detailed.csv", answerKeys)
	if err != nil {
		fmt.Printf("Error processing exam: %v\n", err)
		return
	}

	fmt.Printf("Processed %d students. Average score: %.2f\n", processedStudents, averageScore)
}

func readAnswerKeys(keyFilePath string, totalPoints float64) (map[rune][]AnswerKey, error) {
	file, err := os.Open(keyFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	answerKeys := make(map[rune][]AnswerKey)
	scanner := bufio.NewScanner(file)
	sheetTypes := "ABCDE"
	answerPattern := regexp.MustCompile(`([A-X])(:\d+)?`)

	for i := 0; scanner.Scan() && i < len(sheetTypes); i++ {
		line := scanner.Text()
		matches := answerPattern.FindAllStringSubmatch(line, -1)

		if len(matches) == 0 {
			return nil, fmt.Errorf("no valid answers found in line for sheet %c", sheetTypes[i])
		}

		var keys []AnswerKey
		var totalDefinedPoints float64

		for _, match := range matches {
			answer := match[1]
			var points float64
			if match[2] != "" {
				var err error
				pointsStr := strings.TrimPrefix(match[2], ":")
				points, err = strconv.ParseFloat(pointsStr, 64)
				if err != nil {
					return nil, fmt.Errorf("invalid points format for answer %s in sheet %c", answer, rune(sheetTypes[i]))
				}
				totalDefinedPoints += points
			} else {
				// Points not specified, assign a placeholder value for now
				points = -1
			}
			keys = append(keys, AnswerKey{Answer: answer, Points: points})
		}

		// Distribute remaining points among answers without specified points
		distributeRemainingPoints(keys, totalPoints, totalDefinedPoints)

		answerKeys[rune(sheetTypes[i])] = keys
	}

	return answerKeys, scanner.Err()
}

func distributeRemainingPoints(keys []AnswerKey, totalPoints, totalDefinedPoints float64) {
	remainingPoints := totalPoints - totalDefinedPoints
	var countWithoutPoints int
	for _, key := range keys {
		if key.Points == -1 {
			countWithoutPoints++
		}
	}
	for i := range keys {
		if keys[i].Points == -1 {
			keys[i].Points = remainingPoints / float64(countWithoutPoints)
		}
	}
}

func processExam(datFilePath, csvFilePath, detailedCsvFilePath string, answerKeys map[rune][]AnswerKey) (int, float64, error) {
	dataFile, err := os.Open(datFilePath)
	if err != nil {
		return 0, 0, err
	}
	defer dataFile.Close()

	reader := bufio.NewReader(charmap.Windows1254.NewDecoder().Reader(dataFile))

	csvFile, err := os.Create(csvFilePath)
	if err != nil {
		return 0, 0, err
	}
	defer csvFile.Close()

	detailedCsvFile, err := os.Create(detailedCsvFilePath)
	if err != nil {
		return 0, 0, err
	}
	defer detailedCsvFile.Close()

	var totalScore float64
	var processedStudents int

	questionNumbers := len(answerKeys['A'])

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}

		studentID, studentName, sheetType, answers, ok := parseLine(string(line), questionNumbers)
		if !ok {
			continue
		}

		score := calculateScore(answers, answerKeys[sheetType])
		roundedScore := int(math.Round(score))

		_, err = fmt.Fprintf(csvFile, "%s,%d\n", studentID, roundedScore)
		if err != nil {
			return processedStudents, totalScore / float64(processedStudents), err
		}

		_, err = fmt.Fprintf(detailedCsvFile, "%s,%s,%d\n", studentName, studentID, roundedScore)
		if err != nil {
			return processedStudents, totalScore / float64(processedStudents), err
		}

		processedStudents++
		totalScore += score
	}

	return processedStudents, totalScore / float64(processedStudents), nil
}

func parseLine(line string, questionNumber int) (studentID string, studentName string, sheetType rune, answers string, ok bool) {
	runes := []rune(line)

	// Ensure the runes slice is long enough before slicing
	if len(runes) < 33 {
		fmt.Printf("Line too short: %s\n", line)
		return "", "", ' ', "", false
	}

	// Extract studentName and studentID using rune slicing
	studentName = string(runes[:20])
	studentIDRaw := string(runes[23:33])
	studentID = strings.TrimSpace(studentIDRaw)

	// Ensure there's enough length for at least one character of answers
	if len(runes) <= 33 {
		fmt.Printf("Line too short to include answers: %s\n", line)
		return "", "", ' ', "", false
	}

	// Extract sheetType and answers
	sheetTypeChar := runes[33]
	if sheetTypeChar >= 'A' && sheetTypeChar <= 'E' {
		sheetType = rune(sheetTypeChar)
	} else {
		sheetType = 'A'
	}

	if len(runes) > 34 {
		answers = string(runes[34:(34 + questionNumber)])
	} else {
		answers = ""
	}

	return studentID, studentName, sheetType, answers, true
}

func calculateScore(answers string, keys []AnswerKey) float64 {
	var score float64
	for i, answer := range answers {
		if i >= len(keys) {
			break
		}
		correctAnswer := keys[i].Answer

		// Convert the correctAnswer to a rune if it's a single character
		if len(correctAnswer) == 1 {
			correctRune := rune(correctAnswer[0])

			// Check if the answer is one of the multiple correct answers
			if multipleAnswers, ok := multipleAnswerMap[correctRune]; ok {
				if strings.ContainsRune(multipleAnswers, rune(answer)) {
					score += keys[i].Points
				}
			} else if string(answer) == correctAnswer {
				// Check if the answer matches the key directly
				score += keys[i].Points
			}
		}
	}
	return score
}
