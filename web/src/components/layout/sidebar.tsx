"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import {
  LayoutDashboard,
  KeyRound,
  Users,
  ScrollText,
  Settings,
  Shield,
  LogOut,
  ChevronLeft,
  Eye,
} from "lucide-react";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import { ScrollArea } from "@/components/ui/scroll-area";
import { createClient } from "@/lib/supabase/client";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

const mainNav = [
  { label: "Dashboard", href: "/dashboard", icon: LayoutDashboard },
  { label: "Settings", href: "/dashboard/settings", icon: Settings },
];

const projectNav = [
  { label: "Overview", href: "", icon: Eye },
  { label: "Secrets", href: "/secrets", icon: KeyRound },
  { label: "Members", href: "/members", icon: Users },
  { label: "Audit Log", href: "/audit", icon: ScrollText },
  { label: "Settings", href: "/settings", icon: Settings },
];

interface SidebarProps {
  className?: string;
  onClose?: () => void;
}

export function Sidebar({ className, onClose }: SidebarProps) {
  const pathname = usePathname();
  const router = useRouter();
  const [userEmail, setUserEmail] = useState<string>("");

  const slugMatch = pathname.match(/^\/dashboard\/projects\/([^/]+)/);
  const currentSlug = slugMatch ? slugMatch[1] : null;

  useEffect(() => {
    const supabase = createClient();
    supabase.auth.getUser().then(({ data }) => {
      setUserEmail(data.user?.email || "");
    });
  }, []);

  const handleSignOut = async () => {
    const supabase = createClient();
    await supabase.auth.signOut();
    router.push("/login");
  };

  return (
    <div
      className={cn(
        "flex h-full w-[280px] flex-col border-r bg-sidebar text-sidebar-foreground",
        className
      )}
    >
      {/* Logo */}
      <div className="flex h-16 items-center gap-2 px-6">
        <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary">
          <Shield className="h-4 w-4 text-primary-foreground" />
        </div>
        <span className="text-lg font-semibold tracking-tight">Envault</span>
      </div>

      <Separator />

      <ScrollArea className="flex-1 px-3 py-4">
        {/* Main Navigation */}
        <div className="space-y-1">
          {mainNav.map((item) => {
            const isActive =
              item.href === "/dashboard"
                ? pathname === "/dashboard"
                : pathname.startsWith(item.href);
            return (
              <Link key={item.href} href={item.href} onClick={onClose}>
                <Button
                  variant={isActive ? "secondary" : "ghost"}
                  className={cn(
                    "w-full justify-start gap-3 font-medium",
                    isActive && "bg-sidebar-accent text-sidebar-accent-foreground"
                  )}
                  size="sm"
                >
                  <item.icon className="h-4 w-4" />
                  {item.label}
                </Button>
              </Link>
            );
          })}
        </div>

        {/* Project Navigation */}
        {currentSlug && (
          <>
            <Separator className="my-4" />
            <div className="mb-2 flex items-center gap-2 px-3">
              <Link
                href="/dashboard"
                className="text-muted-foreground hover:text-foreground"
                onClick={onClose}
              >
                <ChevronLeft className="h-4 w-4" />
              </Link>
              <span className="truncate text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                {decodeURIComponent(currentSlug)}
              </span>
            </div>
            <div className="space-y-1">
              {projectNav.map((item) => {
                const href = `/dashboard/projects/${currentSlug}${item.href}`;
                const isActive =
                  item.href === ""
                    ? pathname === `/dashboard/projects/${currentSlug}`
                    : pathname.startsWith(href);
                return (
                  <Link key={item.href} href={href} onClick={onClose}>
                    <Button
                      variant={isActive ? "secondary" : "ghost"}
                      className={cn(
                        "w-full justify-start gap-3 font-medium",
                        isActive &&
                          "bg-sidebar-accent text-sidebar-accent-foreground"
                      )}
                      size="sm"
                    >
                      <item.icon className="h-4 w-4" />
                      {item.label}
                    </Button>
                  </Link>
                );
              })}
            </div>
          </>
        )}
      </ScrollArea>

      <Separator />

      {/* User Section */}
      <div className="p-4">
        <div className="flex items-center gap-3">
          <div className="flex h-8 w-8 items-center justify-center rounded-full bg-primary/10 text-primary">
            <span className="text-sm font-medium">
              {userEmail?.[0]?.toUpperCase() || "U"}
            </span>
          </div>
          <div className="flex-1 truncate">
            <p className="truncate text-sm font-medium">{userEmail || "User"}</p>
          </div>
          <Button
            variant="ghost"
            size="icon"
            className="h-8 w-8 text-muted-foreground hover:text-foreground"
            onClick={handleSignOut}
          >
            <LogOut className="h-4 w-4" />
          </Button>
        </div>
      </div>
    </div>
  );
}
