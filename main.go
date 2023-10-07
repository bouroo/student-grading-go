package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
)

type Grade string

const (
	A Grade = "A"
	B Grade = "B"
	C Grade = "C"
	F Grade = "F"
)

type student struct {
	firstName, lastName, university                string
	test1Score, test2Score, test3Score, test4Score int
}

type studentStat struct {
	student
	finalScore float32
	grade      Grade
}

func parseCSV(filePath string) []student {
	// Open the CSV file
	csvFile, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("failed to open CSV file: %s", err)
		return nil
	}
	defer csvFile.Close()

	// Create a CSV reader
	csvReader := csv.NewReader(csvFile)

	// Create an empty slice to store student data
	students := make([]student, 0)

	// Loop through each line in the CSV file
csvLineLoop:
	for {
		// Read a line from the CSV file
		line, err := csvReader.Read()
		if err != nil {
			// If there is an error reading the line, check if it's end of file
			if err == io.EOF {
				// If it's end of file, break out of the loop
				break
			}
			log.Fatalf("failed to read CSV file: %s", err)
			return nil
		}

		// Create an array to store the test scores
		var testScores [4]int

		// Loop through the test score fields in the line
		for i := 3; i <= 6; i++ {
			// Convert the test score from string to integer
			testScores[i-3], err = strconv.Atoi(line[i])
			if err != nil {
				// If the test score is not a number, skip to the next line
				continue csvLineLoop
			}
		}

		// Create a new student object and append it to the slice
		students = append(students, student{
			firstName:  line[0],
			lastName:   line[1],
			university: line[2],
			test1Score: testScores[0],
			test2Score: testScores[1],
			test3Score: testScores[2],
			test4Score: testScores[3],
		})
	}

	// Return the slice of students
	return students
}

// calculateGrade calculates the final score and grade for each student in the given slice of students.
// It returns a slice of studentStat which contains the student information along with their final score and grade.
func calculateGrade(students []student) []studentStat {
	// Create a slice to store the student statistics
	studentStats := make([]studentStat, len(students))

	// Iterate over each student in the students slice
	for i, student := range students {
		// Calculate the final score by averaging the test scores
		finalScore := float32(student.test1Score+student.test2Score+student.test3Score+student.test4Score) / 4

		// Determine the grade based on the final score
		var grade Grade
		switch {
		case finalScore >= 70:
			grade = A
		case finalScore >= 50:
			grade = B
		case finalScore >= 35:
			grade = C
		case finalScore < 35:
			grade = F
		}

		// Create a studentStat object with the student information, final score, and grade
		studentStats[i] = studentStat{
			student:    student,
			finalScore: finalScore,
			grade:      grade,
		}
	}

	// Return the slice of student statistics
	return studentStats
}

// findOverallTopper is a function that takes in a slice of studentStat structs
// and returns the studentStat struct with the highest finalScore.
func findOverallTopper(gradedStudents []studentStat) studentStat {
	// Initialize the variable "topper" with the first student in the slice.
	topper := gradedStudents[0]

	// Iterate through the remaining students in the slice.
	for i := 1; i < len(gradedStudents); i++ {
		// Check if the current student's finalScore is greater than the current topper's finalScore.
		if gradedStudents[i].finalScore > topper.finalScore {
			// If it is, update the topper to be the current student.
			topper = gradedStudents[i]
		}
	}

	// Return the student with the highest finalScore.
	return topper
}

// findTopperPerUniversity the topper from each university
func findTopperPerUniversity(gs []studentStat) map[string]studentStat {
	// Create a map to store the topper student from each university
	toppers := make(map[string]studentStat)

	// Iterate over the list of students
	for _, student := range gs {
		// Check if there is already a topper student for the university
		if topper, ok := toppers[student.university]; !ok || student.finalScore > topper.finalScore {
			// If there is no topper student or the current student has a higher final score,
			// update the topper student for the university
			toppers[student.university] = student
		}
	}

	// Return the map containing the topper student from each university
	return toppers
}
