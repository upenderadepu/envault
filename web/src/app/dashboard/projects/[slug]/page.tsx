"use client";

import Link from "next/link";
import { KeyRound, Users, ScrollText, Shield, ArrowRight } from "lucide-react";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { useProjectContext } from "@/components/projects/project-context";
import { useSecrets } from "@/hooks/use-secrets";
import { useMembers } from "@/hooks/use-members";
import { useAuditLogs } from "@/hooks/use-audit";
import { formatRelativeTime, actionLabel } from "@/lib/utils";

export default function ProjectOverviewPage() {
  const project = useProjectContext();
  const { data: devSecrets } = useSecrets(project.slug, "development");
  const { data: stagingSecrets } = useSecrets(project.slug, "staging");
  const { data: prodSecrets } = useSecrets(project.slug, "production");
  const { data: members } = useMembers(project.slug);
  const { data: auditData } = useAuditLogs(project.slug, { limit: 5 });

  const stats = [
    {
      label: "Development Secrets",
      value: devSecrets?.length ?? 0,
      icon: KeyRound,
      href: `/dashboard/projects/${project.slug}/secrets?env=development`,
    },
    {
      label: "Staging Secrets",
      value: stagingSecrets?.length ?? 0,
      icon: KeyRound,
      href: `/dashboard/projects/${project.slug}/secrets?env=staging`,
    },
    {
      label: "Production Secrets",
      value: prodSecrets?.length ?? 0,
      icon: Shield,
      href: `/dashboard/projects/${project.slug}/secrets?env=production`,
    },
    {
      label: "Team Members",
      value: members?.length ?? 0,
      icon: Users,
      href: `/dashboard/projects/${project.slug}/members`,
    },
  ];

  return (
    <div className="space-y-8">
      {/* Header */}
      <div>
        <h1 className="text-2xl font-bold tracking-tight">{project.name}</h1>
        <p className="text-muted-foreground font-mono text-sm">{project.slug}</p>
      </div>

      {/* Stats */}
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {stats.map((stat) => (
          <Link key={stat.label} href={stat.href}>
            <Card className="transition-all hover:border-primary/30 hover:shadow-sm cursor-pointer">
              <CardHeader className="flex flex-row items-center justify-between pb-2">
                <CardTitle className="text-sm font-medium text-muted-foreground">
                  {stat.label}
                </CardTitle>
                <stat.icon className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{stat.value}</div>
              </CardContent>
            </Card>
          </Link>
        ))}
      </div>

      {/* Quick Actions */}
      <div className="grid gap-4 sm:grid-cols-3">
        <Link href={`/dashboard/projects/${project.slug}/secrets`}>
          <Card className="cursor-pointer transition-all hover:border-primary/30 group">
            <CardHeader>
              <CardTitle className="flex items-center gap-2 text-base">
                <KeyRound className="h-4 w-4 text-primary" />
                Manage Secrets
                <ArrowRight className="ml-auto h-4 w-4 opacity-0 group-hover:opacity-100 transition-opacity" />
              </CardTitle>
              <CardDescription>
                View, add, or modify environment variables
              </CardDescription>
            </CardHeader>
          </Card>
        </Link>
        <Link href={`/dashboard/projects/${project.slug}/members`}>
          <Card className="cursor-pointer transition-all hover:border-primary/30 group">
            <CardHeader>
              <CardTitle className="flex items-center gap-2 text-base">
                <Users className="h-4 w-4 text-primary" />
                Team Members
                <ArrowRight className="ml-auto h-4 w-4 opacity-0 group-hover:opacity-100 transition-opacity" />
              </CardTitle>
              <CardDescription>
                Invite members and manage access roles
              </CardDescription>
            </CardHeader>
          </Card>
        </Link>
        <Link href={`/dashboard/projects/${project.slug}/audit`}>
          <Card className="cursor-pointer transition-all hover:border-primary/30 group">
            <CardHeader>
              <CardTitle className="flex items-center gap-2 text-base">
                <ScrollText className="h-4 w-4 text-primary" />
                Audit Log
                <ArrowRight className="ml-auto h-4 w-4 opacity-0 group-hover:opacity-100 transition-opacity" />
              </CardTitle>
              <CardDescription>
                Review all actions and access history
              </CardDescription>
            </CardHeader>
          </Card>
        </Link>
      </div>

      {/* Recent Activity */}
      {auditData && auditData.data.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle className="text-base">Recent Activity</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              {auditData.data.map((log) => (
                <div
                  key={log.id}
                  className="flex items-center justify-between text-sm"
                >
                  <div className="flex items-center gap-3">
                    <Badge variant="secondary" className="text-xs">
                      {actionLabel(log.action)}
                    </Badge>
                    <span className="text-muted-foreground font-mono text-xs truncate max-w-[200px]">
                      {log.resource_path}
                    </span>
                  </div>
                  <span className="text-xs text-muted-foreground">
                    {formatRelativeTime(log.created_at)}
                  </span>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  );
}
