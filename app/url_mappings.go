package app

func mapUrls() {
	cfg.RunTime.Router.POST("/jobs", jobHandler.CreateJob)
	cfg.RunTime.Router.GET("/jobs", jobHandler.GetAllJobs)
	cfg.RunTime.Router.GET("jobs/:job_id", jobHandler.GetJobById)
	cfg.RunTime.Router.DELETE("jobs/:job_id", jobHandler.DeleteJobById)
	cfg.RunTime.Router.GET("/jobs/next", jobHandler.GetNextJob)
}
