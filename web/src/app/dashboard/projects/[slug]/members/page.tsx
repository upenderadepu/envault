"use client";

import { UserPlus } from "lucide-react";
import { Button } from "@/components/ui/button";
import { MembersTable } from "@/components/members/members-table";
import { InviteMemberDialog } from "@/components/members/invite-member-dialog";
import { useMembers } from "@/hooks/use-members";
import { useProjectContext } from "@/components/projects/project-context";

export default function MembersPage() {
  const project = useProjectContext();
  const { data: members, isLoading } = useMembers(project.slug);

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Team Members</h1>
          <p className="text-muted-foreground">
            Manage who has access to {project.name}.
          </p>
        </div>
        <InviteMemberDialog slug={project.slug}>
          <Button className="gap-2">
            <UserPlus className="h-4 w-4" />
            Invite Member
          </Button>
        </InviteMemberDialog>
      </div>

      <MembersTable
        members={members}
        isLoading={isLoading}
        slug={project.slug}
      />
    </div>
  );
}
