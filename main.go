package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

//var validEmailDomains = []string{"com", "edu", "net"}
var validGenders = []string{"m", "f", "male", "female"}

type Student struct {
	StudentID string
	Name string
	Age int
	Gender string
	Email string
	CoursesEnrolled []Course
	Grades map[string]int
}

type Course struct {
	CourseID string
	CourseName string
	Instructor string
}

func main() {
	student := getStudent()
	runApp(&student)
}

func showMenu() int {
	fmt.Println("-----Hello, We Welcome you to our program!-----")
	fmt.Println("-----Choose from the list below!-----")
	fmt.Println("1. Enroll in a course")
	fmt.Println("2. Remove a course")
	fmt.Println("3. Add a grade")
	fmt.Println("4. Update a grade")
	fmt.Println("5. View all enrolled courses")
	fmt.Println("6. Get GPA")
	fmt.Println("7. View all grades")
	fmt.Println("8. Exit the app")
	fmt.Println("--------------------------------------")

	fmt.Print("Choose from 1 - 8: ")
	choice, err := strconv.Atoi(readInput())
	if err != nil || choice < 1 || choice > 8 {
		fmt.Println("Ops! Invalid choice! Choose from 1 - 8.")
		return 0
	}
	return choice
}

func getStudent() Student {
	var student Student
	for {
		fmt.Print("Enter your full name: ")
		student.Name = strings.TrimSpace(readInput())
		if !isValidName(student.Name) {
			fmt.Println("Name is invalid!")
			continue
		}

		fmt.Print("Enter email: ")
		student.Email = strings.TrimSpace(readInput())
		if !isValidEmail(student.Email) {
			fmt.Println("Email is invalid!")
			continue
		}

		fmt.Print("Enter your student ID: ")
		student.StudentID = strings.TrimSpace(readInput())
		if !isValidStudentID(student.StudentID) {
			fmt.Println("Student ID is invalid!")
			continue
		}

		fmt.Print("Enter your age: ")
		ageInput := strings.TrimSpace(readInput())
		age, err := strconv.Atoi(ageInput)
		if err != nil {
			fmt.Println("Age must be a number")
			continue
		}
		student.Age = age

		fmt.Print("Enter your gender: ")
		student.Gender = strings.TrimSpace(readInput())
		if !isValidGender(student.Gender) {
			fmt.Println("Gender is invalid!")
			continue
		}

		student.Grades = make(map[string]int)
		return student
	}
}

func runApp(student *Student) {
	for {
		choice := showMenu()
		switch choice {
		case 1:
			student.AddCourse()
		case 2:
			student.RemoveCourse()
		case 3:
			student.AddGrade()
		case 4:
			student.UpdateGrade()
		case 5:
			student.PrintAllCourses()
		case 6:
			student.GetGPA()
		case 7:
			student.PrintGrades()
		case 8:
			fmt.Println("You exited the app")
			return
		default:
			fmt.Println("Invalid choice! Choose from 1 - 8.")
		}
	}
}

func (s *Student) AddCourse() {
	course, err := getCourseInput()
	if err != nil {
		fmt.Println(err)
		return
	}
	s.CoursesEnrolled = append(s.CoursesEnrolled, course)
}

func (s *Student) RemoveCourse() {
	if len(s.CoursesEnrolled) == 0 {
		fmt.Println("No courses enrolled")
		return
	}

	fmt.Print("Enter the ID of the course you want to remove: ")
	courseID := strings.TrimSpace(readInput())

	index, found := s.findCourseIndexByID(courseID)
	if found {
		s.CoursesEnrolled = append(s.CoursesEnrolled[:index], s.CoursesEnrolled[index+1:]...)
		delete(s.Grades, courseID)
		fmt.Printf("Course with ID %s removed!\n", courseID)
	} else {
		fmt.Printf("Course with ID %s not found\n", courseID)
	}
}

func (s *Student) AddGrade() {
	fmt.Print("Enter the course ID to add grade: ")
	courseID := strings.TrimSpace(readInput())

	fmt.Print("Enter the new grade: ")
	gradeInput := strings.TrimSpace(readInput())
	grade, err := strconv.Atoi(gradeInput)
	if err != nil {
		fmt.Println("Grade must be a number")
		return
	}

	if _, found := s.findCourseByID(courseID); found {
		s.Grades[courseID] = grade
		fmt.Printf("Grade %d added to course %s\n", grade, courseID)
	} else {
		fmt.Printf("Course with ID %s not found\n", courseID)
	}
}

func (s *Student) UpdateGrade() {
	fmt.Print("Enter the course ID to update grade: ")
	courseID := strings.TrimSpace(readInput())

	fmt.Print("Enter the new grades: ")
	gradeInput := strings.TrimSpace(readInput())
	grade, err := strconv.Atoi(gradeInput)
	if err != nil {
		fmt.Println("Grade must be number")
		return
	}

	if _, found := s.findCourseByID(courseID); found {
		s.Grades[courseID] = grade
		fmt.Printf("Grade updated to %d for course %s\n", grade, courseID)
	} else {
		fmt.Printf("Course with ID %s not found\n", courseID)
	}
}

func (s *Student) PrintGrades() {
	if len(s.Grades) == 0 {
		fmt.Println("No grades recorded")
		return
	}
	for courseID, grade := range s.Grades {
		fmt.Printf("Course ID: %s, Grade: %d\n", courseID, grade)
	}
}

func (s *Student) GetGPA() {
	if len(s.Grades) == 0 {
		fmt.Println("No grades recorded")
		return
	}
	total := 0
	for _, grade := range s.Grades {
		total += grade
	}
	gpa := float64(total) / float64(len(s.Grades))
	fmt.Printf("Your GPA: %.2f\n", gpa)
}

func (s *Student) PrintAllCourses() {
	if len(s.CoursesEnrolled) == 0 {
		fmt.Println("No course enrolled")
		return
	}
	for _, course := range s.CoursesEnrolled {
		fmt.Println(course.CourseName)
	}
}

func (s *Student) findCourseByID(courseID string) (*Course, bool) {
	for _, course := range s.CoursesEnrolled {
		if course.CourseID == courseID {
			return &course, true
		}
	}
	return nil, false
}

func (s *Student) findCourseIndexByID(courseID string) (int, bool) {
	for i, course := range s.CoursesEnrolled {
		if course.CourseID == courseID {
			return i, true
		}
	}
	return -1, false
}

// get course input
func getCourseInput() (Course, error) {
	fmt.Print("Enter course name: ")
	courseName := strings.TrimSpace(readInput())
	if !isValidCourseName(courseName) {
		return Course{}, errors.New("invalid course name")
	}

	fmt.Print("Enter course ID: ")
	courseID := strings.TrimSpace(readInput())
	if !isValidCourseID(courseID) {
		return Course{}, errors.New("invalid course ID")
	}

	fmt.Print("Enter instructor's name: ")
	instructor := strings.TrimSpace(readInput())
	if !isValidInstructor(instructor) {
		return Course{}, errors.New("invalid instructor name")
	}

	saveCourseToFile(courseID, courseName, instructor)
	return Course{CourseID: courseID, CourseName: courseName, Instructor: instructor}, nil
}

// save course into a file
func saveCourseToFile(courseID, courseName, instructor string) {
	file, err := os.OpenFile("courses.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Ops! failed to write to file: %v\n", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	fileInfo, err := file.Stat()
	if err == nil && fileInfo.Size() == 0 {
		writer.Write([]string{"course_id", "course_name", "instructor"})
	}

	writer.Write([]string{courseID, courseName, instructor})
}

//! validations methods:
func isValidName(name string) bool {
	return regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString(name)
}

func isValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.(com|edu|net)$`
	return regexp.MustCompile(pattern).MatchString(email)
}

func isValidStudentID(studentID string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9/\s]+$`).MatchString(studentID)
}

func isValidGender(gender string) bool {
	for _, validGender := range validGenders {
		if strings.EqualFold(gender, validGender) {
			return true
		}
	}
	return false
}

func isValidCourseName(name string) bool {
	return regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString(name)
}

func isValidCourseID(courseID string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(courseID)
}

func isValidInstructor(name string) bool {
	return regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString(name)
}

// read input
func readInput() string {
	var input string
	fmt.Scanln(&input)
	return input
}