package vitals

import (
	"github.com/gin-gonic/gin"
	"github.com/hanshal101/term-test-monitor/database/model"
)

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

// func CreateAttendance(c *gin.Context) {
// 	var teacher_allocations []model.TeacherAllocation
// 	tx := postgres.DB.Begin()
// 	if err := tx.Find(&teacher_allocations).Error; err != nil {
// 		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
// 		return
// 	}

// 	fmt.Println(len(teacher_allocations))

// 	var teacher_timetable []model.CreateTimeTable
// 	// var result []Result_2
// 	validSuffixes := map[string][]string{
// 		"SY": {"A", "B"},
// 		"TY": {"A", "B"},
// 		"LY": {"A", "B"},
// 	}

// 	var attendance []model.Attendence_Models

// 	for _, t_alloc := range teacher_allocations {
// 		if err := tx.Where("date = ? AND start_time = ? AND end_time = ?", t_alloc.Date, t_alloc.Start_Time, t_alloc.End_Time).Find(&teacher_timetable).Error; err != nil {
// 			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
// 		}
// 		var t_st_alloc []model.AllocationResult

// 		for _, subjects := range teacher_timetable {
// 			if err := tx.Where("class_room = ?", t_alloc.Classroom).Find(&t_st_alloc).Error; err != nil {
// 				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
// 			}
// 			var students []model.StudentsDB
// 			for _, student := range t_st_alloc {
// 				valid := false
// 				for _, suffix := range validSuffixes[subjects.Year] {
// 					if strings.Contains(student.ClassName, subjects.Year+suffix) {
// 						valid = true
// 						break
// 					}
// 				}
// 				if !valid {
// 					continue
// 				}

// 				if err := tx.Where("class = ? AND roll_no BETWEEN ? AND ?", student.ClassName, student.Start, student.End).Find(&students).Error; err != nil {
// 					fmt.Fprintf(os.Stderr, "Error: %v\n", err)
// 				}
// 			}
// 			for _, st := range students {
// 				attper := model.Attendence_Models{
// 					M_Teacher:  t_alloc.Main_Teacher,
// 					C_Teacher:  t_alloc.Co_Teacher,
// 					Classroom:  t_alloc.Classroom,
// 					Date:       t_alloc.Date,
// 					Start_Time: t_alloc.Start_Time,
// 					End_Time:   t_alloc.End_Time,
// 					Subject:    subjects.Subject,
// 					Year:       subjects.Year,
// 					IsPresent:  false,
// 					Supplement: 0,
// 					Name:       st.Name,
// 					RollNo:     st.RollNo,
// 					Class:      st.Class,
// 				}
// 				attendance = append(attendance, attper)
// 			}
// 		}

// 	}
// 	if err := tx.Create(&attendance).Error; err != nil {
// 		fmt.Println("error in creating attendence", err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "error in creating attendance"})
// 		return
// 	}
// 	tx.Commit()

// 	c.JSON(http.StatusOK, gin.H{"success": "db created successfully"})
// }

func CreateAttendance(c *gin.Context) {

}
