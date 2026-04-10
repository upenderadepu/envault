"use client";

import { useState } from "react";
import { AuditFilters } from "@/components/audit/audit-filters";
import { AuditLogTable } from "@/components/audit/audit-log-table";
import { useAuditLogs } from "@/hooks/use-audit";
import { useProjectContext } from "@/components/projects/project-context";

const PAGE_SIZE = 20;

export default function AuditPage() {
  const project = useProjectContext();
  const [action, setAction] = useState("");
  const [offset, setOffset] = useState(0);

  const { data, isLoading } = useAuditLogs(project.slug, {
    action: action || undefined,
    limit: PAGE_SIZE,
    offset,
  });

  const handleActionChange = (newAction: string) => {
    setAction(newAction);
    setOffset(0);
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold tracking-tight">Audit Log</h1>
        <p className="text-muted-foreground">
          Review all actions performed in {project.name}.
        </p>
      </div>

      <AuditFilters action={action} onActionChange={handleActionChange} />

      <AuditLogTable
        logs={data?.data}
        isLoading={isLoading}
        total={data?.total ?? 0}
        limit={PAGE_SIZE}
        offset={offset}
        onPageChange={setOffset}
      />
    </div>
  );
}
