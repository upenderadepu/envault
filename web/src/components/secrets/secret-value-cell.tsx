"use client";

import { useState, useEffect, useCallback } from "react";
import { Eye, EyeOff, Copy, Check, Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { api } from "@/lib/api";
import { toast } from "sonner";

interface SecretValueCellProps {
  slug: string;
  keyName: string;
  environment: string;
}

export function SecretValueCell({
  slug,
  keyName,
  environment,
}: SecretValueCellProps) {
  const [revealed, setRevealed] = useState(false);
  const [value, setValue] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [copied, setCopied] = useState(false);
  const [countdown, setCountdown] = useState(0);

  const hideValue = useCallback(() => {
    setRevealed(false);
    setValue(null);
    setCountdown(0);
  }, []);

  useEffect(() => {
    if (!revealed) return;
    setCountdown(10);
    const interval = setInterval(() => {
      setCountdown((prev) => {
        if (prev <= 1) {
          hideValue();
          return 0;
        }
        return prev - 1;
      });
    }, 1000);
    return () => clearInterval(interval);
  }, [revealed, hideValue]);

  const handleReveal = async () => {
    if (revealed) {
      hideValue();
      return;
    }
    setLoading(true);
    try {
      const secret = await api.getSecret(slug, keyName, environment);
      setValue(secret.value);
      setRevealed(true);
    } catch (err) {
      toast.error(
        err instanceof Error ? err.message : "Failed to reveal secret"
      );
    } finally {
      setLoading(false);
    }
  };

  const handleCopy = async () => {
    if (!value) return;
    await navigator.clipboard.writeText(value);
    setCopied(true);
    toast.success("Copied to clipboard");
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="flex items-center gap-2">
      <span className="font-mono text-sm truncate max-w-[200px]">
        {revealed && value ? value : "••••••••••••"}
      </span>
      <div className="flex items-center gap-1 ml-auto flex-shrink-0">
        {revealed && (
          <span className="text-xs text-muted-foreground tabular-nums mr-1">
            {countdown}s
          </span>
        )}
        <Tooltip>
          <TooltipTrigger asChild>
            <Button
              variant="ghost"
              size="icon"
              className="h-7 w-7"
              onClick={handleReveal}
              disabled={loading}
            >
              {loading ? (
                <Loader2 className="h-3.5 w-3.5 animate-spin" />
              ) : revealed ? (
                <EyeOff className="h-3.5 w-3.5" />
              ) : (
                <Eye className="h-3.5 w-3.5" />
              )}
            </Button>
          </TooltipTrigger>
          <TooltipContent>{revealed ? "Hide" : "Reveal"}</TooltipContent>
        </Tooltip>
        {revealed && (
          <Tooltip>
            <TooltipTrigger asChild>
              <Button
                variant="ghost"
                size="icon"
                className="h-7 w-7"
                onClick={handleCopy}
              >
                {copied ? (
                  <Check className="h-3.5 w-3.5 text-green-500" />
                ) : (
                  <Copy className="h-3.5 w-3.5" />
                )}
              </Button>
            </TooltipTrigger>
            <TooltipContent>Copy</TooltipContent>
          </Tooltip>
        )}
      </div>
    </div>
  );
}
