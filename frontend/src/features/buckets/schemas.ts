import { z } from "zod";

export const createBucketSchema = z.object({
  name: z
    .string()
    .min(3, "Bucket name must be at least 3 characters long.")
    .max(63, "Bucket name must not exceed 63 characters.")
    .regex(
      /^[a-z0-9]([a-z0-9-]*[a-z0-9])?$/,
      "Bucket name must be lowercase letters, numbers, and hyphens, and cannot start or end with a hyphen."
    ),
});

export type CreateBucketFormData = z.infer<typeof createBucketSchema>;
