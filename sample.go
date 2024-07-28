// package main

// import (
// 	"fmt"

// 	"github.com/hanshal101/term-test-monitor/database/model"
// 	"github.com/hanshal101/term-test-monitor/database/postgres"
// )

// func main() {
// 	postgres.PostgresInitializer()
// 	tx := postgres.DB.Begin()
// 	var data = []model.Main_Teachers{
// 		{
// 			Name:  "Arjun Sharma",
// 			Email: "arjun.sharma@example.com",
// 			Phone: "9876543210",
// 		},
// 		{
// 			Name:  "Priya Singh",
// 			Email: "priya.singh@example.com",
// 			Phone: "8765432109",
// 		},
// 		{
// 			Name:  "Vikram Patel",
// 			Email: "vikram.patel@example.com",
// 			Phone: "7654321098",
// 		},
// 		{
// 			Name:  "Anita Rao",
// 			Email: "anita.rao@example.com",
// 			Phone: "6543210987",
// 		},
// 		{
// 			Name:  "Rajesh Nair",
// 			Email: "rajesh.nair@example.com",
// 			Phone: "5432109876",
// 		},
// 		{
// 			Name:  "Sneha Verma",
// 			Email: "sneha.verma@example.com",
// 			Phone: "4321098765",
// 		},
// 		{
// 			Name:  "Amit Kulkarni",
// 			Email: "amit.kulkarni@example.com",
// 			Phone: "3210987654",
// 		},
// 		{
// 			Name:  "Neha Gupta",
// 			Email: "neha.gupta@example.com",
// 			Phone: "2109876543",
// 		},
// 		{
// 			Name:  "Sanjay Iyer",
// 			Email: "sanjay.iyer@example.com",
// 			Phone: "1098765432",
// 		},
// 		{
// 			Name:  "Kavita Desai",
// 			Email: "kavita.desai@example.com",
// 			Phone: "1987654321",
// 		},
// 	}
// 	if err := tx.Create(data).Error; err != nil {
// 		fmt.Println("error in db")
// 	}
// 	var data2 = []model.Co_Teachers{
// 		{
// 			Name:  "Ravi Menon",
// 			Email: "ravi.menon@example.com",
// 			Phone: "9988776655",
// 		},
// 		{
// 			Name:  "Aarti Bhatt",
// 			Email: "aarti.bhatt@example.com",
// 			Phone: "8877665544",
// 		},
// 		{
// 			Name:  "Gopal Sen",
// 			Email: "gopal.sen@example.com",
// 			Phone: "7766554433",
// 		},
// 		{
// 			Name:  "Manisha Kapoor",
// 			Email: "manisha.kapoor@example.com",
// 			Phone: "6655443322",
// 		},
// 		{
// 			Name:  "Suresh Chauhan",
// 			Email: "suresh.chauhan@example.com",
// 			Phone: "5544332211",
// 		},
// 		{
// 			Name:  "Ritu Garg",
// 			Email: "ritu.garg@example.com",
// 			Phone: "4433221100",
// 		},
// 		{
// 			Name:  "Tarun Das",
// 			Email: "tarun.das@example.com",
// 			Phone: "3322110099",
// 		},
// 		{
// 			Name:  "Meena Reddy",
// 			Email: "meena.reddy@example.com",
// 			Phone: "2211009988",
// 		},
// 		{
// 			Name:  "Harish Rao",
// 			Email: "harish.rao@example.com",
// 			Phone: "1100998877",
// 		},
// 		{
// 			Name:  "Pooja Jain",
// 			Email: "pooja.jain@example.com",
// 			Phone: "0099887766",
// 		},
// 	}
// 	if err := tx.Create(data2).Error; err != nil {
// 		fmt.Println("error in db2")
// 	}

// 	tx.Commit()
// }

package main

import (
	"fmt"
)

// Classroom represents a classroom with a number of benches
type Classroom struct {
	index   int
	benches int
}

// AccommodateStudents checks if all students can be accommodated and prints the seating arrangement
func AccommodateStudents(x int, y int, classrooms []Classroom) bool {
	classASeated := 0
	classBSeated := 0
	totalBenches := 0

	for _, classroom := range classrooms {
		if classASeated >= x && classBSeated >= y {
			break
		}

		// Each bench can seat one student from class A and one from class B
		classAStudents := min(classroom.benches, x-classASeated)
		classBStudents := min(classroom.benches, y-classBSeated)

		if classAStudents > 0 {
			fmt.Printf("Classroom %d: Class A: %d-%d students, ", classroom.index+1, classASeated+1, classASeated+classAStudents)
			classASeated += classAStudents
		} else {
			fmt.Printf("Classroom %d: Class A: nil students, ", classroom.index+1)
		}

		if classBStudents > 0 {
			fmt.Printf("Class B: %d-%d students\n", classBSeated+1, classBSeated+classBStudents)
			classBSeated += classBStudents
		} else {
			fmt.Printf("Class B: nil students\n")
		}

		totalBenches += classroom.benches
	}

	// Check if all students have been seated
	return classASeated >= x && classBSeated >= y
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	var x, y, n int
	fmt.Print("Enter the number of students in class A: ")
	fmt.Scan(&x)
	fmt.Print("Enter the number of students in class B: ")
	fmt.Scan(&y)
	fmt.Print("Enter the number of classrooms: ")
	fmt.Scan(&n)

	classrooms := make([]Classroom, n)
	for i := 0; i < n; i++ {
		fmt.Printf("Enter the number of benches in classroom %d: ", i+1)
		fmt.Scan(&classrooms[i].benches)
		classrooms[i].index = i
	}

	if AccommodateStudents(x, y, classrooms) {
		fmt.Println("True - All students can be accommodated.")
	} else {
		fmt.Println("False - Not all students can be accommodated.")
	}
}
