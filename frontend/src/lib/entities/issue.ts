import type { SvelteComponent } from "svelte";
import ArrowDown from "svelte-radix/ArrowDown.svelte";
import ArrowUp from "svelte-radix/ArrowUp.svelte";
import Minus from "svelte-radix/Minus.svelte";
import SewingPinFilled from "svelte-radix/SewingPinFilled.svelte";

export enum IssueSeverity {
  CRITICAL = "CRITICAL",
  HIGH = "HIGH",
  MEDIUM = "MEDIUM",
  LOW = "LOW",
}

export type IssueSeverityDetail = {
  title: string;
  icon: typeof SvelteComponent;
};

export const IssueSeverityDetails: Record<IssueSeverity | string, IssueSeverityDetail> = {
  [IssueSeverity.CRITICAL]: {
    title: "Critical",
    icon: SewingPinFilled,
  },
  [IssueSeverity.HIGH]: {
    title: "High",
    icon: ArrowUp,
  },
  [IssueSeverity.MEDIUM]: {
    title: "Medium",
    icon: Minus,
  },
  [IssueSeverity.LOW]: {
    title: "Low",
    icon: ArrowDown,
  },
};

export const ISSUE_NEW_MAX_DAYS: number = 7;

export type Issue = {
  id: string;
  productID: string;
  sources: Record<string, number>;
  title: string;
  description: string;
  steps: string[];
  severities: Record<string, number>;
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
