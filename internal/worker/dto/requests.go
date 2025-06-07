package dto

type SubmitTaskRequest struct {
	TrackID uint `json:"track_id" binding:"required"`
	Input   uint `json:"input" binding:"required"`
}

type UpdateResultTaskRequest struct {
	ID     uint `json:"id" binding:"required"`
	Result uint `json:"result" binding:"required"`
}

type WorkerRegisterRequest struct {
	BaseUrl string `json:"baseUrl" binding:"required"`
}
