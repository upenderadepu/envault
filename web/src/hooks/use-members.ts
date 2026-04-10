"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api";

export function useMembers(slug: string) {
  return useQuery({
    queryKey: ["members", slug],
    queryFn: () => api.listMembers(slug),
    enabled: !!slug,
  });
}

export function useAddMember(slug: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: { email: string; role: string }) =>
      api.addMember(slug, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["members", slug] });
    },
  });
}

export function useRemoveMember(slug: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (memberId: string) => api.removeMember(slug, memberId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["members", slug] });
    },
  });
}

export function useRotateCredentials(slug: string) {
  return useMutation({
    mutationFn: () => api.rotateCredentials(slug),
  });
}
