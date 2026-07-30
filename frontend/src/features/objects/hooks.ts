import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { objectsApi } from "./api";
import { buildFolderViewModels } from "./mappers";
import { ObjectViewModel } from "./types";

export function useObjects(bucket: string, prefix = "") {
  return useQuery<ObjectViewModel[], Error>({
    queryKey: ["objects", bucket, prefix],
    queryFn: async () => {
      const dtos = await objectsApi.listObjects(bucket, prefix);
      return buildFolderViewModels(dtos, prefix);
    },
    staleTime: 15000,
    enabled: !!bucket,
  });
}

export function useDeleteObject(bucket: string) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (key: string) => {
      return objectsApi.deleteObject(bucket, key);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["objects", bucket] });
      queryClient.invalidateQueries({ queryKey: ["buckets"] });
    },
  });
}
