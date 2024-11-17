package teachers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hanshal101/term-test-monitor/database/model"
	"github.com/hanshal101/term-test-monitor/database/postgres"
	alloc_helper "github.com/hanshal101/term-test-monitor/helpers/alloc_helpers"
	mailClient "github.com/hanshal101/term-test-monitor/internal/mail"
	"golang.org/x/exp/rand"
)

// func CreateTeacherAllocation(c *gin.Context) {
// 	fetch, err := http.Get("http://localhost:3002/api/allotment")
// 	if err != nil {
// 		log.Fatalf("Error in getting Teacher Allocation: %v", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch teacher allocation"})
// 		return
// 	}
// 	defer fetch.Body.Close()

// 	var result []model.TeacherAllocation
// 	if err := json.NewDecoder(fetch.Body).Decode(&result); err != nil {
// 		log.Fatalf("Error decoding response body: %v", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode response body"})
// 		return
// 	}

// 	tx := postgres.DB.Begin()
// 	if err := tx.Where("deleted_at IS NOT NULL").Delete(&model.TeacherAllocation{}).Error; err != nil {
// 		log.Fatalf("Error in DB Deletion: %v", err)
// 	}

// 	for _, item := range result {
// 		teacherAllocation := model.TeacherAllocation{
// 			Date:         item.Date,
// 			Start_Time:   item.Start_Time,
// 			End_Time:     item.End_Time,
// 			Classroom:    item.Classroom,
// 			Main_Teacher: item.Main_Teacher,
// 			Co_Teacher:   item.Co_Teacher,
// 		}
// 		if err := tx.Create(&teacherAllocation).Error; err != nil {
// 			tx.Rollback()
// 			log.Fatalf("Error in DB Creation: %v", err)
// 		}
// 	}
// 	tx.Commit()

// 	// Process result further as needed...

// 	c.JSON(http.StatusOK, gin.H{"status": "Allocation Done"})
// }

func CreateTeacherAllocation(c *gin.Context) {
	tx := postgres.DB.Begin()
	classrooms := alloc_helper.GetClass()
	for _, c := range classrooms {
		fmt.Println(c.ClassRoom)
	}

	var mainTeachers []model.Main_Teachers
	var coTeachers []model.Co_Teachers

	if err := tx.Where("deleted_at IS NULL").Find(&mainTeachers).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in finding main teachers"})
		return
	}
	if err := tx.Where("deleted_at IS NULL").Find(&coTeachers).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in finding co teachers"})
		return
	}

	if len(mainTeachers) == 0 || len(coTeachers) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient number of teachers or subteachers in the database"})
		return
	}

	var results []model.TeacherAllocation
	var mtmap = make(map[string]int)
	var ctmap = make(map[string]int)
	for _, teacher := range mainTeachers {
		mtmap[teacher.Email] = 3
	}
	for _, teacher := range coTeachers {
		ctmap[teacher.Email] = 5
	}

	var subjects []model.CreateTimeTable
	if err := tx.Select("DISTINCT date, start_time, end_time").Find(&subjects).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in getting subjects"})
		return
	}

	// Maintain schedules for main and co-teachers
	var mtSchedule = make(map[string][]model.TeacherAllocation)
	var ctSchedule = make(map[string][]model.TeacherAllocation)

	for _, sub := range subjects {
		for _, classroom := range classrooms {
			fmt.Println("Getting teachers")
			mTeacher, ok := SelectAvailableMTeacher(mainTeachers, mtmap, mtSchedule, sub)
			if !ok {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"error": "error in getting main teacher"})
				return
			}
			fmt.Println("MT", mTeacher.Name)
			cTeacher, ok := SelectAvailableCTeacher(coTeachers, ctmap, ctSchedule, sub)
			if !ok {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"error": "error in getting co teacher"})
				return
			}
			fmt.Println("CT", cTeacher.Name)

			allocation := model.TeacherAllocation{
				Classroom:    classroom.ClassRoom,
				Date:         sub.Date,
				Start_Time:   sub.Start_Time,
				End_Time:     sub.End_Time,
				Main_Teacher: mTeacher.Name,
				Co_Teacher:   cTeacher.Name,
			}

			// Update the teacher schedules
			mtSchedule[mTeacher.Email] = append(mtSchedule[mTeacher.Email], allocation)
			ctSchedule[cTeacher.Email] = append(ctSchedule[cTeacher.Email], allocation)

			results = append(results, allocation)
		}
	}
	fmt.Println("Creating allocations")
	if err := tx.Create(&results).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in creating allocations"})
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, results)
}

func GetTeacherAllocation(c *gin.Context) {
	var teacherAllocations []model.TeacherAllocation
	if err := postgres.DB.Where("deleted_at IS NULL").Find(&teacherAllocations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch teacher allocations"})
		return
	}

	c.JSON(http.StatusOK, teacherAllocations)
}

func DeleteTeacherAllocation(c *gin.Context) {
	id := c.Param("id")

	var data model.TeacherAllocation

	tx := postgres.DB.Begin()
	if err := tx.Where("id = ?", id).Delete(&data).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in deleting DB"})
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"error": "Allocation Deleted successfully"})
}

func CheckMultiples() bool {
	return false
}

// Create a global random number generator with a seed
// var rng = rand.New(rand.NewSource(uint64(time.Now().UnixNano())))

func SelectAvailableMTeacher(t []model.Main_Teachers, hm map[string]int, schedule map[string][]model.TeacherAllocation, sub model.CreateTimeTable) (model.Main_Teachers, bool) {
	if len(t) < 1 {
		return model.Main_Teachers{}, false
	}

	for i := 0; i < len(t); i++ {
		teacher := t[getRandomNumber(0, len(t)-1)]
		email := teacher.Email

		if count, exists := hm[email]; exists {
			if count > 0 {
				// Check if the main teacher is available during the required time slot
				available := true
				for _, alloc := range schedule[email] {
					if alloc.Date == sub.Date && alloc.Start_Time == sub.Start_Time && alloc.End_Time == sub.End_Time {
						available = false
						break
					}
				}
				if available {
					hm[email]--
					if hm[email] == 0 {
						delete(hm, email)
					}
					return teacher, true
				}
			}
		}
	}

	return SelectAvailableMTeacher(t, hm, schedule, sub)
}

func SelectAvailableCTeacher(t []model.Co_Teachers, hm map[string]int, schedule map[string][]model.TeacherAllocation, sub model.CreateTimeTable) (model.Co_Teachers, bool) {
	if len(t) < 1 {
		return model.Co_Teachers{}, false
	}

	for i := 0; i < len(t); i++ {
		teacher := t[getRandomNumber(0, len(t)-1)]
		email := teacher.Email

		if count, exists := hm[email]; exists {
			if count > 0 {
				// Check if the co-teacher is available during the required time slot
				available := true
				for _, alloc := range schedule[email] {
					if alloc.Date == sub.Date && alloc.Start_Time == sub.Start_Time && alloc.End_Time == sub.End_Time {
						available = false
						break
					}
				}
				if available {
					hm[email]--
					if hm[email] == 0 {
						delete(hm, email)
					}
					return teacher, true
				}
			}
		}
	}

	return SelectAvailableCTeacher(t, hm, schedule, sub)
}

func getRandomNumber(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func SendMail(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file upload error"})
		return
	}

	filePath := "./uploads/" + file.Filename
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not save file"})
		return
	}
	defer os.Remove(filePath)

	var teacherAllocations []model.TeacherAllocation
	tx := postgres.DB.Begin()
	if err := tx.Find(&teacherAllocations).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in receiving teacherAllocations"})
		return
	}

	teacherMap := make(map[string][]model.TeacherAllocation)
	for _, alloc := range teacherAllocations {
		teacherMap[alloc.Main_Teacher] = append(teacherMap[alloc.Main_Teacher], alloc)
	}

	for teacherName, allocations := range teacherMap {
		var teacher model.Main_Teachers
		if err := tx.Where("name = ?", teacherName).First(&teacher).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "error in receiving teacher"})
			return
		}

		subject := fmt.Sprintf("Assignment Allocations for %s", teacher.Name)
		body := fmt.Sprintf("Dear %s,\n\nYou have been allocated the following assignments:\n", teacher.Name)

		for _, alloc := range allocations {
			allocationDetails := fmt.Sprintf("Classroom: %s, Date: %s, Start Time: %s, End Time: %s",
				alloc.Classroom, alloc.Date, alloc.Start_Time, alloc.End_Time)
			body += allocationDetails + "\n"
		}
		body += "Best regards,\nYour School Admin"

		// Send the email with the uploaded file as an attachment
		mailClient.MailClient(teacher.Email, subject, body, filePath)
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"status": "emails sent successfully"})
}

func EditAllocation(c *gin.Context) {
	allocID := c.Param("")
	var req model.TeacherAllocation
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in binding"})
		return
	}

	var existAlloc model.TeacherAllocation
	if err := postgres.DB.Where("id = ?", allocID).First(&existAlloc).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error in finding the allocation"})
		return
	}

	existAlloc.Main_Teacher = req.Main_Teacher
	existAlloc.Co_Teacher = req.Co_Teacher

	if err := postgres.DB.Save(existAlloc).Error; err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "error in updating"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "data saved"})
}
