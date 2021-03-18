package models

type UploadRequest struct {
	File    string `json:"file"`
	TraceId string `json:"traceId"`
}
