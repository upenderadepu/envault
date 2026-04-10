"use client";

import { Trash2 } from "lucide-react";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { Badge } from "@/components/ui/badge";
import { SecretValueCell } from "./secret-value-cell";
import { useDeleteSecret } from "@/hooks/use-secrets";
import { formatRelativeTime } from "@/lib/utils";
import { toast } from "sonner";
import type { SecretMetadata } from "@/lib/types";

interface SecretsTableProps {
  secrets: SecretMetadata[] | undefined;
  isLoading: boolean;
  slug: string;
  environment: string;
}

export function SecretsTable({
  secrets,
  isLoading,
  slug,
  environment,
}: SecretsTableProps) {
  const deleteSecret = useDeleteSecret(slug);

  const handleDelete = async (key: string) => {
    try {
      await deleteSecret.mutateAsync({ key, environment });
      toast.success(`Secret "${key}" deleted`);
    } catch (err) {
      toast.error(
        err instanceof Error ? err.message : "Failed to delete secret"
      );
    }
  };

  if (isLoading) {
    return (
      <div className="space-y-3">
        {[...Array(3)].map((_, i) => (
          <Skeleton key={i} className="h-12 w-full" />
        ))}
      </div>
    );
  }

  if (!secrets || secrets.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center py-12 text-center">
        <p className="text-sm text-muted-foreground">
          No secrets in this environment yet.
        </p>
      </div>
    );
  }

  return (
    <div className="rounded-lg border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead className="w-[250px]">Key</TableHead>
            <TableHead>Value</TableHead>
            <TableHead className="w-[80px]">Version</TableHead>
            <TableHead className="w-[140px]">Modified</TableHead>
            <TableHead className="w-[60px]"></TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {secrets.map((secret) => (
            <TableRow key={secret.id}>
              <TableCell className="font-mono font-medium text-sm">
                {secret.key_name}
              </TableCell>
              <TableCell>
                <SecretValueCell
                  slug={slug}
                  keyName={secret.key_name}
                  environment={environment}
                />
              </TableCell>
              <TableCell>
                <Badge variant="outline" className="text-xs tabular-nums">
                  v{secret.vault_version}
                </Badge>
              </TableCell>
              <TableCell className="text-sm text-muted-foreground">
                {formatRelativeTime(secret.last_modified_at)}
              </TableCell>
              <TableCell>
                <AlertDialog>
                  <AlertDialogTrigger asChild>
                    <Button
                      variant="ghost"
                      size="icon"
                      className="h-8 w-8 text-muted-foreground hover:text-destructive"
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </AlertDialogTrigger>
                  <AlertDialogContent>
                    <AlertDialogHeader>
                      <AlertDialogTitle>Delete Secret</AlertDialogTitle>
                      <AlertDialogDescription>
                        Are you sure you want to delete{" "}
                        <code className="font-mono font-semibold">
                          {secret.key_name}
                        </code>
                        ? This action cannot be undone.
                      </AlertDialogDescription>
                    </AlertDialogHeader>
                    <AlertDialogFooter>
                      <AlertDialogCancel>Cancel</AlertDialogCancel>
                      <AlertDialogAction
                        onClick={() => handleDelete(secret.key_name)}
                        className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
                      >
                        Delete
                      </AlertDialogAction>
                    </AlertDialogFooter>
                  </AlertDialogContent>
                </AlertDialog>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
}
