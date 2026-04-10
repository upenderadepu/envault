"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api";

export function useSecrets(slug: string, environment: string) {
  return useQuery({
    queryKey: ["secrets", slug, environment],
    queryFn: () => api.listSecrets(slug, environment),
    enabled: !!slug && !!environment,
  });
}

export function useSetSecret(slug: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: { environment: string; key: string; value: string }) =>
      api.setSecret(slug, data),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({
        queryKey: ["secrets", slug, variables.environment],
      });
    },
  });
}

export function useDeleteSecret(slug: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (params: { key: string; environment: string }) =>
      api.deleteSecret(slug, params.key, params.environment),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({
        queryKey: ["secrets", slug, variables.environment],
      });
    },
  });
}

export function useBulkSetSecrets(slug: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: {
      environment: string;
      secrets: Record<string, string>;
    }) => api.bulkSetSecrets(slug, data),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({
        queryKey: ["secrets", slug, variables.environment],
      });
    },
  });
}
