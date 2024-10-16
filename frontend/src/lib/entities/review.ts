import Acceptance from "$lib/components/icons/Acceptance.svelte";
import Anger from "$lib/components/icons/Anger.svelte";
import Annoyance from "$lib/components/icons/Annoyance.svelte";
import Anticipation from "$lib/components/icons/Anticipation.svelte";
import Apprehension from "$lib/components/icons/Apprehension.svelte";
import Boredom from "$lib/components/icons/Boredom.svelte";
import Disgust from "$lib/components/icons/Disgust.svelte";
import Distraction from "$lib/components/icons/Distraction.svelte";
import Fear from "$lib/components/icons/Fear.svelte";
import Interest from "$lib/components/icons/Interest.svelte";
import Joy from "$lib/components/icons/Joy.svelte";
import Pensiveness from "$lib/components/icons/Pensiveness.svelte";
import Sadness from "$lib/components/icons/Sadness.svelte";
import Serenity from "$lib/components/icons/Serenity.svelte";
import Surprise from "$lib/components/icons/Surprise.svelte";
import Trust from "$lib/components/icons/Trust.svelte";
import Angry from "lucide-svelte/icons/angry";
import Frown from "lucide-svelte/icons/frown";
import Smile from "lucide-svelte/icons/smile";
import SmilePlus from "lucide-svelte/icons/smile-plus";
import type { SvelteComponent } from "svelte";
import ArrowDown from "svelte-radix/ArrowDown.svelte";
import ArrowUp from "svelte-radix/ArrowUp.svelte";
import Minus from "svelte-radix/Minus.svelte";
import type { Feedback } from "./feedback";

export enum ReviewSentiment {
  POSITIVE = "POSITIVE",
  NEUTRAL = "NEUTRAL",
  NEGATIVE = "NEGATIVE",
}

export type ReviewSentimentDetail = {
  title: string;
  icon: typeof SvelteComponent;
};

export const ReviewSentimentDetails: Record<ReviewSentiment | string, ReviewSentimentDetail> = {
  [ReviewSentiment.POSITIVE]: {
    title: "Positive",
    icon: ArrowUp,
  },
  [ReviewSentiment.NEUTRAL]: {
    title: "Neutral",
    icon: Minus,
  },
  [ReviewSentiment.NEGATIVE]: {
    title: "Negative",
    icon: ArrowDown,
  },
};

export enum ReviewEmotion {
  SERENITY = "SERENITY",
  JOY = "JOY",
  ACCEPTANCE = "ACCEPTANCE",
  TRUST = "TRUST",
  APPREHENSION = "APPREHENSION",
  FEAR = "FEAR",
  DISTRACTION = "DISTRACTION",
  SURPRISE = "SURPRISE",
  PENSIVENESS = "PENSIVENESS",
  SADNESS = "SADNESS",
  BOREDOM = "BOREDOM",
  DISGUST = "DISGUST",
  ANNOYANCE = "ANNOYANCE",
  ANGER = "ANGER",
  INTEREST = "INTEREST",
  ANTICIPATION = "ANTICIPATION",
}

export type ReviewEmotionDetail = {
  title: string;
  icon: typeof SvelteComponent;
};

export const ReviewEmotionDetails: Record<ReviewEmotion | string, ReviewEmotionDetail> = {
  [ReviewEmotion.TRUST]: {
    title: "Trust",
    icon: Trust,
  },
  [ReviewEmotion.ACCEPTANCE]: {
    title: "Acceptance",
    icon: Acceptance,
  },
  [ReviewEmotion.FEAR]: {
    title: "Fear",
    icon: Fear,
  },
  [ReviewEmotion.APPREHENSION]: {
    title: "Apprehension",
    icon: Apprehension,
  },
  [ReviewEmotion.SURPRISE]: {
    title: "Surprise",
    icon: Surprise,
  },
  [ReviewEmotion.DISTRACTION]: {
    title: "Distraction",
    icon: Distraction,
  },
  [ReviewEmotion.SADNESS]: {
    title: "Sadness",
    icon: Sadness,
  },
  [ReviewEmotion.PENSIVENESS]: {
    title: "Pensiveness",
    icon: Pensiveness,
  },
  [ReviewEmotion.DISGUST]: {
    title: "Disgust",
    icon: Disgust,
  },
  [ReviewEmotion.BOREDOM]: {
    title: "Boredom",
    icon: Boredom,
  },
  [ReviewEmotion.ANGER]: {
    title: "Anger",
    icon: Anger,
  },
  [ReviewEmotion.ANNOYANCE]: {
    title: "Annoyance",
    icon: Annoyance,
  },
  [ReviewEmotion.ANTICIPATION]: {
    title: "Anticipation",
    icon: Anticipation,
  },
  [ReviewEmotion.INTEREST]: {
    title: "Interest",
    icon: Interest,
  },
  [ReviewEmotion.JOY]: {
    title: "Joy",
    icon: Joy,
  },
  [ReviewEmotion.SERENITY]: {
    title: "Serenity",
    icon: Serenity,
  },
};

export enum ReviewIntention {
  RETAIN = "RETAIN",
  CHURN = "CHURN",
  RETAIN_AND_RECOMMEND = "RETAIN_AND_RECOMMEND",
  CHURN_AND_DISCOURAGE = "CHURN_AND_DISCOURAGE",
}

export type ReviewIntentionDetail = {
  title: string;
  icon: typeof SvelteComponent;
};

export const ReviewIntentionDetails: Record<ReviewIntention | string, ReviewIntentionDetail> = {
  [ReviewIntention.RETAIN]: {
    title: "Retain",
    icon: Smile,
  },
  [ReviewIntention.CHURN]: {
    title: "Churn",
    icon: Frown,
  },
  [ReviewIntention.RETAIN_AND_RECOMMEND]: {
    title: "Retain and Recommend",
    icon: SmilePlus,
  },
  [ReviewIntention.CHURN_AND_DISCOURAGE]: {
    title: "Churn and Discourage",
    icon: Angry,
  },
};

export const NO_INTENTION: string = "UNKNOWN";

export type Review = {
  id: string;
  productID: string;
  feedback: Feedback;
  keywords: string[];
  sentiment: string;
  emotions: string[];
  intention: string;
  category: string;
  quality?: number;
};
