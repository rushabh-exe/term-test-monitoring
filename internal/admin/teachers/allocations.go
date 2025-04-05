package teachers

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hanshal101/term-test-monitor/database/model"
	"github.com/hanshal101/term-test-monitor/database/postgres"
	alloc_helper "github.com/hanshal101/term-test-monitor/helpers/alloc_helpers"
	mailClient "github.com/hanshal101/term-test-monitor/internal/mail"
	"gorm.io/gorm"
)

func CreateTeacherAllocation(c *gin.Context) {
	tx := postgres.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if err := tx.Error; err != nil {

			fmt.Println("Transaction rolled back due to error:", err)
		}
	}()

	classrooms := alloc_helper.GetClass()
	if len(classrooms) == 0 {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "No classrooms found"})
		return
	}
	fmt.Println("Found classrooms:", len(classrooms))
	var mainTeachers []model.Main_Teachers
	var coTeachers []model.Co_Teachers

	if err := tx.Where("deleted_at IS NULL").Find(&mainTeachers).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch main teachers"})
		return
	}
	if err := tx.Where("deleted_at IS NULL").Find(&coTeachers).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch co teachers"})
		return
	}

	var allocationLimits []model.AllocationCount
	if err := tx.Find(&allocationLimits).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch allocation limits"})
		return
	}

	mtLimit := 0
	ctLimit := 0
	for _, limit := range allocationLimits {
		countInt, err := strconv.Atoi(limit.Count)
		if err != nil {
			tx.Rollback()
			fmt.Printf("Invalid count value '%s' for type '%s': %v\n", limit.Count, limit.Type, err)
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid count '%s' for type '%s'", limit.Count, limit.Type)})
			return
		}
		if limit.Type == "teachingstaff" {
			mtLimit = countInt
		} else if limit.Type == "nonteachingstaff" {
			ctLimit = countInt

		}
	}

	if mtLimit <= 0 || ctLimit <= 0 {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Allocation limits not set or invalid"})
		return
	}
	fmt.Printf("Allocation limits - Main: %d, Co: %d\n", mtLimit, ctLimit)

	if len(mainTeachers) == 0 {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "No main teachers found in the database"})
		return
	}
	if len(coTeachers) == 0 {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "No co-teachers found in the database"})
		return
	}
	fmt.Printf("Found teachers - Main: %d, Co: %d\n", len(mainTeachers), len(coTeachers))

	var mtmap = make(map[string]int)
	var ctmap = make(map[string]int)
	for _, teacher := range mainTeachers {
		mtmap[teacher.Email] = mtLimit
	}
	for _, teacher := range coTeachers {
		ctmap[teacher.Email] = ctLimit
	}

	var timeSlots []model.CreateTimeTable
	type TimeSlot struct {
		Date      string
		StartTime string `gorm:"column:start_time"`
		EndTime   string `gorm:"column:end_time"`
	}
	if err := tx.Model(&model.CreateTimeTable{}).Distinct("date", "start_time", "end_time").Find(&timeSlots).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch time slots from timetable"})
		return
	}
	if len(timeSlots) == 0 {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "No time slots found in the timetable"})
		return
	}
	fmt.Printf("Found %d distinct time slots\n", len(timeSlots))

	var results []model.TeacherAllocation
	var mtSchedule = make(map[string]map[string]bool)
	var ctSchedule = make(map[string]map[string]bool)

	for _, slot := range timeSlots {
		dateTimeSlotKey := fmt.Sprintf("%s_%s_%s", slot.Date, slot.Start_Time, slot.End_Time)

		shuffledMainTeachers := make([]model.Main_Teachers, len(mainTeachers))
		copy(shuffledMainTeachers, mainTeachers)
		rand.Shuffle(len(shuffledMainTeachers), func(i, j int) {
			shuffledMainTeachers[i], shuffledMainTeachers[j] = shuffledMainTeachers[j], shuffledMainTeachers[i]
		})

		shuffledCoTeachers := make([]model.Co_Teachers, len(coTeachers))
		copy(shuffledCoTeachers, coTeachers)
		rand.Shuffle(len(shuffledCoTeachers), func(i, j int) {
			shuffledCoTeachers[i], shuffledCoTeachers[j] = shuffledCoTeachers[j], shuffledCoTeachers[i]
		})

		assignedMT := make(map[string]bool)
		assignedCT := make(map[string]bool)

		fmt.Printf("\nProcessing Slot: %s\n", dateTimeSlotKey)
		for _, classroom := range classrooms {
			fmt.Printf("  Assigning for Classroom: %s\n", classroom.ClassRoom)

			mTeacher, ok := SelectAvailableMTeacher(shuffledMainTeachers, mtmap, mtSchedule, dateTimeSlotKey, assignedMT)
			if !ok {
				tx.Rollback()
				errorMsg := fmt.Sprintf("Could not find available Main Teacher for slot %s, classroom %s. Insufficient teachers or all available are busy/used up.", dateTimeSlotKey, classroom.ClassRoom)
				fmt.Println("ERROR:", errorMsg)
				c.JSON(http.StatusConflict, gin.H{"error": errorMsg})
				return
			}
			fmt.Printf("    -> MT: %s (Email: %s, Remaining: %d)\n", mTeacher.Name, mTeacher.Email, mtmap[mTeacher.Email])
			assignedMT[mTeacher.Email] = true

			cTeacher, ok := SelectAvailableCTeacher(shuffledCoTeachers, ctmap, ctSchedule, dateTimeSlotKey, assignedCT, mTeacher.Email) // Pass main teacher email to avoid conflict if needed
			if !ok {
				tx.Rollback()
				errorMsg := fmt.Sprintf("Could not find available Co-Teacher for slot %s, classroom %s. Insufficient teachers or all available are busy/used up.", dateTimeSlotKey, classroom.ClassRoom)
				fmt.Println("ERROR:", errorMsg)
				c.JSON(http.StatusConflict, gin.H{"error": errorMsg})
				return
			}
			fmt.Printf("    -> CT: %s (Email: %s, Remaining: %d)\n", cTeacher.Name, cTeacher.Email, ctmap[cTeacher.Email])
			assignedCT[cTeacher.Email] = true

			allocation := model.TeacherAllocation{
				Classroom:    classroom.ClassRoom,
				Date:         slot.Date,
				Start_Time:   slot.Start_Time,
				End_Time:     slot.End_Time,
				Main_Teacher: mTeacher.Name,
				Co_Teacher:   cTeacher.Name,
			}

			if _, exists := mtSchedule[mTeacher.Email]; !exists {
				mtSchedule[mTeacher.Email] = make(map[string]bool)
			}
			mtSchedule[mTeacher.Email][dateTimeSlotKey] = true

			if _, exists := ctSchedule[cTeacher.Email]; !exists {
				ctSchedule[cTeacher.Email] = make(map[string]bool)
			}
			ctSchedule[cTeacher.Email][dateTimeSlotKey] = true

			results = append(results, allocation)
		}
	}

	fmt.Println("\nCreating allocations in DB...")
	if len(results) == 0 {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"message": "No allocations needed or generated."})
		return
	}

	if err := tx.Create(&results).Error; err != nil {
		tx.Rollback()
		fmt.Println("ERROR creating allocations in DB:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating allocations"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		fmt.Println("ERROR committing transaction:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database transaction commit error"})
		return
	}

	fmt.Println("Allocation process completed successfully.")
	c.JSON(http.StatusOK, gin.H{"status": "Allocation successful", "allocations_created": len(results)}) // Return status and count
}

func SelectAvailableMTeacher(
	shuffledTeachers []model.Main_Teachers,
	allocCounts map[string]int,
	schedule map[string]map[string]bool,
	dateTimeSlotKey string,
	alreadyAssigned map[string]bool,
) (model.Main_Teachers, bool) {

	for _, teacher := range shuffledTeachers {
		email := teacher.Email

		if alreadyAssigned[email] {
			continue
		}

		count, exists := allocCounts[email]
		if !exists || count <= 0 {
			continue
		}

		if teacherSchedule, scheduleExists := schedule[email]; scheduleExists {
			if teacherSchedule[dateTimeSlotKey] {
				continue
			}
		}

		allocCounts[email]--
		return teacher, true
	}

	return model.Main_Teachers{}, false
}

func SelectAvailableCTeacher(
	shuffledTeachers []model.Co_Teachers,
	allocCounts map[string]int,
	schedule map[string]map[string]bool,
	dateTimeSlotKey string,
	alreadyAssigned map[string]bool,
	assignedMainTeacherEmail string,
) (model.Co_Teachers, bool) {

	for _, teacher := range shuffledTeachers {
		email := teacher.Email

		if email == assignedMainTeacherEmail {
			continue
		}

		if alreadyAssigned[email] {
			continue
		}

		count, exists := allocCounts[email]
		if !exists || count <= 0 {
			continue
		}

		if teacherSchedule, scheduleExists := schedule[email]; scheduleExists {
			if teacherSchedule[dateTimeSlotKey] {
				continue
			}
		}

		allocCounts[email]--
		return teacher, true
	}

	return model.Co_Teachers{}, false
}

func GetTeacherAllocation(c *gin.Context) {
	var teacherAllocations []model.TeacherAllocation
	if err := postgres.DB.Where("deleted_at IS NULL").Order("date, start_time, classroom").Find(&teacherAllocations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch teacher allocations"})
		return
	}
	if teacherAllocations == nil {
		teacherAllocations = []model.TeacherAllocation{}
	}
	c.JSON(http.StatusOK, teacherAllocations)
}

func DeleteTeacherAllocation(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing allocation ID"})
		return
	}

	var data model.TeacherAllocation
	if err := postgres.DB.First(&data, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Allocation not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking allocation"})
		}
		return
	}

	if err := postgres.DB.Delete(&model.TeacherAllocation{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting allocation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "Allocation deleted successfully"})
}

func SendMail(c *gin.Context) {
	file, err := c.FormFile("file")
	isAttachment := true
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			fmt.Println("No file uploaded for email attachment.")
			isAttachment = false
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File upload error: " + err.Error()})
			return
		}
	}

	filePath := ""
	if isAttachment {
		filePath = "./uploads/" + file.Filename
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save uploaded file"})
			return
		}
		defer func() {
			if filePath != "" {
				if err := os.Remove(filePath); err != nil {
					fmt.Printf("Warning: could not remove temporary upload file %s: %v\n", filePath, err)
				}
			}
		}()
	}

	var teacherAllocations []model.TeacherAllocation
	if err := postgres.DB.Where("deleted_at IS NULL").Find(&teacherAllocations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching teacher allocations"})
		return
	}

	if len(teacherAllocations) == 0 {
		c.JSON(http.StatusOK, gin.H{"status": "No allocations found to send emails for."})
		return
	}

	teacherMap := make(map[string][]model.TeacherAllocation)
	mainTeacherEmails := make(map[string]string)

	for _, alloc := range teacherAllocations {
		var teacher model.Main_Teachers
		if err := postgres.DB.Where("name = ? AND deleted_at IS NULL", alloc.Main_Teacher).First(&teacher).Error; err != nil {

			fmt.Printf("Warning: Could not find teacher details for %s: %v\n", alloc.Main_Teacher, err)
			continue
		}
		teacherMap[teacher.Email] = append(teacherMap[teacher.Email], alloc)
		mainTeacherEmails[teacher.Email] = teacher.Name
	}

	fmt.Printf("Sending emails to %d teachers...\n", len(teacherMap))
	for email, allocations := range teacherMap {
		teacherName := mainTeacherEmails[email]

		subject := fmt.Sprintf("Exam Duty Allocations for %s", teacherName)
		body := fmt.Sprintf("Dear %s,\n\nYou have been allocated the following exam duties:\n\n", teacherName)

		for i, alloc := range allocations {
			allocationDetails := fmt.Sprintf("%d. Classroom: %s, Date: %s, Time: %s - %s",
				i+1, alloc.Classroom, alloc.Date, alloc.Start_Time, alloc.End_Time)
			body += allocationDetails + "\n"
		}
		body += "\nPlease find attached document for further details (if provided).\n\nBest regards,\nExam Coordinator" // Adjust signature

		fmt.Printf("Sending mail to %s (%s)\n", teacherName, email)

		attachmentPath := ""
		if isAttachment {
			attachmentPath = filePath
		}
		mailClient.MailClient(email, subject, body, attachmentPath)
		if err != nil {
			fmt.Printf("ERROR sending mail to %s (%s): %v\n", teacherName, email, err)
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "Emails dispatch process completed."})
}

func EditAllocation(c *gin.Context) {
	allocID := c.Param("allocID")
	if allocID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing allocation ID"})
		return
	}

	var req model.TeacherAllocation
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	if req.Main_Teacher == "" || req.Co_Teacher == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Main Teacher and Co-Teacher names cannot be empty"})
		return
	}

	var existAlloc model.TeacherAllocation
	tx := postgres.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if err := tx.Where("id = ? AND deleted_at IS NULL", allocID).First(&existAlloc).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Allocation not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error finding the allocation"})
		}
		return
	}

	var tempTeacher model.Main_Teachers
	if err := tx.Where("name = ? AND deleted_at IS NULL", req.Main_Teacher).First(&tempTeacher).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Main Teacher '%s' not found or inactive", req.Main_Teacher)})
		return
	}
	var tempCoTeacher model.Co_Teachers
	if err := tx.Where("name = ? AND deleted_at IS NULL", req.Co_Teacher).First(&tempCoTeacher).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Co-Teacher '%s' not found or inactive", req.Co_Teacher)})
		return
	}
	existAlloc.Main_Teacher = req.Main_Teacher
	existAlloc.Co_Teacher = req.Co_Teacher

	if err := tx.Save(&existAlloc).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating allocation"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving updated allocation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "Allocation updated successfully"})
}
