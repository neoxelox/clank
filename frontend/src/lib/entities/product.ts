export enum ProductLanguage {
  ENGLISH = "ENGLISH",
}

export const NO_CATEGORY: string = "UNKNOWN";

export type ProductSettings = Record<PropertyKey, never>;

export type Product = {
  id: string;
  organizationID: string;
  name: string;
  picture: string;
  language: string;
  context: string;
  categories: string[];
  release: string;
  settings: ProductSettings;
  usage: number;
};
