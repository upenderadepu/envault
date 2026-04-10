"use client";

import Link from "next/link";
import { FolderKey, Users, Clock } from "lucide-react";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { formatRelativeTime } from "@/lib/utils";
import type { Project } from "@/lib/types";

interface ProjectCardProps {
  project: Project;
}

export function ProjectCard({ project }: ProjectCardProps) {
  return (
    <Link href={`/dashboard/projects/${project.slug}`}>
      <Card className="group cursor-pointer transition-all hover:border-primary/30 hover:shadow-md hover:shadow-primary/5">
        <CardHeader className="pb-3">
          <div className="flex items-start justify-between">
            <div className="flex items-center gap-3">
              <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10 text-primary group-hover:bg-primary/15 transition-colors">
                <FolderKey className="h-5 w-5" />
              </div>
              <div>
                <CardTitle className="text-base">{project.name}</CardTitle>
                <p className="text-sm text-muted-foreground font-mono">
                  {project.slug}
                </p>
              </div>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <div className="flex items-center gap-4 text-sm text-muted-foreground">
            {project.environments && (
              <div className="flex items-center gap-1.5">
                <Badge variant="secondary" className="text-xs">
                  {project.environments.length} env
                  {project.environments.length !== 1 ? "s" : ""}
                </Badge>
              </div>
            )}
            {project.team_members && (
              <div className="flex items-center gap-1.5">
                <Users className="h-3.5 w-3.5" />
                <span>
                  {project.team_members.length} member
                  {project.team_members.length !== 1 ? "s" : ""}
                </span>
              </div>
            )}
            <div className="ml-auto flex items-center gap-1.5">
              <Clock className="h-3.5 w-3.5" />
              <span>{formatRelativeTime(project.updated_at)}</span>
            </div>
          </div>
        </CardContent>
      </Card>
    </Link>
  );
}
