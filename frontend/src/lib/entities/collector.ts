import Amazon from "$lib/components/icons/Amazon.svelte";
import AppStore from "$lib/components/icons/AppStore.svelte";
import IAgora from "$lib/components/icons/IAgora.svelte";
import PlayStore from "$lib/components/icons/PlayStore.svelte";
import Trustpilot from "$lib/components/icons/Trustpilot.svelte";
import type { SvelteComponent } from "svelte";

export enum CollectorType {
  TRUSTPILOT = "TRUSTPILOT",
  PLAY_STORE = "PLAY_STORE",
  APP_STORE = "APP_STORE",
  AMAZON = "AMAZON",
  IAGORA = "IAGORA",
}

export interface BaseCollectorSettings {}

export type TrustpilotCollectorSettings = BaseCollectorSettings & {
  domain: string;
};

export type PlayStoreCollectorSettings = BaseCollectorSettings & {
  appID: string;
};

export type AppStoreCollectorSettings = BaseCollectorSettings & {
  appID: string;
};

export type AmazonCollectorSettings = BaseCollectorSettings & {
  asin: string;
};

export type IAgoraCollectorSettings = BaseCollectorSettings & {
  institution: string;
};

export type CollectorSettings =
  | TrustpilotCollectorSettings
  | PlayStoreCollectorSettings
  | AppStoreCollectorSettings
  | AmazonCollectorSettings
  | IAgoraCollectorSettings;

export type Collector = {
  id: string;
  productID: string;
  type: string;
  settings: CollectorSettings;
};

export type CollectorDetail = {
  title: string;
  description: string;
  icon: typeof SvelteComponent;
  base?: string;
};

export const CollectorDetails: Record<CollectorType | string, CollectorDetail> = {
  [CollectorType.TRUSTPILOT]: {
    title: "Trustpilot",
    description: "Collects reviews from your company page on Trustpilot (trustpilot.com).",
    icon: Trustpilot,
    base: "https://www.trustpilot.com/review/",
  },
  [CollectorType.PLAY_STORE]: {
    title: "Play Store",
    description: "Collects reviews from your application page on Play Store (all storefronts).",
    icon: PlayStore,
    base: "https://play.google.com/store/apps/details?id=",
  },
  [CollectorType.APP_STORE]: {
    title: "App Store",
    description: "Collects reviews from your application page on App Store (all storefronts).",
    icon: AppStore,
    base: "https://apps.apple.com/app/id",
  },
  [CollectorType.AMAZON]: {
    title: "Amazon",
    description: "Collects reviews from your product page on Amazon (all storefronts).",
    icon: Amazon,
    base: "https://www.amazon.com/dp/",
  },
  [CollectorType.IAGORA]: {
    title: "iAgora",
    description: "Collects reviews from your institution page on iAgora (iagora.com).",
    icon: IAgora,
    base: "https://www.iagora.com/studies/uni/",
  },
};
