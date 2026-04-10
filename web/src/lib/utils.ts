import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function formatDate(iso: string): string {
  return new Date(iso).toLocaleDateString("en-US", {
    year: "numeric",
    month: "short",
    day: "numeric",
  });
}

export function formatDateTime(iso: string): string {
  return new Date(iso).toLocaleString("en-US", {
    year: "numeric",
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

export function formatRelativeTime(iso: string): string {
  const now = Date.now();
  const then = new Date(iso).getTime();
  const diff = now - then;
  const seconds = Math.floor(diff / 1000);
  const minutes = Math.floor(seconds / 60);
  const hours = Math.floor(minutes / 60);
  const days = Math.floor(hours / 24);

  if (seconds < 60) return "just now";
  if (minutes < 60) return `${minutes}m ago`;
  if (hours < 24) return `${hours}h ago`;
  if (days < 30) return `${days}d ago`;
  return formatDate(iso);
}

export function roleColor(role: string): "default" | "secondary" | "destructive" | "outline" {
  switch (role) {
    case "admin":
      return "default";
    case "developer":
      return "secondary";
    case "ci":
      return "outline";
    default:
      return "secondary";
  }
}

export function actionLabel(action: string): string {
  const labels: Record<string, string> = {
    "project.create": "Project Created",
    "project.delete": "Project Deleted",
    "secret.read": "Secret Read",
    "secret.write": "Secret Updated",
    "secret.delete": "Secret Deleted",
    "member.invite": "Member Invited",
    "member.remove": "Member Removed",
    "credentials.rotate": "Credentials Rotated",
  };
  return labels[action] || action;
}
