# Dat to CSV Converter for Sekonic OMR Software

## Overview
This tool serves as a backup solution for processing data from the Sekonic SR-1800 Optical Mark Reader (OMR) software used at our workplace. Sometimes the official software may fail to generate spreadsheet files from scanned optical papers, even though it can produce `.dat` files with the scan results. This converter takes those `.dat` files, along with corresponding `.key` files containing the answer keys, and generates a `.csv` file listing student IDs and their calculated scores.

![Sekonic OMR Software Screenshot](images/sekonic_ss.png)
*Note: Replace "path/to/screenshot.jpg" with the actual path to the screenshot image.*

## Installation

### Prerequisites
- **Go Programming Language**: This tool is implemented in Go. Ensure you have Go installed on your system to compile and execute the program. Go can be downloaded from the [official Go website](https://golang.org/dl/).

### Compiling the Program
1. **Clone the Repository**: First, clone the repository to your local machine or download the source code directly. Use the following command to clone the repository:
    ```sh
    git clone https://github.com/gusanmaz/opticsv.git
    cd optikcsv
    ```

2. **Build the Program**: Compile the source code using the Go compiler. This will generate an executable file named `dat-to-csv-converter` (or `dat-to-csv-converter.exe` on Windows systems).
    ```sh
    go build -o opticsv
    ```

### Running the Program
The program requires the `-filename` argument to specify the base name of the `.dat` and `.key` files (without their extensions). For instance, if your files are named `exam1.dat` and `exam1.key`, execute the program as follows:
```sh
./opticsv -filename exam1
```
This command will process the exam1.dat and exam1.key files located in the current directory and generate an exam1.csv file with student scores.
Usage

The tool is designed to be straightforward. Specify the base filename of your .dat and .key files using the -filename flag when running the program from the command line. The tool will automatically calculate the scores based on the provided answer keys and generate a .csv file with the results.
Configuration

The converter automatically determines the number of questions from the length of the answer keys in the .key file and assumes each exam's total score is 100 points, distributing the points evenly across all questions.
Notes

The tool gracefully handles cases where student numbers are incomplete or certain answers are left unmarked, ensuring these instances are reported but not included in the final .csv file.
For exams with multiple versions of answer sheets, the .key file can contain multiple lines, each corresponding to a different version (A, B, C, D, etc.).

### Author

Güvenç Usanmaz

### License

This project is licensed under the MIT License - see the LICENSE.md file for details.