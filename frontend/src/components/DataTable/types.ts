import type { Component } from "vue";

export type ColumnPosition = "left" | "center" | "right";

export interface Column {
  key: string;
  label: string;
  sortable?: boolean;
  thPosition?: ColumnPosition;
  tdPosition?: ColumnPosition;
  component?: Component;
  componentProps?: Record<string, any>;
}

export interface SortState {
  key: string;
  direction: "asc" | "desc";
}

export interface PaginationMeta {
  currentPage: number;
  pageSize: number;
  totalItems: number;
  totalPages: number;
}
