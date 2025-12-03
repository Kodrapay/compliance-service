package dto

type AuditLogRequest struct {
	ActorID   string `json:"actor_id"`
	ActorRole string `json:"actor_role"`
	Action    string `json:"action"`
	Entity    string `json:"entity"`
	EntityID  string `json:"entity_id"`
	IP        string `json:"ip"`
}

type AuditLogResponse struct {
	ID string `json:"id"`
}
