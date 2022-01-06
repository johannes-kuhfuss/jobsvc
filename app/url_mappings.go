package app

func mapUrls() {
	cfg.RunTime.Router.POST("/jobs", jobHandler.CreateJob)
	cfg.RunTime.Router.GET("/jobs", jobHandler.GetAllJobs)
	cfg.RunTime.Router.GET("/jobs/:job_id", jobHandler.GetJobById)
	cfg.RunTime.Router.DELETE("/jobs/:job_id", jobHandler.DeleteJobById)
	cfg.RunTime.Router.DELETE("/jobs", jobHandler.DeleteAllJobs)
	cfg.RunTime.Router.PUT("/jobs/:job_id", jobHandler.UpdateJob)
	cfg.RunTime.Router.PUT("/jobs/:job_id/status", jobHandler.SetStatusById)
	cfg.RunTime.Router.PUT("/jobs/:job_id/history", jobHandler.SetHistoryById)
	cfg.RunTime.Router.PUT("/jobs/dequeue", jobHandler.Dequeue)
}
