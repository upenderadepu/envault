"use client";

import { useState } from "react";
import { Loader2, Copy, Check, AlertTriangle } from "lucide-react";
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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useAddMember } from "@/hooks/use-members";
import { toast } from "sonner";

interface InviteMemberDialogProps {
  slug: string;
  children: React.ReactNode;
}

export function InviteMemberDialog({
  slug,
  children,
}: InviteMemberDialogProps) {
  const [open, setOpen] = useState(false);
  const [email, setEmail] = useState("");
  const [role, setRole] = useState("developer");
  const [vaultToken, setVaultToken] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);
  const addMember = useAddMember(slug);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const result = await addMember.mutateAsync({ email, role });
      setVaultToken(result.vault_token);
      toast.success("Member invited successfully");
    } catch (err) {
      toast.error(
        err instanceof Error ? err.message : "Failed to invite member"
      );
    }
  };

  const handleCopy = async () => {
    if (vaultToken) {
      await navigator.clipboard.writeText(vaultToken);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    }
  };

  const handleClose = () => {
    setOpen(false);
    setEmail("");
    setRole("developer");
    setVaultToken(null);
    setCopied(false);
  };

  return (
    <Dialog open={open} onOpenChange={(o) => (o ? setOpen(true) : handleClose())}>
      <DialogTrigger asChild>{children}</DialogTrigger>
      <DialogContent className="sm:max-w-md">
        {!vaultToken ? (
          <>
            <DialogHeader>
              <DialogTitle>Invite Member</DialogTitle>
              <DialogDescription>
                Add a team member and assign their access role.
              </DialogDescription>
            </DialogHeader>
            <form onSubmit={handleSubmit} className="space-y-4">
              <div className="space-y-2">
                <label className="text-sm font-medium">Email</label>
                <Input
                  type="email"
                  placeholder="teammate@company.com"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  required
                  autoFocus
                />
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium">Role</label>
                <Select value={role} onValueChange={setRole}>
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="admin">
                      <div>
                        <span className="font-medium">Admin</span>
                        <span className="text-muted-foreground ml-2 text-xs">
                          Full access
                        </span>
                      </div>
                    </SelectItem>
                    <SelectItem value="developer">
                      <div>
                        <span className="font-medium">Developer</span>
                        <span className="text-muted-foreground ml-2 text-xs">
                          Read & write secrets
                        </span>
                      </div>
                    </SelectItem>
                    <SelectItem value="ci">
                      <div>
                        <span className="font-medium">CI</span>
                        <span className="text-muted-foreground ml-2 text-xs">
                          Read-only
                        </span>
                      </div>
                    </SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div className="flex justify-end gap-3">
                <Button type="button" variant="outline" onClick={handleClose}>
                  Cancel
                </Button>
                <Button type="submit" disabled={addMember.isPending}>
                  {addMember.isPending && (
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  )}
                  Invite
                </Button>
              </div>
            </form>
          </>
        ) : (
          <>
            <DialogHeader>
              <DialogTitle>Member Invited</DialogTitle>
              <DialogDescription>
                Share this Vault token securely with the new member.
              </DialogDescription>
            </DialogHeader>
            <div className="space-y-4">
              <div className="flex items-start gap-2 rounded-lg border border-yellow-500/30 bg-yellow-500/5 p-3">
                <AlertTriangle className="mt-0.5 h-4 w-4 text-yellow-600 dark:text-yellow-400 flex-shrink-0" />
                <p className="text-sm text-yellow-700 dark:text-yellow-300">
                  This token will not be shown again. Share it securely with the
                  invited member.
                </p>
              </div>
              <div className="flex items-center gap-2">
                <code className="flex-1 rounded-md bg-muted px-3 py-2 font-mono text-xs break-all">
                  {vaultToken}
                </code>
                <Button variant="outline" size="icon" onClick={handleCopy}>
                  {copied ? (
                    <Check className="h-4 w-4 text-green-500" />
                  ) : (
                    <Copy className="h-4 w-4" />
                  )}
                </Button>
              </div>
              <div className="flex justify-end">
                <Button onClick={handleClose}>Done</Button>
              </div>
            </div>
          </>
        )}
      </DialogContent>
    </Dialog>
  );
}
