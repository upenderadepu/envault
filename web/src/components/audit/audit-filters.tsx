"use client";

import { X } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

const ACTION_TYPES = [
  { value: "project.create", label: "Project Created" },
  { value: "project.delete", label: "Project Deleted" },
  { value: "secret.read", label: "Secret Read" },
  { value: "secret.write", label: "Secret Updated" },
  { value: "secret.delete", label: "Secret Deleted" },
  { value: "member.invite", label: "Member Invited" },
  { value: "member.remove", label: "Member Removed" },
  { value: "credentials.rotate", label: "Credentials Rotated" },
];

interface AuditFiltersProps {
  action: string;
  onActionChange: (action: string) => void;
}

export function AuditFilters({ action, onActionChange }: AuditFiltersProps) {
  return (
    <div className="flex items-center gap-3">
      <Select value={action || "all"} onValueChange={(v) => onActionChange(v === "all" ? "" : v)}>
        <SelectTrigger className="w-[200px]">
          <SelectValue placeholder="Filter by action" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem value="all">All Actions</SelectItem>
          {ACTION_TYPES.map((type) => (
            <SelectItem key={type.value} value={type.value}>
              {type.label}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
      {action && (
        <Button
          variant="ghost"
          size="sm"
          onClick={() => onActionChange("")}
          className="gap-1 text-muted-foreground"
        >
          <X className="h-3.5 w-3.5" />
          Clear
        </Button>
      )}
    </div>
  );
}
