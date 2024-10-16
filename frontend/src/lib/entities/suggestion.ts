import type { SvelteComponent } from "svelte";
import ArrowDown from "svelte-radix/ArrowDown.svelte";
import ArrowUp from "svelte-radix/ArrowUp.svelte";
import Minus from "svelte-radix/Minus.svelte";
import SewingPinFilled from "svelte-radix/SewingPinFilled.svelte";

export enum SuggestionImportance {
  CRITICAL = "CRITICAL",
  HIGH = "HIGH",
  MEDIUM = "MEDIUM",
  LOW = "LOW",
}

export type SuggestionImportanceDetail = {
  title: string;
  icon: typeof SvelteComponent;
};

export const SuggestionImportanceDetails: Record<SuggestionImportance | string, SuggestionImportanceDetail> = {
  [SuggestionImportance.CRITICAL]: {
    title: "Critical",
    icon: SewingPinFilled,
  },
  [SuggestionImportance.HIGH]: {
    title: "High",
    icon: ArrowUp,
  },
  [SuggestionImportance.MEDIUM]: {
    title: "Medium",
    icon: Minus,
  },
  [SuggestionImportance.LOW]: {
    title: "Low",
    icon: ArrowDown,
  },
};

export const SUGGESTION_NEW_MAX_DAYS: number = 7;

export type Suggestion = {
  id: string;
  productID: string;
  sources: Record<string, number>;
  title: string;
  description: string;
  reason: string;
  importances: Record<string, number>;
  priority: number;
  categories: Record<string, number>;
  releases: Record<string, number>;
  customers: number;
  assigneeID?: string;
  quality?: number;
  firstSeenAt: Date;
  lastSeenAt: Date;
  archivedAt?: Date;
};
