package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hanshal101/term-test-monitor/database/postgres"
	alloc_helper "github.com/hanshal101/term-test-monitor/helpers/alloc_helpers"
	"github.com/hanshal101/term-test-monitor/helpers/auth"
	"github.com/hanshal101/term-test-monitor/internal/admin"
	students "github.com/hanshal101/term-test-monitor/internal/admin/students"
	"github.com/hanshal101/term-test-monitor/internal/admin/teachers"
	"github.com/hanshal101/term-test-monitor/internal/admin/vitals"
	"github.com/hanshal101/term-test-monitor/internal/teacher"
	"github.com/hanshal101/term-test-monitor/internal/teacher/attendence"
	middleware "github.com/hanshal101/term-test-monitor/middleware/auth"
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
			create.POST("/timetable/:year", admin.CreateTimeTable)
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
				vital.GET("/teachers/:type", vitals.GetTeachers)
				vital.POST("/createTeacher", vitals.CreateTeacher)
				vital.PUT("/teachers/:type", vitals.EditTeacher)
				vital.DELETE("/teachers/:type/:email", vitals.DeleteTeacher)
			}
		}
		get := adminGroup.Group("/get")
		{
			get.GET("/timetable", admin.GetTT)
			get.GET("/timetable/:Year", admin.GetTTbyYear)
			get.DELETE("/timetable/:year", admin.DeleteTimeTable)
			get.GET("/student/allocation", students.GetAllocation)
			get.DELETE("/student/allocation/:id", students.DeleteAllocation)
			get.GET("/teacher/allocation", teachers.GetTeacherAllocation)
			get.DELETE("/teacher/allocation/:id", teachers.DeleteTeacherAllocation)
		}
	}

	teacherGroup := r.Group("/teacher")
	teacherGroup.Use(middleware.TeacherAuthMiddleware())
	{
		teacherGroup.GET("/", teacher.BaseGET)
		teacherGroup.POST("/login", auth.IsTeacher, Generic)
		teacherGroup.GET("/getAttendence", attendence.Test3)
		teacherGroup.POST("/getAttendence", attendence.CreateAttendence)
		teacherGroup.PUT("/getAttendence", attendence.EditAttendance)
	}

	api := r.Group("/api")
	{
		api.GET("/teachers", alloc_helper.GetTeachers)
		api.GET("/classroom", alloc_helper.GetClass)
	}

	r.Run(":3001")
}

func Generic(c *gin.Context) {
	c.JSON(200, gin.H{"message": c.Request.Method})
}
