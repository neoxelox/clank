import * as entities from "$lib/entities";

export type ProductSettings = Record<PropertyKey, never>;

export type Product = {
  id: string;
  organization_id: string;
  name: string;
  picture: string;
  language: string;
  context: string;
  categories: string[];
  release: string;
  settings: ProductSettings;
  usage: number;
};

export function toProduct(product: Product): entities.Product {
  return {
    id: product.id,
    organizationID: product.organization_id,
    name: product.name,
    picture: product.picture,
    language: product.language,
    context: product.context,
    categories: product.categories,
    release: product.release,
    settings: {},
    usage: product.usage,
  };
}

export type ListProductsResponse = {
  products: Product[];
};

export type GetProductResponse = Product;

export type PostProductRequest = {
  name: string;
  picture?: string;
  language: string;
  context?: string;
  categories?: string[];
  release?: string;
};

export type PostProductResponse = Product;

export type PutProductRequest = {
  name?: string;
  picture?: string;
  context?: string;
  categories?: string[];
  release?: string;
};

export type PutProductResponse = Product;

export type DeleteProductResponse = Record<PropertyKey, never>;
