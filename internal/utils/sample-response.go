package utils

type Success struct {
	Success bool   `json:"success" example:"true"`
	Result  string `json:"result" example:"message"`
	Message string `json:"message" example:"message"`
	Status  int    `json:"status" example:"200"`
	Error   string `json:"error" example:"null"`
}
type BadRequest struct {
	Success bool   `json:"success" example:"false"`
	Result  string `json:"result" example:"null"`
	Message string `json:"message" example:"message"`
	Status  int    `json:"status" example:"400"`
	Error   string `json:"error" example:"string"`
}
type InternalServerError struct {
	Success bool   `json:"success" example:"false"`
	Result  string `json:"result" example:"null"`
	Message string `json:"message" example:"message"`
	Status  int    `json:"status" example:"500"`
	Error   string `json:"error" example:"message"`
}
type SuccessDeploy struct {
	Success bool `json:"success" example:"true"`
	Result  struct {
		Cloudflare string `json:"cloudflare" example:"success add domain sub.example.com"`
		Portainer  string `json:"portainer" example:"success deploy portainer"`
		Jenkins    string `json:"jenkins" example:"success deploy jenkins with status build in proccess"`
	} `json:"result"`
	Message string `json:"message" example:"message"`
	Status  int    `json:"status" example:"200"`
	Error   string `json:"error" example:"null"`
}
type FoundFail struct {
	Success bool   `json:"success" example:"false"`
	Result  string `json:"result" example:"null"`
	Message string `json:"message" example:"message"`
	Status  int    `json:"status" example:"302"`
	Error   string `json:"error" example:"message"`
}
type FoundSuccess struct {
	Success bool   `json:"success" example:"true"`
	Result  string `json:"result" example:"message"`
	Message string `json:"message" example:"message"`
	Status  int    `json:"status" example:"302"`
	Error   string `json:"error" example:"null"`
}
