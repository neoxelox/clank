import * as entities from "$lib/entities";

export type OrganizationSettings = {
  domain_sign_in: boolean;
  is_domain_sign_in_supported: boolean;
};

export function toOrganizationSettings(settings: OrganizationSettings): entities.OrganizationSettings {
  return {
    domainSignIn: settings.domain_sign_in,
    isDomainSignInSupported: settings.is_domain_sign_in_supported,
  };
}

export type OrganizationCapacity = {
  included: number;
  extra: number;
};

export type Organization = {
  id: string;
  name: string;
  picture: string;
  domain: string;
  settings: OrganizationSettings;
  plan: string;
  trial_ends_at?: Date;
  capacity: OrganizationCapacity;
  usage: number;
};

export function toOrganization(organization: Organization): entities.Organization {
  return {
    id: organization.id,
    name: organization.name,
    picture: organization.picture,
    domain: organization.domain,
    settings: toOrganizationSettings(organization.settings),
    plan: organization.plan,
    trialEndsAt: organization.trial_ends_at,
    capacity: {
      included: organization.capacity.included,
      extra: organization.capacity.extra,
    },
    usage: organization.usage,
  };
}

export type GetOrganizationResponse = Organization;

export type PutOrganizationRequest = {
  name?: string;
  picture?: string;
};

export type PutOrganizationResponse = Organization;

export type DeleteOrganizationResponse = Record<PropertyKey, never>;

export type PutOrganizationSettingsRequest = {
  domain_sign_in?: boolean;
};

export type PutOrganizationSettingsResponse = OrganizationSettings;
