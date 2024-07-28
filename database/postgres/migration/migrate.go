package main

import (
	"github.com/hanshal101/term-test-monitor/database/model"
	"github.com/hanshal101/term-test-monitor/database/postgres"
)

func init() {
	postgres.PostgresInitializer()
}

func main() {
	postgres.DB.AutoMigrate(&model.AllocationResult{})
	postgres.DB.AutoMigrate(&model.StudentsDB{})
	postgres.DB.AutoMigrate(&model.CreateTimeTable{})
	postgres.DB.AutoMigrate(&model.Main_Teachers{})
	postgres.DB.AutoMigrate(&model.Co_Teachers{})
	postgres.DB.AutoMigrate(&model.TeacherAllocation{})
	postgres.DB.AutoMigrate(&model.AttPerSt{})
	postgres.DB.AutoMigrate(&model.AttendanceModel{})
	postgres.DB.AutoMigrate(&model.Attendence_Models{})
	postgres.DB.AutoMigrate(&model.Subject{})
	postgres.DB.AutoMigrate(&model.PaperModel{})
}
