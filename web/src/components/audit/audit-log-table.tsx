"use client";

import { useState } from "react";
import { ChevronDown, ChevronRight } from "lucide-react";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { formatDateTime, actionLabel } from "@/lib/utils";
import type { AuditLog } from "@/lib/types";

interface AuditLogTableProps {
  logs: AuditLog[] | undefined;
  isLoading: boolean;
  total: number;
  limit: number;
  offset: number;
  onPageChange: (offset: number) => void;
}

function actionVariant(
  action: string
): "default" | "secondary" | "destructive" | "outline" {
  if (action.includes("delete") || action.includes("remove"))
    return "destructive";
  if (action.includes("create") || action.includes("invite")) return "default";
  if (action.includes("write") || action.includes("rotate"))
    return "secondary";
  return "outline";
}

function ExpandableRow({ log }: { log: AuditLog }) {
  const [expanded, setExpanded] = useState(false);

  return (
    <>
      <TableRow
        className="cursor-pointer hover:bg-muted/50"
        onClick={() => setExpanded(!expanded)}
      >
        <TableCell className="w-8">
          {expanded ? (
            <ChevronDown className="h-4 w-4 text-muted-foreground" />
          ) : (
            <ChevronRight className="h-4 w-4 text-muted-foreground" />
          )}
        </TableCell>
        <TableCell className="text-sm text-muted-foreground whitespace-nowrap">
          {formatDateTime(log.created_at)}
        </TableCell>
        <TableCell>
          <Badge variant={actionVariant(log.action)} className="text-xs">
            {actionLabel(log.action)}
          </Badge>
        </TableCell>
        <TableCell className="font-mono text-sm text-muted-foreground truncate max-w-[200px]">
          {log.resource_path}
        </TableCell>
        <TableCell className="text-sm text-muted-foreground">
          {log.user_id ? log.user_id.slice(0, 8) + "..." : "System"}
        </TableCell>
      </TableRow>
      {expanded && (
        <TableRow>
          <TableCell colSpan={5} className="bg-muted/30 p-4">
            <div className="space-y-2 text-sm">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <span className="font-medium text-muted-foreground">
                    Request ID
                  </span>
                  <p className="font-mono text-xs mt-1">
                    {log.request_id || "—"}
                  </p>
                </div>
                <div>
                  <span className="font-medium text-muted-foreground">
                    IP Address
                  </span>
                  <p className="font-mono text-xs mt-1">
                    {log.ip_address || "—"}
                  </p>
                </div>
              </div>
              {log.metadata &&
                Object.keys(log.metadata).length > 0 && (
                  <div>
                    <span className="font-medium text-muted-foreground">
                      Metadata
                    </span>
                    <pre className="mt-1 rounded-md bg-muted p-3 text-xs font-mono overflow-x-auto">
                      {JSON.stringify(log.metadata, null, 2)}
                    </pre>
                  </div>
                )}
            </div>
          </TableCell>
        </TableRow>
      )}
    </>
  );
}

export function AuditLogTable({
  logs,
  isLoading,
  total,
  limit,
  offset,
  onPageChange,
}: AuditLogTableProps) {
  if (isLoading) {
    return (
      <div className="space-y-3">
        {[...Array(5)].map((_, i) => (
          <Skeleton key={i} className="h-12 w-full" />
        ))}
      </div>
    );
  }

  if (!logs || logs.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-12 text-center">
        <p className="text-sm text-muted-foreground">No audit logs found.</p>
      </div>
    );
  }

  const totalPages = Math.ceil(total / limit);
  const currentPage = Math.floor(offset / limit) + 1;

  return (
    <div className="space-y-4">
      <div className="rounded-lg border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="w-8"></TableHead>
              <TableHead className="w-[180px]">Timestamp</TableHead>
              <TableHead className="w-[160px]">Action</TableHead>
              <TableHead>Resource</TableHead>
              <TableHead className="w-[120px]">User</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {logs.map((log) => (
              <ExpandableRow key={log.id} log={log} />
            ))}
          </TableBody>
        </Table>
      </div>

      {totalPages > 1 && (
        <div className="flex items-center justify-between">
          <p className="text-sm text-muted-foreground">
            Showing {offset + 1}–{Math.min(offset + limit, total)} of {total}
          </p>
          <div className="flex items-center gap-2">
            <Button
              variant="outline"
              size="sm"
              disabled={offset === 0}
              onClick={() => onPageChange(Math.max(0, offset - limit))}
            >
              Previous
            </Button>
            <span className="text-sm text-muted-foreground">
              Page {currentPage} of {totalPages}
            </span>
            <Button
              variant="outline"
              size="sm"
              disabled={offset + limit >= total}
              onClick={() => onPageChange(offset + limit)}
            >
              Next
            </Button>
          </div>
        </div>
      )}
    </div>
  );
}
