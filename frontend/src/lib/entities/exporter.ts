import type { SvelteComponent } from "svelte";

export enum ExporterType {
  SLACK = "SLACK",
  JIRA = "JIRA",
}

export interface BaseExporterSettings {}

export type SlackExporterSettings = BaseExporterSettings & {
  channel: string;
};

export type JiraExporterSettings = BaseExporterSettings & {
  board: string;
};

export type ExporterSettings = SlackExporterSettings | JiraExporterSettings;

export type Exporter = {
  id: string;
  productID: string;
  type: string;
  settings: ExporterSettings;
};

export type ExporterDetail = {
  title: string;
  description: string;
  icon: typeof SvelteComponent;
};

export const ExporterDetails: Record<ExporterType | string, ExporterDetail> = {};
