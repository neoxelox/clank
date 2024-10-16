import * as entities from "$lib/entities";

export type FeedbackCustomer = {
  email?: string;
  name: string;
  picture: string;
  location?: string;
  verified?: boolean;
  reviews?: number;
  link?: string;
};

export type FeedbackMetadata = {
  rating?: number;
  media?: string[];
  verified?: boolean;
  votes?: number;
  link?: string;
};

export type Feedback = {
  id: string;
  product_id: string;
  source: string;
  customer: FeedbackCustomer;
  content: string;
  language: string;
  translation: string;
  release: string;
  metadata: FeedbackMetadata;
  posted_at: Date;
};

export function toFeedback(feedback: Feedback): entities.Feedback {
  return {
    id: feedback.id,
    productID: feedback.product_id,
    source: feedback.source,
    customer: {
      email: feedback.customer.email,
      name: feedback.customer.name,
      picture: feedback.customer.picture,
      location: feedback.customer.location,
      verified: feedback.customer.verified,
      reviews: feedback.customer.reviews,
      link: feedback.customer.link,
    },
    content: feedback.content,
    language: feedback.language,
    translation: feedback.translation,
    release: feedback.release,
    metadata: {
      rating: feedback.metadata.rating,
      media: feedback.metadata.media,
      verified: feedback.metadata.verified,
      votes: feedback.metadata.votes,
      link: feedback.metadata.link,
    },
    postedAt: feedback.posted_at,
  };
}
