package attendence

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hanshal101/term-test-monitor/database/model"
	"github.com/hanshal101/term-test-monitor/database/postgres"
)

// func Test(c *gin.Context) {
// 	session_cookie, err := c.Cookie("lolTest")
// 	if err != nil {
// 		c.JSON(500, gin.H{"error": "Failed to get session cookie"})
// 		return
// 	}
// 	if session_cookie == "" {
// 		c.JSON(400, gin.H{"error": "No session cookie provided"})
// 		return
// 	}

// 	var teacherAlloc []model.TeacherAllocation
// 	var STAll []model.AllocationResult
// 	var Students []model.StudentsDB
// 	var SubjectsTT []model.CreateTimeTable
// 	var result []map[string]interface{}
// 	tx := postgres.DB.Begin()
// 	tx.Where("main_teacher = ?", session_cookie).Find(&teacherAlloc)

// 	fmt.Println(len(teacherAlloc))

// 	for _, arr := range teacherAlloc {
// 		if err := tx.Where("date = ? AND start_time = ? AND end_time = ?", arr.Date, arr.Start_Time, arr.End_Time).Find(&SubjectsTT).Error; err != nil {
// 			errorMsg := fmt.Sprintf("Failed to fetch subjects for date %s, start time %s, end time %s: %s", arr.Date, arr.Start_Time, arr.End_Time, err)
// 			log.Println(errorMsg) // Log the detailed error message
// 			c.JSON(500, gin.H{"error": "Failed to fetch subjects"})
// 			return
// 		}
// 		teacher_classroom := arr.Classroom
// 		if err := tx.Where("class_room = ?", teacher_classroom).Find(&STAll).Error; err != nil {
// 			errorMsg := fmt.Sprintf("Failed to fetch allocation results for classroom %s: %s", teacher_classroom, err)
// 			log.Println(errorMsg) // Log the detailed error message
// 			c.JSON(500, gin.H{"error": "Failed to fetch allocation results"})
// 			return
// 		}

// 		fmt.Println(len(STAll))

// 		for _, arr2 := range STAll {
// 			if err := tx.Where("class = ? AND roll_no BETWEEN ? AND ?", arr2.ClassName, arr2.Start, arr2.End).Find(&Students).Error; err != nil {
// 				errorMsg := fmt.Sprintf("Failed to fetch students for class %s, roll number range %d - %d: %s", arr2.ClassName, arr2.Start, arr2.End, err)
// 				log.Println(errorMsg) // Log the detailed error message
// 				c.JSON(500, gin.H{"error": "Failed to fetch students"})
// 				return
// 			}

//				data := map[string]interface{}{
//					"date":         arr.Date,
//					"start_time":   arr.Start_Time,
//					"end_time":     arr.End_Time,
//					"classroom":    arr.Classroom,
//					"subject":      SubjectsTT,
//					"allocation":   arr2,
//					"students":     Students,
//					"main_teacher": session_cookie,
//				}
//				result = append(result, data)
//			}
//		}
//		tx.Commit()
//		c.JSON(200, result)
//	}
type Result struct {
	Name   string
	Rollno int
}

type Out struct {
	ExamDetails []Exam_subject
	Class_Name  string
	Class_Room  string
	Students    []Result
}

type TeacherOut struct {
	Main_Teacher string
	Co_Teacher   string
	Date         string
	Start_Time   string
	End_Time     string
	Out          []Out
}

type Exam_subject struct {
	Name string
	Year string
}

type SubjectOut struct {
	Date    string
	Teacher []TeacherOut
}

// func Test(c *gin.Context) {
// 	var subject_teacher []model.CreateTimeTable
// 	var subjectCount []model.CreateTimeTable
// 	year := "SY"
// 	if err := postgres.DB.Where("year = ?", year).Find(&subjectCount).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch main teachers"})
// 		return
// 	}
// 	print(len(subjectCount), "\n")

// 	var syaStudents []model.StudentsDB
// 	var subjectOut []SubjectOut

// 	for _, perSubAlloc := range subjectCount {
// 		// var TeacherResult []TeacherOut
// 		match_date := perSubAlloc.Date
// 		match_st_time := perSubAlloc.Start_Time
// 		match_ed_time := perSubAlloc.End_Time
// 		var TeacherAlloc []model.TeacherAllocation
// 		if err := postgres.DB.Where("date = ? AND start_time = ? AND end_time = ?", match_date, match_st_time, match_ed_time).Find(&TeacherAlloc).Error; err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch main teachers"})
// 			return
// 		}
// 		if err := postgres.DB.Where("date = ? AND start_time = ?", match_date, match_st_time).Find(&subject_teacher); err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch main teachers"})
// 			return
// 		}

// 		for _, teacherAlloc := range TeacherAlloc {
// 			var result []Out
// 			classRoom := teacherAlloc.Classroom
// 			var stAlloc []model.AllocationResult
// 			if err := postgres.DB.Where("class_room = ?", classRoom).Find(&stAlloc).Error; err != nil {
// 				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch main teachers"})
// 				return
// 			}
// 			for _, allocations := range stAlloc {

// 				if err := postgres.DB.Where("class = ? AND roll_no BETWEEN ? AND ?", allocations.ClassName, allocations.Start, allocations.End).Find(&syaStudents).Error; err != nil {
// 					c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch %s students", allocations.ClassName)})
// 					return
// 				}
// 				var studentNames []Result
// 				for _, student := range syaStudents {
// 					studentNames = append(studentNames, Result{student.Name, student.RollNo})
// 				}
// 				result = append(result, Out{Class_Name: allocations.ClassName, Class_Room: allocations.ClassRoom, Students: studentNames})
// 			}
// 			// TeacherResult = append(TeacherResult, TeacherOut{Main_Teacher: teacherAlloc.Main_Teacher, Co_Teacher: teacherAlloc.Co_Teacher, Start_Time: teacherAlloc.Start_Time, End_Time: teacherAlloc.End_Time, Out: result, Date: teacherAlloc.Date})
// 		}
// 		// subjectOut = append(subjectOut, SubjectOut{Subject: perSubAlloc.Subject, Year: perSubAlloc.Year, Date: perSubAlloc.Date, Teacher: TeacherResult})
// 	}
// 	c.JSON(http.StatusOK, subjectOut)
// }

type ST_Data struct {
	Subject  string `json:"subject"`
	Year     string `json:"year"`
	Students []model.StudentsDB
}

type Result_2 struct {
	M_Teacher   string    `json:"m_teacher"`
	C_Teacher   string    `json:"c_teacher"`
	Classroom   string    `json:"class_room"`
	Date        string    `json:"date"`
	Start_Time  string    `json:"start_time"`
	End_Time    string    `json:"end_time"`
	StudentData []ST_Data `json:"student_data"`
}

// func Test2(c *gin.Context) {
// 	teacher_name := "John Doe"
// 	var student_data []ST_Data
// 	var teacher_allocations []model.TeacherAllocation
// 	tx := postgres.DB.Begin()
// 	if err := tx.Where("main_teacher = ?", teacher_name).Find(&teacher_allocations).Error; err != nil {
// 		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
// 	}
// 	print(len(teacher_allocations), "\n")

// 	var teacher_timetable []model.CreateTimeTable
// 	for _, allocations := range teacher_allocations {
// 		if err := tx.Where("date = ? AND start_time = ? AND end_time = ?", allocations.Date, allocations.Start_Time, allocations.End_Time).Find(&teacher_timetable).Error; err != nil {
// 			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
// 		}
// 	}
// 	print(len(teacher_timetable), "\n")

// 	var teacher_st_allocation []model.AllocationResult
// 	for _, st_allocation := range teacher_allocations {
// 		if err := tx.Where("class_room = ?", st_allocation.Classroom).Find(&teacher_st_allocation).Error; err != nil {
// 			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
// 		}
// 	}
// 	print(len(teacher_st_allocation), "\n")

// 	for _, alloc := range teacher_st_allocation {
// 		var value map[string]string
// 		value["SYA"] = "SY"
// 		value["SYB"] = "SY"
// 		value["TYA"] = "TY"
// 		value["TYB"] = "TY"
// 		value["LYA"] = "LY"
// 		value["LYB"] = "LY"

// 		var students []model.StudentsDB
// 		if err := tx.Where("class = ?", alloc.ClassName).Find(&students).Error; err != nil {
// 			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
// 		}
// 		var data ST_Data
// 		data = ST_Data{
// 			// Subject:  subject.Subject,
// 			// Year:     subject.Year,
// 			Students: students,
// 		}
// 		student_data = append(student_data, data)
// 	}

// 	tx.Commit()

// }

func Test3(c *gin.Context) {
	teacher_name := "Test"
	var teacher_allocations []model.TeacherAllocation
	tx := postgres.DB.Begin()
	if err := tx.Where("main_teacher = ?", teacher_name).Find(&teacher_allocations).Error; err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	fmt.Println(len(teacher_allocations))

	var teacher_timetable []model.CreateTimeTable
	var result []Result_2
	validSuffixes := map[string][]string{
		"SY": {"A", "B"},
		"TY": {"A", "B"},
		"LY": {"A", "B"},
	}

	for _, t_alloc := range teacher_allocations {
		if err := tx.Where("date = ? AND start_time = ? AND end_time = ?", t_alloc.Date, t_alloc.Start_Time, t_alloc.End_Time).Find(&teacher_timetable).Error; err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		var t_st_alloc []model.AllocationResult
		var st_data []ST_Data

		for _, subjects := range teacher_timetable {
			if err := tx.Where("class_room = ?", t_alloc.Classroom).Find(&t_st_alloc).Error; err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			}
			var students []model.StudentsDB
			for _, student := range t_st_alloc {
				valid := false
				for _, suffix := range validSuffixes[subjects.Year] {
					if strings.Contains(student.ClassName, subjects.Year+suffix) {
						valid = true
						break
					}
				}
				if !valid {
					continue
				}

				if err := tx.Where("class = ? AND roll_no BETWEEN ? AND ?", student.ClassName, student.Start, student.End).Find(&students).Error; err != nil {
					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				}
			}
			data := ST_Data{
				Subject:  subjects.Subject,
				Year:     subjects.Year,
				Students: students,
			}
			st_data = append(st_data, data)
		}
		data := Result_2{
			M_Teacher:   t_alloc.Main_Teacher,
			C_Teacher:   t_alloc.Co_Teacher,
			Classroom:   t_alloc.Classroom,
			Date:        t_alloc.Date,
			Start_Time:  t_alloc.Start_Time,
			End_Time:    t_alloc.End_Time,
			StudentData: st_data,
		}
		result = append(result, data)
	}
	tx.Commit()
	c.JSON(200, result)
}

func SendAttendence(c *gin.Context) {

}

type Attendence_Models struct {
	M_Teacher  string `json:"m_teacher"`
	C_Teacher  string `json:"c_teacher"`
	Classroom  string `json:"class_room"`
	Date       string `json:"date"`
	Start_Time string `json:"start_time"`
	End_Time   string `json:"end_time"`
	Subject    string `json:"subject"`
	Year       string `json:"year"`
	Name       string `json:"name"`
	RollNo     int    `json:"roll_no"`
	Class      string `json:"class"`
	IsPresent  bool   `json:"is_present"`
}

func CreateAttendence(c *gin.Context) {
	var attendence_req []model.Attendence_Models
	if err := c.BindJSON(&attendence_req); err != nil {
		fmt.Fprintf(os.Stderr, "Error in binding :\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in json format"})
	}
	tx := postgres.DB.Begin()
	if err := tx.Create(&attendence_req).Error; err != nil {
		fmt.Fprintf(os.Stderr, "Error in binding :\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in creating attendence"})
	}
	tx.Commit()
	c.JSON(http.StatusCreated, gin.H{"error": "attendence created"})
}

func EditAttendance(c *gin.Context) {
	var attendanceReq []model.Attendence_Models
	if err := c.BindJSON(&attendanceReq); err != nil {
		fmt.Fprintf(os.Stderr, "Error in binding: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	tx := postgres.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			fmt.Fprintf(os.Stderr, "Error occurred: %v\n", r)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
	}()

	for _, req := range attendanceReq {
		// Define the conditions for updating the record
		conditions := &model.Attendence_Models{
			M_Teacher:  req.M_Teacher,
			C_Teacher:  req.C_Teacher,
			Classroom:  req.Classroom,
			Class:      req.Class,
			Year:       req.Year,
			Date:       req.Date,
			Start_Time: req.Start_Time,
			End_Time:   req.End_Time,
			Name:       req.Name,
			RollNo:     req.RollNo,
			Subject:    req.Subject,
			// Exclude IsPresent field from conditions
		}

		// Update the IsPresent field without fetching the record
		if err := tx.Model(&model.Attendence_Models{}).Where(conditions).Update("is_present", req.IsPresent).Error; err != nil {
			tx.Rollback()
			fmt.Fprintf(os.Stderr, "Error updating record: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update record"})
			return
		}
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Attendance records updated successfully"})
}
