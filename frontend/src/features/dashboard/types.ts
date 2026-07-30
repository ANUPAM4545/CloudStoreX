export interface DashboardStatsViewModel {
  totalBuckets: number;
  activeProvider: string;
  providerStatus: "success" | "warning" | "error";
  maxUploadSizeMB: number;
  storagePolicy: string;
}
