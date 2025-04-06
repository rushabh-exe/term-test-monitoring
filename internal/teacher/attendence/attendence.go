package attendence

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hanshal101/term-test-monitor/database/model"
	"github.com/hanshal101/term-test-monitor/database/postgres"
	"github.com/hanshal101/term-test-monitor/helpers/auth"
)

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

type ST_Data struct {
	Subject  string `json:"subject"`
	Year     string `json:"year"`
	Class    string `json:"class"`
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

func Test3(c *gin.Context) {
	teacher, err := auth.GetTeacher(c)
	if err != nil {
		fmt.Println("error in cookie")
		return
	}
	var teacher_allocations []model.TeacherAllocation
	tx := postgres.DB.Begin()
	if err := tx.Where("main_teacher = ?", teacher.Name).Find(&teacher_allocations).Error; err != nil {
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
			var className string
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
				className = student.ClassName
			}
			data := ST_Data{
				Subject:  subjects.Subject,
				Year:     subjects.Year,
				Students: students,
				Class:    className,
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

func CreateAttendence(c *gin.Context) {
	var attendence_req []model.Attendence_Models
	if err := c.BindJSON(&attendence_req); err != nil {
		fmt.Fprintf(os.Stderr, "Error in binding : %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in json format"})
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

	for _, req := range attendence_req {
		// Check if record already exists
		var existingRecord model.Attendence_Models
		conditions := model.Attendence_Models{
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
		}

		result := tx.Where(conditions).First(&existingRecord)
		if result.Error == nil {
			// Record exists, update it
			if err := tx.Model(&existingRecord).Updates(map[string]interface{}{
				"is_present": req.IsPresent,
				"supplement": req.Supplement,
			}).Error; err != nil {
				tx.Rollback()
				fmt.Fprintf(os.Stderr, "Error updating record: %v\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update record"})
				return
			}
		} else {
			// Record doesn't exist, create new
			if err := tx.Create(&req).Error; err != nil {
				tx.Rollback()
				fmt.Fprintf(os.Stderr, "Error creating record: %v\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create record"})
				return
			}
		}
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Attendance records processed successfully"})
}

// func EditAttendance(c *gin.Context) {
// 	var attendanceReq []model.Attendence_Models
// 	if err := c.BindJSON(&attendanceReq); err != nil {
// 		fmt.Fprintf(os.Stderr, "Error in binding: %v\n", err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
// 		return
// 	}

// 	tx := postgres.DB.Begin()
// 	defer func() {
// 		if r := recover(); r != nil {
// 			tx.Rollback()
// 			fmt.Fprintf(os.Stderr, "Error occurred: %v\n", r)
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
// 		}
// 	}()

// 	for _, req := range attendanceReq {
// 		// Define the conditions for updating the record
// 		conditions := &model.Attendence_Models{
// 			M_Teacher:  req.M_Teacher,
// 			C_Teacher:  req.C_Teacher,
// 			Classroom:  req.Classroom,
// 			Class:      req.Class,
// 			Year:       req.Year,
// 			Date:       req.Date,
// 			Start_Time: req.Start_Time,
// 			End_Time:   req.End_Time,
// 			Name:       req.Name,
// 			RollNo:     req.RollNo,
// 			Subject:    req.Subject,
// 			// Exclude IsPresent field from conditions
// 		}

// 		// Update the IsPresent field without fetching the record
// 		if err := tx.Model(&model.Attendence_Models{}).Where(conditions).Update("is_present", req.IsPresent).Error; err != nil {
// 			tx.Rollback()
// 			fmt.Fprintf(os.Stderr, "Error updating record: %v\n", err)
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update record"})
// 			return
// 		}
// 	}

// 	tx.Commit()
// 	c.JSON(http.StatusOK, gin.H{"message": "Attendance records updated successfully"})
// }
