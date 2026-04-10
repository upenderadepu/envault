"use client";

import { Plus, FolderKey } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { useProjects } from "@/hooks/use-projects";
import { ProjectCard } from "@/components/projects/project-card";
import { CreateProjectDialog } from "@/components/projects/create-project-dialog";

export default function DashboardPage() {
  const { data: projects, isLoading } = useProjects();

  return (
    <div>
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Projects</h1>
          <p className="text-muted-foreground">
            Manage your secrets across all projects.
          </p>
        </div>
        <CreateProjectDialog>
          <Button className="gap-2">
            <Plus className="h-4 w-4" />
            New Project
          </Button>
        </CreateProjectDialog>
      </div>

      {isLoading ? (
        <div className="mt-8 grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {[...Array(3)].map((_, i) => (
            <Skeleton key={i} className="h-[140px] rounded-xl" />
          ))}
        </div>
      ) : projects && projects.length > 0 ? (
        <div className="mt-8 grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {projects.map((project) => (
            <ProjectCard key={project.id} project={project} />
          ))}
        </div>
      ) : (
        <div className="mt-16 flex flex-col items-center justify-center text-center">
          <div className="flex h-20 w-20 items-center justify-center rounded-2xl bg-muted">
            <FolderKey className="h-10 w-10 text-muted-foreground" />
          </div>
          <h3 className="mt-6 text-lg font-semibold">No projects yet</h3>
          <p className="mt-2 max-w-sm text-sm text-muted-foreground">
            Create your first project to start managing secrets securely with
            HashiCorp Vault.
          </p>
          <CreateProjectDialog>
            <Button className="mt-6 gap-2">
              <Plus className="h-4 w-4" />
              Create Your First Project
            </Button>
          </CreateProjectDialog>
        </div>
      )}
    </div>
  );
}
