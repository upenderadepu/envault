"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import {
  Loader2,
  Copy,
  Check,
  AlertTriangle,
  RefreshCw,
  Trash2,
} from "lucide-react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
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
import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";
import { useProjectContext } from "@/components/projects/project-context";
import { useDeleteProject } from "@/hooks/use-projects";
import { useRotateCredentials } from "@/hooks/use-members";
import { toast } from "sonner";

export default function ProjectSettingsPage() {
  const project = useProjectContext();
  const router = useRouter();
  const deleteProject = useDeleteProject();
  const rotateCredentials = useRotateCredentials(project.slug);

  const [confirmName, setConfirmName] = useState("");
  const [newToken, setNewToken] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);

  const handleRotate = async () => {
    try {
      const result = await rotateCredentials.mutateAsync();
      setNewToken(result.vault_token);
      toast.success("Credentials rotated");
    } catch (err) {
      toast.error(
        err instanceof Error ? err.message : "Failed to rotate credentials"
      );
    }
  };

  const handleCopy = async () => {
    if (newToken) {
      await navigator.clipboard.writeText(newToken);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    }
  };

  const handleDelete = async () => {
    try {
      await deleteProject.mutateAsync(project.slug);
      toast.success("Project deleted");
      router.push("/dashboard");
    } catch (err) {
      toast.error(
        err instanceof Error ? err.message : "Failed to delete project"
      );
    }
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold tracking-tight">Project Settings</h1>
        <p className="text-muted-foreground">
          Manage settings for {project.name}.
        </p>
      </div>

      {/* Project Info */}
      <Card>
        <CardHeader>
          <CardTitle className="text-base">Project Information</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid gap-4 sm:grid-cols-2">
            <div>
              <label className="text-sm font-medium text-muted-foreground">
                Name
              </label>
              <p className="mt-1 font-medium">{project.name}</p>
            </div>
            <div>
              <label className="text-sm font-medium text-muted-foreground">
                Slug
              </label>
              <p className="mt-1 font-mono text-sm">{project.slug}</p>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Rotate Credentials */}
      <Card>
        <CardHeader>
          <CardTitle className="text-base">Rotate Credentials</CardTitle>
          <CardDescription>
            Generate a new Vault token. The old token will be revoked
            immediately.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          {newToken && (
            <div className="space-y-3">
              <div className="flex items-start gap-2 rounded-lg border border-yellow-500/30 bg-yellow-500/5 p-3">
                <AlertTriangle className="mt-0.5 h-4 w-4 text-yellow-600 dark:text-yellow-400 flex-shrink-0" />
                <p className="text-sm text-yellow-700 dark:text-yellow-300">
                  Save this token. It will not be shown again.
                </p>
              </div>
              <div className="flex items-center gap-2">
                <code className="flex-1 rounded-md bg-muted px-3 py-2 font-mono text-xs break-all">
                  {newToken}
                </code>
                <Button variant="outline" size="icon" onClick={handleCopy}>
                  {copied ? (
                    <Check className="h-4 w-4 text-green-500" />
                  ) : (
                    <Copy className="h-4 w-4" />
                  )}
                </Button>
              </div>
            </div>
          )}
          <Button
            variant="outline"
            onClick={handleRotate}
            disabled={rotateCredentials.isPending}
            className="gap-2"
          >
            {rotateCredentials.isPending ? (
              <Loader2 className="h-4 w-4 animate-spin" />
            ) : (
              <RefreshCw className="h-4 w-4" />
            )}
            Rotate Credentials
          </Button>
        </CardContent>
      </Card>

      <Separator />

      {/* Danger Zone */}
      <Card className="border-destructive/30">
        <CardHeader>
          <CardTitle className="text-base text-destructive">
            Danger Zone
          </CardTitle>
          <CardDescription>
            Permanently delete this project and all its secrets.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <AlertDialog>
            <AlertDialogTrigger asChild>
              <Button variant="destructive" className="gap-2">
                <Trash2 className="h-4 w-4" />
                Delete Project
              </Button>
            </AlertDialogTrigger>
            <AlertDialogContent>
              <AlertDialogHeader>
                <AlertDialogTitle>Delete Project</AlertDialogTitle>
                <AlertDialogDescription>
                  This action cannot be undone. All secrets, environments, and
                  team member access will be permanently removed. Type{" "}
                  <code className="font-mono font-semibold">
                    {project.slug}
                  </code>{" "}
                  to confirm.
                </AlertDialogDescription>
              </AlertDialogHeader>
              <Input
                placeholder={project.slug}
                value={confirmName}
                onChange={(e) => setConfirmName(e.target.value)}
                className="font-mono"
              />
              <AlertDialogFooter>
                <AlertDialogCancel onClick={() => setConfirmName("")}>
                  Cancel
                </AlertDialogCancel>
                <AlertDialogAction
                  onClick={handleDelete}
                  disabled={
                    confirmName !== project.slug || deleteProject.isPending
                  }
                  className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
                >
                  {deleteProject.isPending && (
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  )}
                  Delete Permanently
                </AlertDialogAction>
              </AlertDialogFooter>
            </AlertDialogContent>
          </AlertDialog>
        </CardContent>
      </Card>
    </div>
  );
}
