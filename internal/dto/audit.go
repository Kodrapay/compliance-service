package dto

type AuditLogRequest struct {
	ActorID   int    `json:"actor_id"`
	ActorRole string `json:"actor_role"`
	Action    string `json:"action"`
	Entity    string `json:"entity"`
	EntityID  int    `json:"entity_id"`
	IP        string `json:"ip"`
}

type AuditLogResponse struct {
	ID int `json:"id"`
}
