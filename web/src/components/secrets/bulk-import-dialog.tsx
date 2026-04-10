"use client";

import { useState, useCallback } from "react";
import { Loader2, Upload } from "lucide-react";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Badge } from "@/components/ui/badge";
import { useBulkSetSecrets } from "@/hooks/use-secrets";
import { toast } from "sonner";

interface BulkImportDialogProps {
  slug: string;
  defaultEnvironment?: string;
  children: React.ReactNode;
}

function parseEnvContent(content: string): Record<string, string> {
  const secrets: Record<string, string> = {};
  const lines = content.split("\n");

  for (const line of lines) {
    const trimmed = line.trim();
    if (!trimmed || trimmed.startsWith("#")) continue;

    const eqIndex = trimmed.indexOf("=");
    if (eqIndex === -1) continue;

    const key = trimmed.slice(0, eqIndex).trim();
    let value = trimmed.slice(eqIndex + 1).trim();

    // Remove surrounding quotes
    if (
      (value.startsWith('"') && value.endsWith('"')) ||
      (value.startsWith("'") && value.endsWith("'"))
    ) {
      value = value.slice(1, -1);
    }

    if (key && /^[A-Za-z_][A-Za-z0-9_]*$/.test(key)) {
      secrets[key] = value;
    }
  }

  return secrets;
}

export function BulkImportDialog({
  slug,
  defaultEnvironment = "development",
  children,
}: BulkImportDialogProps) {
  const [open, setOpen] = useState(false);
  const [environment, setEnvironment] = useState(defaultEnvironment);
  const [content, setContent] = useState("");
  const bulkSet = useBulkSetSecrets(slug);

  const parsed = content ? parseEnvContent(content) : {};
  const keyCount = Object.keys(parsed).length;

  const handleFileUpload = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const file = e.target.files?.[0];
      if (!file) return;
      const reader = new FileReader();
      reader.onload = (evt) => {
        setContent(evt.target?.result as string);
      };
      reader.readAsText(file);
    },
    []
  );

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (keyCount === 0) return;

    try {
      await bulkSet.mutateAsync({ environment, secrets: parsed });
      toast.success(`${keyCount} secret${keyCount !== 1 ? "s" : ""} imported`);
      setOpen(false);
      setContent("");
    } catch (err) {
      toast.error(
        err instanceof Error ? err.message : "Failed to import secrets"
      );
    }
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>{children}</DialogTrigger>
      <DialogContent className="sm:max-w-lg">
        <DialogHeader>
          <DialogTitle>Bulk Import Secrets</DialogTitle>
          <DialogDescription>
            Paste .env file content or upload a file.
          </DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <label className="text-sm font-medium">Environment</label>
            <Select value={environment} onValueChange={setEnvironment}>
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="development">Development</SelectItem>
                <SelectItem value="staging">Staging</SelectItem>
                <SelectItem value="production">Production</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-2">
            <div className="flex items-center justify-between">
              <label className="text-sm font-medium">
                Environment Variables
              </label>
              <label className="cursor-pointer">
                <input
                  type="file"
                  accept=".env,.env.*,.txt"
                  className="hidden"
                  onChange={handleFileUpload}
                />
                <span className="inline-flex items-center gap-1.5 text-sm text-primary hover:underline">
                  <Upload className="h-3.5 w-3.5" />
                  Upload file
                </span>
              </label>
            </div>
            <Textarea
              placeholder={`# Paste your .env content\nDATABASE_URL=postgres://...\nREDIS_URL=redis://...\nAPI_KEY=sk_live_...`}
              value={content}
              onChange={(e) => setContent(e.target.value)}
              className="font-mono text-sm min-h-[160px]"
            />
          </div>

          {keyCount > 0 && (
            <div className="rounded-md border bg-muted/50 p-3">
              <p className="text-sm font-medium mb-2">
                Preview ({keyCount} key{keyCount !== 1 ? "s" : ""})
              </p>
              <div className="flex flex-wrap gap-1.5">
                {Object.keys(parsed).map((key) => (
                  <Badge key={key} variant="secondary" className="font-mono text-xs">
                    {key}
                  </Badge>
                ))}
              </div>
            </div>
          )}

          <div className="flex justify-end gap-3">
            <Button
              type="button"
              variant="outline"
              onClick={() => setOpen(false)}
            >
              Cancel
            </Button>
            <Button
              type="submit"
              disabled={bulkSet.isPending || keyCount === 0}
            >
              {bulkSet.isPending && (
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              )}
              Import {keyCount > 0 ? `${keyCount} Secret${keyCount !== 1 ? "s" : ""}` : "Secrets"}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}
