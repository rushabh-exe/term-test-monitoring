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

// package main

// import (
// 	"fmt"
// )

// // Classroom represents a classroom with a number of benches
// type Classroom struct {
// 	index   int
// 	benches int
// }

// // AccommodateStudents checks if all students can be accommodated and prints the seating arrangement
// func AccommodateStudents(x int, y int, classrooms []Classroom) bool {
// 	classASeated := 0
// 	classBSeated := 0
// 	totalBenches := 0

// 	for _, classroom := range classrooms {
// 		if classASeated >= x && classBSeated >= y {
// 			break
// 		}

// 		// Each bench can seat one student from class A and one from class B
// 		classAStudents := min(classroom.benches, x-classASeated)
// 		classBStudents := min(classroom.benches, y-classBSeated)

// 		if classAStudents > 0 {
// 			fmt.Printf("Classroom %d: Class A: %d-%d students, ", classroom.index+1, classASeated+1, classASeated+classAStudents)
// 			classASeated += classAStudents
// 		} else {
// 			fmt.Printf("Classroom %d: Class A: nil students, ", classroom.index+1)
// 		}

// 		if classBStudents > 0 {
// 			fmt.Printf("Class B: %d-%d students\n", classBSeated+1, classBSeated+classBStudents)
// 			classBSeated += classBStudents
// 		} else {
// 			fmt.Printf("Class B: nil students\n")
// 		}

// 		totalBenches += classroom.benches
// 	}

// 	// Check if all students have been seated
// 	return classASeated >= x && classBSeated >= y
// }

// func min(a, b int) int {
// 	if a < b {
// 		return a
// 	}
// 	return b
// }

// func main() {
// 	var x, y, n int
// 	fmt.Print("Enter the number of students in class A: ")
// 	fmt.Scan(&x)
// 	fmt.Print("Enter the number of students in class B: ")
// 	fmt.Scan(&y)
// 	fmt.Print("Enter the number of classrooms: ")
// 	fmt.Scan(&n)

// 	classrooms := make([]Classroom, n)
// 	for i := 0; i < n; i++ {
// 		fmt.Printf("Enter the number of benches in classroom %d: ", i+1)
// 		fmt.Scan(&classrooms[i].benches)
// 		classrooms[i].index = i
// 	}

// 	if AccommodateStudents(x, y, classrooms) {
// 		fmt.Println("True - All students can be accommodated.")
// 	} else {
// 		fmt.Println("False - Not all students can be accommodated.")
// 	}
// }

// package main

// import (
// 	"fmt"

// 	"github.com/hanshal101/term-test-monitor/database/model"
// 	"github.com/hanshal101/term-test-monitor/database/postgres"
// )

// func main() {
// 	postgres.PostgresInitializer()
// 	fmt.Println("HELLO")
// 	var students = []model.StudentsDB{
// 		{Name: "Rajat Nair", RollNo: 1, Email: "rajat.nair@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Tanvi Desai", RollNo: 2, Email: "tanvi.desai@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Aman Mishra", RollNo: 3, Email: "aman.mishra@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Riya Das", RollNo: 4, Email: "riya.das@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Arvind Sen", RollNo: 5, Email: "arvind.sen@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Neha Kapoor", RollNo: 6, Email: "neha.kapoor@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Pranav Reddy", RollNo: 7, Email: "pranav.reddy@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Sanya Iyer", RollNo: 8, Email: "sanya.iyer@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Vikas Kumar", RollNo: 9, Email: "vikas.kumar@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Arjun Bhatt", RollNo: 10, Email: "arjun.bhatt@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Maya Nair", RollNo: 11, Email: "maya.nair@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Karthik Menon", RollNo: 12, Email: "karthik.menon@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Pooja Singh", RollNo: 13, Email: "pooja.singh@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Rohit Mehta", RollNo: 14, Email: "rohit.mehta@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Ishita Sharma", RollNo: 15, Email: "ishita.sharma@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Aditya Rao", RollNo: 16, Email: "aditya.rao@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Nidhi Deshmukh", RollNo: 17, Email: "nidhi.deshmukh@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Vikram Joshi", RollNo: 18, Email: "vikram.joshi@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Ayesha Qureshi", RollNo: 19, Email: "ayesha.qureshi@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Rahul Shah", RollNo: 20, Email: "rahul.shah@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Shruti Agarwal", RollNo: 21, Email: "shruti.agarwal@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Rakesh Patel", RollNo: 22, Email: "rakesh.patel@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Kavya Rao", RollNo: 23, Email: "kavya.rao@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Ishaan Singh", RollNo: 24, Email: "ishaan.singh@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Aditi Sharma", RollNo: 25, Email: "aditi.sharma@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Mohan Gupta", RollNo: 26, Email: "mohan.gupta@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Ananya Menon", RollNo: 27, Email: "ananya.menon@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Vikas Kapoor", RollNo: 28, Email: "vikas.kapoor@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Ritu Jain", RollNo: 29, Email: "ritu.jain@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Rajiv Reddy", RollNo: 30, Email: "rajiv.reddy@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Siddharth Verma", RollNo: 31, Email: "siddharth.verma@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Neeraj Joshi", RollNo: 32, Email: "neeraj.joshi@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Amrita Das", RollNo: 33, Email: "amrita.das@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Amit Patel", RollNo: 34, Email: "amit.patel@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Sneha Reddy", RollNo: 35, Email: "sneha.reddy@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Harsh Sharma", RollNo: 36, Email: "harsh.sharma@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Aarav Gupta", RollNo: 37, Email: "aarav.gupta@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Rekha Shah", RollNo: 38, Email: "rekha.shah@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Krishna Rao", RollNo: 39, Email: "krishna.rao@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Varun Nair", RollNo: 40, Email: "varun.nair@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Priya Kapoor", RollNo: 41, Email: "priya.kapoor@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Sahil Jain", RollNo: 42, Email: "sahil.jain@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Arti Sen", RollNo: 43, Email: "arti.sen@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Nitin Kumar", RollNo: 44, Email: "nitin.kumar@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Surbhi Mehta", RollNo: 45, Email: "surbhi.mehta@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Ashwin Rao", RollNo: 46, Email: "ashwin.rao@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Sneha Desai", RollNo: 47, Email: "sneha.desai@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Nisha Iyer", RollNo: 48, Email: "nisha.iyer@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Vikas Singh", RollNo: 49, Email: "vikas.singh@example.com", Class: "SYB", Department: "EXTC"},
// 		{Name: "Raj Sharma", RollNo: 50, Email: "raj.sharma@example.com", Class: "SYB", Department: "EXTC"},
// 	}

// 	if err := postgres.DB.Create(&students).Error; err != nil {
// 		fmt.Printf("error in students: %v\n", err)
// 		return
// 	}
// 	fmt.Println("DATA SAVED")
// }

// package main

// import (
// 	"encoding/base64"
// 	"encoding/json"
// 	"fmt"
// 	"os"

// 	"github.com/hanshal101/term-test-monitor/database/model"
// )

// func main() {
// 	var teacher model.Main_Teachers
// 	cookie := "eyJJRCI6NCwiQ3JlYXRlZEF0IjoiMjAyNC0wNy0yNFQxNzo1OTowNi4zMTY1OTcrMDU6MzAiLCJVcGRhdGVkQXQiOiIyMDI0LTA3LTI0VDE3OjU5OjA2LjMxNjU5NyswNTozMCIsIkRlbGV0ZWRBdCI6bnVsbCwibmFtZSI6IkFuaXRhIFJhbyIsImVtYWlsIjoibWVodGEuaGFuc2hhbDEwQGdtYWlsLmNvbSIsInBobm8iOiI2NTQzMjEwOTg3In0="

// 	decodedData, err := base64.StdEncoding.DecodeString(cookie)
// 	if err != nil {
// 		fmt.Println("Error decoding base64 data:", err)
// 		return
// 	}

// 	if err := json.Unmarshal([]byte(decodedData), &teacher); err != nil {
// 		fmt.Fprintf(os.Stderr, "Error : %v", err)
// 		return
// 	}

// 	fmt.Println("name :", teacher.Name)
// 	fmt.Println("email:", teacher.Email)
// }

package main

// import (
// 	"log"

// 	"github.com/wneessen/go-mail"
// )

// func main() {
// 	m := mail.NewMsg()
// 	if err := m.From("rishisheshe@outlook.com"); err != nil {
// 		log.Fatalf("failed to set From address: %s", err)
// 	}
// 	if err := m.To("rushabh.mevada@somaiya.edu"); err != nil {
// 		log.Fatalf("failed to set To address: %s", err)
// 	}
// 	m.Subject("This is my first mail with go-mail!")
// 	m.SetBodyString(mail.TypeTextPlain, "Do you like this mail? I certainly do!")
// 	c, err := mail.NewClient("smtp-mail.outlook.com", mail.WithPort(587), mail.WithSMTPAuth(mail.SMTPAuthLogin),
// 		mail.WithUsername("rishisheshe@outlook.com"), mail.WithPassword("rishi@sheshe"))
// 	if err != nil {
// 		log.Fatalf("failed to create mail client: %s", err)
// 	}
// 	if err := c.DialAndSend(m); err != nil {
// 		log.Fatalf("failed to send mail: %s", err)
// 	}
// }
