import * as entities from "$lib/entities";

export type UserSettings = Record<PropertyKey, never>;

export type User = {
  id: string;
  organization_id: string;
  name: string;
  picture: string;
  email: string;
  role: string;
  settings: UserSettings;
  left_at?: Date;
};

export function toUser(user: User): entities.User {
  return {
    id: user.id,
    organizationID: user.organization_id,
    name: user.name,
    picture: user.picture,
    email: user.email,
    role: user.role,
    settings: {},
    leftAt: user.left_at,
  };
}

export type GetMeResponse = User;

export type PutMeRequest = {
  name?: string;
  picture?: string;
};

export type PutMeResponse = User;

export type DeleteMeResponse = Record<PropertyKey, never>;

export type ListUsersResponse = {
  users: User[];
};

export type PutUserRequest = {
  role?: string;
};

export type PutUserResponse = User;

export type DeleteUserResponse = Record<PropertyKey, never>;

export type Invitation = {
  id: string;
  organization_id: string;
  email: string;
  role: string;
  expires_at: Date;
};

export function toInvitation(invitation: Invitation): entities.Invitation {
  return {
    id: invitation.id,
    organizationID: invitation.organization_id,
    email: invitation.email,
    role: invitation.role,
    expiresAt: invitation.expires_at,
  };
}

export type ListInvitationsResponse = {
  invitations: Invitation[];
};

export type PostInvitationRequest = {
  email: string;
  role: string;
};

export type PostInvitationResponse = Invitation;

export type DeleteInvitationResponse = Record<PropertyKey, never>;
