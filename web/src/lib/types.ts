export interface User {
  id: string;
  supabase_uid: string;
  email: string;
  created_at: string;
}

export interface Project {
  id: string;
  name: string;
  slug: string;
  vault_mount_path: string;
  owner_id: string;
  owner?: User;
  environments?: Environment[];
  team_members?: TeamMember[];
  created_at: string;
  updated_at: string;
  deleted_at?: string | null;
}

export interface Environment {
  id: string;
  project_id: string;
  name: "development" | "staging" | "production";
  is_production: boolean;
  created_at: string;
}

export interface SecretMetadata {
  id: string;
  project_id: string;
  environment_id: string;
  environment?: Environment;
  key_name: string;
  vault_path: string;
  created_by_id: string;
  vault_version: number;
  last_modified_at: string;
  created_at: string;
}

export interface SecretValue {
  key: string;
  value: string;
  version: number;
  last_modified_at: string;
}

export interface TeamMember {
  id: string;
  project_id: string;
  user_id: string;
  user?: User;
  role: "admin" | "developer" | "ci";
  vault_policy_name?: string;
  is_active: boolean;
  invited_at: string;
  joined_at?: string;
  created_at: string;
  updated_at: string;
}

export interface AuditLog {
  id: string;
  project_id: string;
  user_id?: string;
  action: string;
  resource_path: string;
  ip_address?: string;
  user_agent?: string;
  request_id?: string;
  metadata: Record<string, unknown>;
  created_at: string;
}

export interface PaginatedAuditResponse {
  data: AuditLog[];
  total: number;
  limit: number;
  offset: number;
}

export interface CreateProjectResponse {
  project: Project;
  vault_token: string;
}

export interface AddMemberResponse {
  member: TeamMember;
  vault_token: string;
}

export interface RotateCredentialsResponse {
  vault_token: string;
}

export interface ApiError {
  error: string;
}
