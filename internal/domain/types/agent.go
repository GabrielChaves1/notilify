package types

import "github.com/google/uuid"

type AgentID uuid.UUID

func NewAgentIDFromString(agentIDStr string) (AgentID, error) {
	id, err := uuid.Parse(agentIDStr)
	if err != nil {
		return AgentID{}, err
	}
	return AgentID(id), nil
}
