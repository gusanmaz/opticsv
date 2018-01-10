package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// If a question has multiple answers it is represented with a different letter than multiple choice letters
// E.g. if both choice A and D is correct answer of the regarding question is denoted with letter H
var multipleAnswerMap = map[byte]string{
	'F': "AB",
	'G': "AC",
	'H': "AD",
	'I': "AE",
	'J': "BC",
	'K': "BD",
	'L': "BE",
	'M': "CD",
	'N': "CE",
	'O': "DE",
	'X': "ABCDE",
}

type ExamInfo struct {
	AnswerKeys            []string
	QuestionNumber        int
	SingleQuestionPoint   float64
	StudentNumberLen      int
	DatFilePath           string
	CsvFilePath           string
	MultipleChoiceLetters string
	ExamSheetLetters      string
}

/* Change values of ExamInfo variables according to your needs */

var orgunExamInfo = ExamInfo{
	AnswerKeys: []string{
		"CEAABCDCBCEEABEACBEBECDACEBDDECCDAAEDELABAXCB",
		"DACEBDDECCDCEAABCDCBCEEDELABAXCBABEACBEBECAAE",
		"BAXCBBCEEABEACBDEBECCEAABCDCAAEACEBDDECCDDELA",
		"BCEEDDECCDACEBDDELABAXCBABEACBEBECAAECEAABCDC"},
	QuestionNumber:        45,
	SingleQuestionPoint:   2,
	DatFilePath:           "oop1.dat",
	CsvFilePath:           "oop1.csv",
	MultipleChoiceLetters: "ABCDE",
	ExamSheetLetters:      "ABCD",
	StudentNumberLen:      10,
}

var geceExamInfo = ExamInfo{
	AnswerKeys: []string{
		"EDBDDDDEAABADCDEDEBECDCACBDAEECDAAAEECAABAXCB",
		"DCDEDEAEECDDBDDDDEAABABAXCBEBECDCACBDEECAAAAA",
		"BAXCBEBECDEDBDDDEDCDEDDCACBDAEECDAAAEECAAAABA",
		"AEECDAAABDDDEDDEDBAXCBEBECDEDCDCACBDAABAEECAA"},
	QuestionNumber:        45,
	SingleQuestionPoint:   2,
	DatFilePath:           "oop2.dat",
	CsvFilePath:           "oop2.csv",
	MultipleChoiceLetters: "ABCDE",
	ExamSheetLetters:      "ABCD",
	StudentNumberLen:      10,
}

// Readln returns a single line (without the ending \n)
// from the input buffered reader.
// An error is returned iff there is an error with the
// buffered reader.
func Readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}

func compare(teacher, student string) int {
	var correctCount = 0

	for question, answer := range teacher {
		multipleAnswers, ok := multipleAnswerMap[byte(answer)]
		if ok {
			if strings.Index(multipleAnswers, string(student[question])) != -1 {
				correctCount++
			}
		} else {
			if byte(answer) == student[question] {
				correctCount++
			}
		}
	}
	return correctCount
}

func getPoint(answerKeys []string, studentAnswers string, sheetType string, ppq float64) int {
	// Modify sheetLetters string if your exam sheets have different bear different set of letters!
	sheetLetters := "ABCDE"
	sheetIndex := strings.Index(sheetLetters, sheetType)
	// In case sheet type is not marked on student answer key
	if sheetIndex == -1 {
		sheetIndex = 0
	}
	answerKey := answerKeys[sheetIndex]
	correctCnt := compare(answerKey, studentAnswers)
	return int(float64(correctCnt) * ppq)
}

func processExam(info ExamInfo) {
	/* Below could be found an excerpt of content of a anonymised dat file:
	    JOHN DOE            2  9174968307DBDECEAECCDAAEBCDACABAACBABEADBEBECAAECEADBCDC
		JANE DOE            2  2142656541AADAEDCBDCDECACBADCDEBAEACDBABAEECAAEAAAABAAEB
		FOO BAR             2  3439648052DBDEBEBECBDACEBDDCCABAACBADEAEBEBDCAAADBAAADBE
		GOO FOO             2  5657676027BAACEDCBECBDCEAAEADDBCECBDEABAADBABEDCECABAAAE
		BAR JANE            2  7758686991BBBDEBCCBBDDADDEBADCBAEBEBACCCACECAABEBCCDEACB
	*/

	dataFile, err := os.Open(info.DatFilePath)
	defer dataFile.Close()

	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		os.Exit(1)
	}

	csvFile, err := os.Create("orgun.csv")
	defer csvFile.Close()

	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		os.Exit(1)
	}

	answerKeys := info.AnswerKeys

	regExpStr := "([0-9]{" + strconv.Itoa(info.StudentNumberLen) + "})([" + info.ExamSheetLetters + "])([* " +
		info.MultipleChoiceLetters + "]{" + strconv.Itoa(info.QuestionNumber) + "})"
	//var matchExp = regexp.MustCompile(`([0-9]{10})([ABCD])([ABCDE* ]{50})`)
	var matchExp = regexp.MustCompile(regExpStr)

	dataReader := bufio.NewReader(dataFile)
	s, e := Readln(dataReader)
	studentCnt := 0
	gradeSum := 0
	for e == nil {
		matches := matchExp.FindStringSubmatch(s)
		studentNo := matches[1]
		letter := matches[2]
		studentAnswers := matches[3]

		point := getPoint(answerKeys, studentAnswers, letter, info.SingleQuestionPoint)
		fmt.Println(studentNo + "->" + strconv.Itoa(point))
		fileLine := studentNo + ";" + strconv.Itoa(point) + "\n"
		csvFile.WriteString(fileLine)

		s, e = Readln(dataReader)
		studentCnt++
		gradeSum += point
	}

	avg := float64(gradeSum) / float64(studentCnt)
	avgStr := strconv.FormatFloat(avg, 'f', 6, 64)
	fmt.Println("Student Number: " + strconv.Itoa(studentCnt) + " Avg: " + avgStr)
}

func main() {

	examInfos := []ExamInfo{orgunExamInfo, geceExamInfo}

	for _, examInfoVal := range examInfos {
		processExam(examInfoVal)
	}

}
