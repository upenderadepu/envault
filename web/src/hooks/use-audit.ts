"use client";

import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";

export function useAuditLogs(
  slug: string,
  params: { action?: string; limit?: number; offset?: number } = {}
) {
  return useQuery({
    queryKey: ["audit", slug, params],
    queryFn: () => api.listAuditLogs(slug, params),
    enabled: !!slug,
  });
}
