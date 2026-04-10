import { createClient } from "@/lib/supabase/client";
import type {
  Project,
  SecretMetadata,
  SecretValue,
  TeamMember,
  PaginatedAuditResponse,
  CreateProjectResponse,
  AddMemberResponse,
  RotateCredentialsResponse,
} from "@/lib/types";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

class ApiClient {
  private async getToken(): Promise<string> {
    const supabase = createClient();
    const { data } = await supabase.auth.getSession();
    const token = data.session?.access_token;
    if (!token) throw new Error("Not authenticated");
    return token;
  }

  private async request<T>(
    path: string,
    options: RequestInit = {}
  ): Promise<T> {
    const token = await this.getToken();
    const res = await fetch(`${API_URL}${path}`, {
      ...options,
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
        ...options.headers,
      },
    });

    if (!res.ok) {
      const body = await res.json().catch(() => ({ error: "Request failed" }));
      throw new Error(body.error || `Request failed with status ${res.status}`);
    }

    if (res.status === 204) return undefined as T;
    return res.json();
  }

  // Projects
  async listProjects(): Promise<Project[]> {
    return this.request<Project[]>("/api/v1/projects");
  }

  async createProject(name: string): Promise<CreateProjectResponse> {
    return this.request<CreateProjectResponse>("/api/v1/projects", {
      method: "POST",
      body: JSON.stringify({ name }),
    });
  }

  async getProject(slug: string): Promise<Project> {
    return this.request<Project>(`/api/v1/projects/${slug}`);
  }

  async deleteProject(slug: string): Promise<void> {
    return this.request(`/api/v1/projects/${slug}`, { method: "DELETE" });
  }

  // Secrets
  async listSecrets(
    slug: string,
    environment: string = "development"
  ): Promise<SecretMetadata[]> {
    return this.request<SecretMetadata[]>(
      `/api/v1/projects/${slug}/secrets?environment=${environment}`
    );
  }

  async getSecret(
    slug: string,
    key: string,
    environment: string = "development"
  ): Promise<SecretValue> {
    return this.request<SecretValue>(
      `/api/v1/projects/${slug}/secrets/${key}?environment=${environment}`
    );
  }

  async setSecret(
    slug: string,
    data: { environment: string; key: string; value: string }
  ): Promise<SecretMetadata> {
    return this.request<SecretMetadata>(`/api/v1/projects/${slug}/secrets`, {
      method: "POST",
      body: JSON.stringify(data),
    });
  }

  async bulkSetSecrets(
    slug: string,
    data: { environment: string; secrets: Record<string, string> }
  ): Promise<SecretMetadata[]> {
    return this.request<SecretMetadata[]>(
      `/api/v1/projects/${slug}/secrets/bulk`,
      {
        method: "POST",
        body: JSON.stringify(data),
      }
    );
  }

  async deleteSecret(
    slug: string,
    key: string,
    environment: string = "development"
  ): Promise<void> {
    return this.request(
      `/api/v1/projects/${slug}/secrets/${key}?environment=${environment}`,
      { method: "DELETE" }
    );
  }

  // Members
  async listMembers(slug: string): Promise<TeamMember[]> {
    return this.request<TeamMember[]>(`/api/v1/projects/${slug}/members`);
  }

  async addMember(
    slug: string,
    data: { email: string; role: string }
  ): Promise<AddMemberResponse> {
    return this.request<AddMemberResponse>(
      `/api/v1/projects/${slug}/members`,
      {
        method: "POST",
        body: JSON.stringify(data),
      }
    );
  }

  async removeMember(slug: string, memberId: string): Promise<void> {
    return this.request(`/api/v1/projects/${slug}/members/${memberId}`, {
      method: "DELETE",
    });
  }

  async rotateCredentials(slug: string): Promise<RotateCredentialsResponse> {
    return this.request<RotateCredentialsResponse>(
      `/api/v1/projects/${slug}/rotate`,
      { method: "POST" }
    );
  }

  // Audit
  async listAuditLogs(
    slug: string,
    params: { action?: string; limit?: number; offset?: number } = {}
  ): Promise<PaginatedAuditResponse> {
    const searchParams = new URLSearchParams();
    if (params.action) searchParams.set("action", params.action);
    if (params.limit) searchParams.set("limit", String(params.limit));
    if (params.offset) searchParams.set("offset", String(params.offset));
    const qs = searchParams.toString();
    return this.request<PaginatedAuditResponse>(
      `/api/v1/projects/${slug}/audit${qs ? `?${qs}` : ""}`
    );
  }
}

export const api = new ApiClient();
