package dtos

type ReusableCodeGetByCodeReq struct {
	Code string `json:"code" binding:"required"`
}
