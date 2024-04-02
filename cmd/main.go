package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hanshal101/term-test-monitor/database/postgres"
	"github.com/hanshal101/term-test-monitor/helpers/auth"
	"github.com/hanshal101/term-test-monitor/internal/admin"
	students "github.com/hanshal101/term-test-monitor/internal/admin/students"
	"github.com/hanshal101/term-test-monitor/internal/admin/teachers"
	"github.com/hanshal101/term-test-monitor/internal/admin/vitals"
	"github.com/hanshal101/term-test-monitor/internal/teacher"
	"github.com/hanshal101/term-test-monitor/internal/teacher/attendence"
)

func init() {
	postgres.PostgresInitializer()
}

func main() {
	r := gin.Default()
	r.Use(cors.Default())
	adminGroup := r.Group("/admin")
	{
		adminGroup.GET("/", admin.BaseGET)
		adminGroup.GET("/:year/:subject/:class", students.DashboardAttendence)
		adminGroup.DELETE("/:year/:subject/:class", students.DeleteAttendence)
		adminGroup.PUT("/:year/:subject/:class", students.EditAttendence)
		create := adminGroup.Group("/create")
		{
			create.POST("/timetable", admin.CreateTimeTable)
			student := create.Group("/student")
			{
				student.POST("/dualAllocation", students.DualAllocation)
				student.POST("/singleAllocation", students.SingleAllocation)
			}
			teacher := create.Group("/teacher")
			{
				teacher.POST("/allocation", teachers.CreateTeacherAllocation)
			}
			vital := create.Group("/vitals")
			{
				vital.GET("/", vitals.Base)
				vital.GET("/:year", vitals.GetSubject)
				vital.POST("/:year", vitals.CreateSubject)
				vital.DELETE("/:year/:subject", vitals.DeleteSubject)
			}
		}
		get := adminGroup.Group("/get")
		{
			get.GET("/timetable", admin.GetTT)
			get.GET("/timetable/:Year", admin.GetTTbyYear)
			get.GET("/student/allocation", students.GetAllocation)
			get.GET("/teacher/allocation", teachers.GetTeacherAllocation)
		}
	}

	teacherGroup := r.Group("/teacher")
	{
		teacherGroup.GET("/", teacher.BaseGET)
		teacherGroup.POST("/login", auth.IsTeacher, Generic)
		teacherGroup.GET("/getAttendence", attendence.Test3)
		teacherGroup.POST("/getAttendence", attendence.CreateAttendence)
		teacherGroup.PUT("/getAttendence", attendence.EditAttendance)

	}

	r.Run(":3001")
}

func Generic(c *gin.Context) {
	c.JSON(200, gin.H{"message": c.Request.Method})
}
