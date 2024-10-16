package user

import "time"

type UserPayloadSettings struct {
}

type UserPayload struct {
	ID             string              `json:"id"`
	OrganizationID string              `json:"organization_id"`
	Name           string              `json:"name"`
	Picture        string              `json:"picture"`
	Email          string              `json:"email"`
	Role           string              `json:"role"`
	Settings       UserPayloadSettings `json:"settings"`
	LeftAt         *time.Time          `json:"left_at"`
}

func NewUserPayload(user User) *UserPayload {
	return &UserPayload{
		ID:             user.ID,
		OrganizationID: user.OrganizationID,
		Name:           user.Name,
		Picture:        user.Picture,
		Email:          user.Email,
		Role:           user.Role,
		Settings:       UserPayloadSettings{},
		LeftAt:         user.DeletedAt,
	}
}

type InvitationPayload struct {
	ID             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	Email          string    `json:"email"`
	Role           string    `json:"role"`
	ExpiresAt      time.Time `json:"expires_at"`
}

func NewInvitationPayload(invitation Invitation) *InvitationPayload {
	return &InvitationPayload{
		ID:             invitation.ID,
		OrganizationID: invitation.OrganizationID,
		Email:          invitation.Email,
		Role:           invitation.Role,
		ExpiresAt:      invitation.ExpiresAt,
	}
}
