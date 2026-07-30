import { describe, it, expect } from "vitest";
import { mapBucketDTOToViewModel } from "../mappers";
import { createBucketSchema } from "../schemas";

describe("mapBucketDTOToViewModel", () => {
  it("should format ISO timestamp cleanly", () => {
    const dto = {
      name: "prod-media",
      created_at: "2026-07-30T10:00:00Z",
    };
    const vm = mapBucketDTOToViewModel(dto);
    expect(vm.name).toBe("prod-media");
    expect(vm.createdAtFormatted).toContain("Jul 30, 2026");
  });
});

describe("createBucketSchema", () => {
  it("should validate valid bucket names", () => {
    expect(createBucketSchema.safeParse({ name: "my-bucket-123" }).success).toBe(true);
    expect(createBucketSchema.safeParse({ name: "abc" }).success).toBe(true);
  });

  it("should reject invalid bucket names", () => {
    expect(createBucketSchema.safeParse({ name: "ab" }).success).toBe(false);
    expect(createBucketSchema.safeParse({ name: "MY-BUCKET" }).success).toBe(false);
    expect(createBucketSchema.safeParse({ name: "-start-with-hyphen" }).success).toBe(false);
  });
});
