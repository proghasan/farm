import type { Component } from "vue";

export interface Column {
  key: string;
  label: string;
  sortable?: boolean;
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
