package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hanshal101/term-test-monitor/database/postgres"
	migrate "github.com/hanshal101/term-test-monitor/database/postgres/migration"
	alloc_helper "github.com/hanshal101/term-test-monitor/helpers/alloc_helpers"
	"github.com/hanshal101/term-test-monitor/helpers/auth"
	"github.com/hanshal101/term-test-monitor/internal/admin"
	students "github.com/hanshal101/term-test-monitor/internal/admin/students"
	"github.com/hanshal101/term-test-monitor/internal/admin/teachers"
	"github.com/hanshal101/term-test-monitor/internal/admin/vitals"
	"github.com/hanshal101/term-test-monitor/internal/dqc"
	"github.com/hanshal101/term-test-monitor/internal/teacher"
	"github.com/hanshal101/term-test-monitor/internal/teacher/attendence"
	dqcT "github.com/hanshal101/term-test-monitor/internal/teacher/dqc"
	"github.com/hanshal101/term-test-monitor/internal/teacher/papers"
	middleware "github.com/hanshal101/term-test-monitor/middleware/auth"
)

func init() {
	postgres.PostgresInitializer()
	migrate.Migrate()
}

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: true,
	}))

	adminGroup := r.Group("/admin") //
	{
		adminGroup.GET("/", admin.BaseGET) //
		adminGroup.GET("/:year/:subject/:class", students.DashboardAttendence)
		adminGroup.DELETE("/:year/:subject/:class", students.DeleteAttendence)
		adminGroup.PUT("/:year/:subject/:class", students.EditAttendence)
		create := adminGroup.Group("/create") //
		{
			create.POST("/timetable/:year/:sem", admin.CreateTimeTable) //
			student := create.Group("/student")                         //
			{
				student.POST("/dualAllocation", students.DualAllocation)     //
				student.POST("/singleAllocation", students.SingleAllocation) //
			}
			teacher := create.Group("/teacher")
			{
				teacher.POST("/allocation", teachers.CreateTeacherAllocation) //
				teacher.POST("/papers/:reqID/:req", teachers.MakePaperRequests)
			}
			vital := create.Group("/vitals") //
			{
				vital.GET("/", vitals.Base)                                  //
				vital.GET("/:year", vitals.GetSubject)                       //
				vital.POST("/:year", vitals.CreateSubject)                   //
				vital.DELETE("/:year/:subject", vitals.DeleteSubject)        //
				vital.GET("/teachers/:type", vitals.GetTeachers)             //
				vital.POST("/createTeacher", vitals.CreateTeacher)           //
				vital.PUT("/teachers/:type", vitals.EditTeacher)             //
				vital.DELETE("/teachers/:type/:email", vitals.DeleteTeacher) //
			}
		}
		get := adminGroup.Group("/get")
		{
			get.GET("/timetable", admin.GetTT)                                      //
			get.GET("/timetable/:Year", admin.GetTTbyYear)                          //
			get.DELETE("/timetable/:year", admin.DeleteTimeTable)                   //
			get.GET("/student/allocation", students.GetAllocation)                  //
			get.DELETE("/student/allocation/:id", students.DeleteAllocation)        //
			get.GET("/teacher/allocation", teachers.GetTeacherAllocation)           //
			get.DELETE("/teacher/allocation/:id", teachers.DeleteTeacherAllocation) //
			get.GET("/teacher/papers", teachers.GetPaperRequests)                   //bad requests issue
		}
	}

	teacherGroup := r.Group("/teacher")
	teacherGroup.POST("/login", auth.IsTeacherAuth)
	teacherGroup.Use(middleware.TeacherAuthMiddleware())
	{
		teacherGroup.GET("/", teacher.BaseGET)                           //
		teacherGroup.GET("/getAttendence", attendence.Test3)             //
		teacherGroup.POST("/getAttendence", attendence.CreateAttendence) //
		teacherGroup.PUT("/getAttendence", attendence.EditAttendance)    //
		teacherGroup.GET("/papers", papers.GetPaperRequest)              //
		teacherGroup.POST("/papers", papers.CreatePaperRequest)          //
		teacherGroup.DELETE("/papers/:reqID", papers.DeletePaperRequest) //
		teacherGroup.GET("/dqc/reviews", dqcT.GetReviewRequest)          //
		teacherGroup.POST("/dqc/reviews", dqcT.CreateDQCReview)          //
		teacherGroup.DELETE("/dqc/reviews/:reqID", dqcT.DeleteDQCReview) //
	}

	api := r.Group("/api") //not doing anything
	{
		api.GET("/teachers", alloc_helper.GetTeachers) //not doing anything
		// api.GET("/classroom", alloc_helper.GetClass)
	}

	dqcGroup := r.Group("/dqc")
	dqcGroup.POST("/login", auth.IsDqcAuth)
	dqcGroup.Use(middleware.DQCAuthMiddleware()) //something isssue
	{
		dqcGroup.GET("/", func(c *gin.Context) { c.String(200, "You are at DQC routes") })
		dqcGroup.GET("/requests", dqc.GetReviews)
		dqcGroup.GET("/requests/:reqID", dqc.GetReviewbyID)
		dqcGroup.POST("/requests/:reqID/:req", dqc.MakeReviewRequest)
		// dqc.GET("/timeline")
	}

	r.Run(":5051")
}

func Generic(c *gin.Context) {
	c.JSON(200, gin.H{"message": c.Request.Method})
}
