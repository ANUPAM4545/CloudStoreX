import { format, parseISO, isValid } from "date-fns";
import { BucketDTO, BucketViewModel } from "./types";

export function mapBucketDTOToViewModel(dto: BucketDTO): BucketViewModel {
  let formatted = "Unknown date";
  if (dto.created_at) {
    const parsed = parseISO(dto.created_at);
    if (isValid(parsed)) {
      formatted = format(parsed, "MMM d, yyyy • HH:mm");
    } else {
      formatted = dto.created_at;
    }
  }

  return {
    name: dto.name,
    createdAtFormatted: formatted,
    createdAtRaw: dto.created_at || "",
  };
}
