export enum UserRole {
  ADMIN = "ADMIN",
  MEMBER = "MEMBER",
}

export type UserSettings = Record<PropertyKey, never>;

export type User = {
  id: string;
  organizationID: string;
  name: string;
  picture: string;
  email: string;
  role: string;
  settings: UserSettings;
  leftAt?: Date;
};

export type Invitation = {
  id: string;
  organizationID: string;
  email: string;
  role: string;
  expiresAt: Date;
};
