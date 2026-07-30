import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { bucketsApi } from "./api";
import { mapBucketDTOToViewModel } from "./mappers";
import { BucketViewModel } from "./types";

export function useBuckets() {
  return useQuery<BucketViewModel[], Error>({
    queryKey: ["buckets"],
    queryFn: async () => {
      const dtos = await bucketsApi.listBuckets();
      return dtos.map(mapBucketDTOToViewModel);
    },
    staleTime: 30000,
  });
}

export function useCreateBucket() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (name: string) => {
      return bucketsApi.createBucket(name);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["buckets"] });
    },
  });
}

export function useDeleteBucket() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (name: string) => {
      return bucketsApi.deleteBucket(name);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["buckets"] });
    },
  });
}
