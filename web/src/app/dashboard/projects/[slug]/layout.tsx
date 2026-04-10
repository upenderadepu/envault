"use client";

import { useProject } from "@/hooks/use-projects";
import { Skeleton } from "@/components/ui/skeleton";
import { ProjectContext } from "@/components/projects/project-context";

export default function ProjectLayout({
  children,
  params,
}: {
  children: React.ReactNode;
  params: { slug: string };
}) {
  const { data: project, isLoading, error } = useProject(params.slug);

  if (isLoading) {
    return (
      <div className="space-y-6">
        <Skeleton className="h-8 w-48" />
        <Skeleton className="h-64 w-full" />
      </div>
    );
  }

  if (error || !project) {
    return (
      <div className="flex flex-col items-center justify-center py-20">
        <h2 className="text-xl font-semibold">Project not found</h2>
        <p className="mt-2 text-muted-foreground">
          The project you&apos;re looking for doesn&apos;t exist or you don&apos;t have access.
        </p>
      </div>
    );
  }

  return (
    <ProjectContext.Provider value={project}>
      {children}
    </ProjectContext.Provider>
  );
}
