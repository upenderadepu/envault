"use client";

import { useState } from "react";
import { Loader2 } from "lucide-react";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useSetSecret } from "@/hooks/use-secrets";
import { toast } from "sonner";

interface SetSecretDialogProps {
  slug: string;
  defaultEnvironment?: string;
  children: React.ReactNode;
}

const KEY_REGEX = /^[A-Za-z_][A-Za-z0-9_]*$/;

export function SetSecretDialog({
  slug,
  defaultEnvironment = "development",
  children,
}: SetSecretDialogProps) {
  const [open, setOpen] = useState(false);
  const [environment, setEnvironment] = useState(defaultEnvironment);
  const [key, setKey] = useState("");
  const [value, setValue] = useState("");
  const [keyError, setKeyError] = useState("");
  const setSecret = useSetSecret(slug);

  const validateKey = (k: string) => {
    if (!k) {
      setKeyError("");
      return;
    }
    if (!KEY_REGEX.test(k)) {
      setKeyError(
        "Must start with a letter or underscore, and contain only alphanumeric characters and underscores."
      );
    } else {
      setKeyError("");
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!KEY_REGEX.test(key)) return;

    try {
      await setSecret.mutateAsync({ environment, key, value });
      toast.success(`Secret "${key}" saved`);
      setOpen(false);
      setKey("");
      setValue("");
      setKeyError("");
    } catch (err) {
      toast.error(
        err instanceof Error ? err.message : "Failed to save secret"
      );
    }
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>{children}</DialogTrigger>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Add Secret</DialogTitle>
          <DialogDescription>
            Set a new environment variable for this project.
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
            <label className="text-sm font-medium">Key</label>
            <Input
              placeholder="DATABASE_URL"
              value={key}
              onChange={(e) => {
                setKey(e.target.value);
                validateKey(e.target.value);
              }}
              required
              maxLength={256}
              className="font-mono"
              autoFocus
            />
            {keyError && (
              <p className="text-xs text-destructive">{keyError}</p>
            )}
          </div>
          <div className="space-y-2">
            <label className="text-sm font-medium">Value</label>
            <Textarea
              placeholder="Enter secret value..."
              value={value}
              onChange={(e) => setValue(e.target.value)}
              required
              className="font-mono min-h-[80px]"
            />
          </div>
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
              disabled={setSecret.isPending || !!keyError || !key}
            >
              {setSecret.isPending && (
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              )}
              Save Secret
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}
