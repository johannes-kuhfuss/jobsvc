package app

func mapUrls() {
	appCfg.RunTime.Router.POST("/jobs", jobHandler.CreateJob)
	appCfg.RunTime.Router.GET("/jobs", jobHandler.GetAllJobs)
	/*
		router.GET("jobs/:job_id", )
		router.POST("/jobs", )
		router.DELETE("jobs/:job_id", )
		router.GET("/jobs/next", )
	*/
}
