export interface ObjectDTO {
  key: string;
  bucket: string;
  size: number;
  last_modified?: string;
  etag?: string;
  content_type?: string;
}

export interface ObjectViewModel {
  key: string;
  name: string;
  isFolder: boolean;
  prefix: string;
  sizeFormatted: string;
  sizeRaw: number;
  lastModifiedFormatted: string;
  contentType?: string;
}
