import * as entities from "$lib/entities";

export interface BaseCollectorSettings {}

export type TrustpilotCollectorSettings = BaseCollectorSettings & {
  domain: string;
};

export type PlayStoreCollectorSettings = BaseCollectorSettings & {
  app_id: string;
};

export type AppStoreCollectorSettings = BaseCollectorSettings & {
  app_id: string;
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

export function toCollectorSettings(settings: CollectorSettings, type: string): entities.CollectorSettings {
  switch (type) {
    case "TRUSTPILOT":
      return <entities.TrustpilotCollectorSettings>{
        domain: settings.domain,
      };
    case "PLAY_STORE":
      return <entities.PlayStoreCollectorSettings>{
        appID: settings.app_id,
      };
    case "APP_STORE":
      return <entities.AppStoreCollectorSettings>{
        appID: settings.app_id,
      };
    case "AMAZON":
      return <entities.AmazonCollectorSettings>{
        asin: settings.asin,
      };
    case "IAGORA":
      return <entities.IAgoraCollectorSettings>{
        institution: settings.institution,
      };
    default:
      throw new Error(`Unknown collector type ${type}`);
  }
}

export type Collector = {
  id: string;
  product_id: string;
  type: string;
  settings: CollectorSettings;
};

export function toCollector(collector: Collector): entities.Collector {
  return {
    id: collector.id,
    productID: collector.product_id,
    type: collector.type,
    settings: toCollectorSettings(collector.settings, collector.type),
  };
}

export type ListCollectorsResponse = {
  collectors: Collector[];
};

export type PostCollectorRequest = {
  type: string;
};

export type PostTrustpilotCollectorRequest = PostCollectorRequest & {
  domain: string;
};

export type PostPlayStoreCollectorRequest = PostCollectorRequest & {
  app_id: string;
};

export type PostAppStoreCollectorRequest = PostCollectorRequest & {
  app_id: string;
};

export type PostAmazonCollectorRequest = PostCollectorRequest & {
  asin: string;
};

export type PostIAgoraCollectorRequest = PostCollectorRequest & {
  institution: string;
};

export type PostCollectorResponse = Collector;

export type DeleteCollectorResponse = Record<PropertyKey, never>;
