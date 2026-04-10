"use client";

import { useSearchParams, useRouter, usePathname } from "next/navigation";
import { Plus, Upload } from "lucide-react";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Button } from "@/components/ui/button";
import { SecretsTable } from "@/components/secrets/secrets-table";
import { SetSecretDialog } from "@/components/secrets/set-secret-dialog";
import { BulkImportDialog } from "@/components/secrets/bulk-import-dialog";
import { useSecrets } from "@/hooks/use-secrets";
import { useProjectContext } from "@/components/projects/project-context";

const ENVIRONMENTS = ["development", "staging", "production"] as const;

export default function SecretsPage() {
  const project = useProjectContext();
  const searchParams = useSearchParams();
  const router = useRouter();
  const pathname = usePathname();

  const activeEnv = searchParams.get("env") || "development";

  const { data: secrets, isLoading } = useSecrets(project.slug, activeEnv);

  const setEnv = (env: string) => {
    const params = new URLSearchParams(searchParams.toString());
    params.set("env", env);
    router.replace(`${pathname}?${params.toString()}`);
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Secrets</h1>
          <p className="text-muted-foreground">
            Manage environment variables for {project.name}.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <BulkImportDialog slug={project.slug} defaultEnvironment={activeEnv}>
            <Button variant="outline" className="gap-2">
              <Upload className="h-4 w-4" />
              <span className="hidden sm:inline">Bulk Import</span>
            </Button>
          </BulkImportDialog>
          <SetSecretDialog slug={project.slug} defaultEnvironment={activeEnv}>
            <Button className="gap-2">
              <Plus className="h-4 w-4" />
              <span className="hidden sm:inline">Add Secret</span>
            </Button>
          </SetSecretDialog>
        </div>
      </div>

      <Tabs value={activeEnv} onValueChange={setEnv}>
        <TabsList>
          {ENVIRONMENTS.map((env) => (
            <TabsTrigger key={env} value={env} className="capitalize">
              {env}
            </TabsTrigger>
          ))}
        </TabsList>
        {ENVIRONMENTS.map((env) => (
          <TabsContent key={env} value={env} className="mt-4">
            <SecretsTable
              secrets={activeEnv === env ? secrets : undefined}
              isLoading={activeEnv === env ? isLoading : true}
              slug={project.slug}
              environment={env}
            />
          </TabsContent>
        ))}
      </Tabs>
    </div>
  );
}
