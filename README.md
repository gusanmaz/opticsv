# Dat to CSV Converter for Sekonic OMR software
Software installed for Sekonic SR-1800 Optical Mark Reader(OMR) at my workplace sometimes doesn't function properly. From time to time it fails to generate a spreadsheet file from optical papers it read. Luckily it generates a dat file from readings even if fails to generate a proper spreadsheet file.

Interestingly it seems impossible to identify name of this OMR program from it's GUI. In liue you could find a screenshot of this program which could be helpful to figure out if they are using a similar OMR program so that this little Go code migth be handy.

![Sekonic OMR Software SS](/images/sekonic_ss.png)

# USAGE

Check *ExamInfo* variables defined above function definitions. Modify values of these variables that would reflect your case.

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



### Author
Guvenc Usanmaz

### License
This project is licensed under the MIT License.