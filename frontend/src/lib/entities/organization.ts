export enum OrganizationPlan {
  ENTERPRISE = "ENTERPRISE",
  BUSINESS = "BUSINESS",
  STARTER = "STARTER",
  TRIAL = "TRIAL",
  DEMO = "DEMO",
}

export type OrganizationPlanDetail = {
  title: string;
};

export const OrganizationPlanDetails: Record<OrganizationPlan | string, OrganizationPlanDetail> = {
  [OrganizationPlan.ENTERPRISE]: {
    title: "Enterprise",
  },
  [OrganizationPlan.BUSINESS]: {
    title: "Business",
  },
  [OrganizationPlan.STARTER]: {
    title: "Starter",
  },
  [OrganizationPlan.TRIAL]: {
    title: "Free trial",
  },
  [OrganizationPlan.DEMO]: {
    title: "Demo",
  },
};

export type OrganizationSettings = {
  domainSignIn: boolean;
  isDomainSignInSupported: boolean;
};

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
  trialEndsAt?: Date;
  capacity: OrganizationCapacity;
  usage: number;
};
