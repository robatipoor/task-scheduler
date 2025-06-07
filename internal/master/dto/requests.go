package dto

type SubmitTaskRequest struct {
	TrackUID  string `json:"track_uid" binding:"required"`
	Input    uint   `json:"input" binding:"required,min=1,max=100000"`
	Priority uint   `json:"priority" binding:"required,min=1,max=100"`
}

type WorkerUpdateResultTaskRequest struct {
	TrackID uint `json:"track_id" binding:"required"`
	Result  uint `json:"result" binding:"required"`
}

type WorkerRegisterRequest struct {
	BaseUrl string `json:"base_url" binding:"required"`
}
