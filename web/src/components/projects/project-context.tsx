"use client";

import { createContext, useContext } from "react";
import type { Project } from "@/lib/types";

export const ProjectContext = createContext<Project | null>(null);

export function useProjectContext() {
  const ctx = useContext(ProjectContext);
  if (!ctx)
    throw new Error("useProjectContext must be used within a project layout");
  return ctx;
}
