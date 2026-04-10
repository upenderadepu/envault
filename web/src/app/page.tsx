"use client";

import Link from "next/link";
import { useTheme } from "next-themes";
import { useEffect, useState } from "react";
import {
  Shield,
  Lock,
  Users,
  GitBranch,
  Terminal,
  Activity,
  ArrowRight,
  Eye,
  EyeOff,
  ChevronRight,
  Fingerprint,
  ShieldCheck,
  Layers,
  KeyRound,
  Moon,
  Sun,
  Zap,
  BarChart3,
  Clock,
} from "lucide-react";
import { Button } from "@/components/ui/button";

/* ─── Navbar ─── */
function Navbar() {
  const [scrolled, setScrolled] = useState(false);
  const { theme, setTheme } = useTheme();
  const [mounted, setMounted] = useState(false);

  useEffect(() => { setMounted(true); }, []);
  useEffect(() => {
    const handler = () => setScrolled(window.scrollY > 20);
    window.addEventListener("scroll", handler);
    return () => window.removeEventListener("scroll", handler);
  }, []);

  return (
    <nav
      className={`fixed top-0 z-50 w-full transition-all duration-300 ${
        scrolled
          ? "border-b border-border/50 bg-background/80 backdrop-blur-xl"
          : "bg-transparent"
      }`}
    >
      <div className="mx-auto flex h-16 max-w-6xl items-center justify-between px-6">
        <Link href="/" className="flex items-center gap-2.5">
          <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary">
            <Shield className="h-4 w-4 text-primary-foreground" />
          </div>
          <span className="text-lg font-bold tracking-tight">Envault</span>
        </Link>

        <div className="hidden items-center gap-8 md:flex">
          <a href="#features" className="text-sm text-muted-foreground transition-colors hover:text-foreground">Features</a>
          <a href="#how-it-works" className="text-sm text-muted-foreground transition-colors hover:text-foreground">How it Works</a>
          <a href="#security" className="text-sm text-muted-foreground transition-colors hover:text-foreground">Security</a>
          <a href="#cli" className="text-sm text-muted-foreground transition-colors hover:text-foreground">CLI</a>
        </div>

        <div className="flex items-center gap-3">
          {mounted && (
            <Button
              variant="ghost"
              size="icon"
              className="h-9 w-9"
              onClick={() => setTheme(theme === "dark" ? "light" : "dark")}
            >
              {theme === "dark" ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
            </Button>
          )}
          <Link href="/login">
            <Button variant="ghost" size="sm">Sign In</Button>
          </Link>
          <Link href="/login">
            <Button size="sm" className="gap-1.5">
              Get Started <ArrowRight className="h-3.5 w-3.5" />
            </Button>
          </Link>
        </div>
      </div>
    </nav>
  );
}

/* ─── Hero Visual — Animated Vault Graphic ─── */
function HeroVisual() {
  const [revealed, setRevealed] = useState(false);

  return (
    <div className="relative mx-auto w-full max-w-lg">
      {/* Glow behind */}
      <div className="absolute inset-0 -z-10 rounded-3xl bg-primary/10 blur-3xl animate-pulse-glow" />

      {/* Main card */}
      <div className="relative overflow-hidden rounded-2xl border border-border/60 bg-card/80 shadow-2xl shadow-primary/5 backdrop-blur-sm">
        {/* Title bar */}
        <div className="flex items-center gap-2 border-b border-border/50 px-5 py-3.5">
          <div className="flex gap-1.5">
            <div className="h-3 w-3 rounded-full bg-red-400/80" />
            <div className="h-3 w-3 rounded-full bg-yellow-400/80" />
            <div className="h-3 w-3 rounded-full bg-green-400/80" />
          </div>
          <span className="ml-2 text-xs text-muted-foreground font-mono">my-project / production</span>
        </div>

        {/* Secret rows */}
        <div className="divide-y divide-border/30 p-1">
          {[
            { key: "DATABASE_URL", value: "postgresql://prod:s3cur3@db.internal:5432/app" },
            { key: "STRIPE_SECRET_KEY", value: "sk_live_51H7a...redacted...xG4k" },
            { key: "JWT_SIGNING_KEY", value: "eyJhbGciOiJSUzI1NiIsInR5cCI6Ikp..." },
            { key: "AWS_SECRET_ACCESS_KEY", value: "wJalrXUtnFEMI/K7MDENG/bPxRfi..." },
          ].map((secret, i) => (
            <div key={secret.key} className="flex items-center justify-between px-4 py-3 transition-colors hover:bg-muted/30" style={{ animationDelay: `${i * 100}ms` }}>
              <div className="flex items-center gap-3">
                <KeyRound className="h-3.5 w-3.5 text-primary/60" />
                <span className="font-mono text-sm font-medium">{secret.key}</span>
              </div>
              <div className="flex items-center gap-2">
                <span className="font-mono text-xs text-muted-foreground">
                  {revealed ? secret.value : "\u2022\u2022\u2022\u2022\u2022\u2022\u2022\u2022\u2022\u2022\u2022\u2022\u2022\u2022\u2022\u2022\u2022\u2022\u2022\u2022"}
                </span>
                <button
                  onClick={() => setRevealed(!revealed)}
                  className="rounded-md p-1 text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
                >
                  {revealed ? <EyeOff className="h-3.5 w-3.5" /> : <Eye className="h-3.5 w-3.5" />}
                </button>
              </div>
            </div>
          ))}
        </div>

        {/* Status bar */}
        <div className="flex items-center justify-between border-t border-border/50 px-5 py-2.5">
          <div className="flex items-center gap-2">
            <div className="h-2 w-2 rounded-full bg-green-500 animate-pulse" />
            <span className="text-xs text-muted-foreground">Vault connected</span>
          </div>
          <div className="flex items-center gap-4 text-xs text-muted-foreground">
            <span className="flex items-center gap-1"><Lock className="h-3 w-3" /> AES-256</span>
            <span className="flex items-center gap-1"><Activity className="h-3 w-3" /> 4 secrets</span>
          </div>
        </div>
      </div>

      {/* Orbiting elements */}
      <div className="absolute -right-4 -top-4 animate-orbit">
        <div className="flex h-10 w-10 items-center justify-center rounded-xl border border-border/60 bg-card shadow-lg">
          <ShieldCheck className="h-5 w-5 text-primary" />
        </div>
      </div>
      <div className="absolute -bottom-4 -left-4 animate-orbit-reverse">
        <div className="flex h-10 w-10 items-center justify-center rounded-xl border border-border/60 bg-card shadow-lg">
          <Fingerprint className="h-5 w-5 text-primary" />
        </div>
      </div>
    </div>
  );
}

/* ─── Feature Card ─── */
function FeatureCard({ icon: Icon, title, description }: {
  icon: React.ElementType;
  title: string;
  description: string;
}) {
  return (
    <div className="group relative rounded-2xl border border-border/50 bg-card/50 p-6 transition-all duration-300 hover:border-primary/30 hover:bg-card hover:shadow-lg hover:shadow-primary/5">
      <div className="mb-4 flex h-11 w-11 items-center justify-center rounded-xl bg-primary/10 text-primary transition-colors group-hover:bg-primary group-hover:text-primary-foreground">
        <Icon className="h-5 w-5" />
      </div>
      <h3 className="mb-2 text-base font-semibold">{title}</h3>
      <p className="text-sm leading-relaxed text-muted-foreground">{description}</p>
    </div>
  );
}

/* ─── Step Card ─── */
function StepCard({ step, title, description, code }: {
  step: number;
  title: string;
  description: string;
  code: string;
}) {
  return (
    <div className="relative flex flex-col items-center text-center">
      <div className="mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-primary text-lg font-bold text-primary-foreground shadow-lg shadow-primary/25">
        {step}
      </div>
      <h3 className="mb-2 text-base font-semibold">{title}</h3>
      <p className="mb-4 text-sm text-muted-foreground">{description}</p>
      <div className="w-full rounded-xl border border-border/50 bg-muted/30 px-4 py-3">
        <code className="font-mono text-xs text-primary">{code}</code>
      </div>
    </div>
  );
}

/* ─── CLI Preview ─── */
function CLIPreview() {
  return (
    <div className="mx-auto max-w-2xl overflow-hidden rounded-2xl border border-border/60 bg-card shadow-2xl shadow-black/10">
      {/* Title bar */}
      <div className="flex items-center gap-2 border-b border-border/50 bg-muted/30 px-5 py-3">
        <div className="flex gap-1.5">
          <div className="h-3 w-3 rounded-full bg-red-400/80" />
          <div className="h-3 w-3 rounded-full bg-yellow-400/80" />
          <div className="h-3 w-3 rounded-full bg-green-400/80" />
        </div>
        <span className="ml-2 text-xs text-muted-foreground font-mono">terminal</span>
      </div>

      {/* Terminal content */}
      <div className="space-y-4 p-5 font-mono text-sm">
        <div>
          <span className="text-primary">$</span>
          <span className="ml-2 text-foreground">envault init my-saas-app</span>
        </div>
        <div className="pl-2 text-muted-foreground">
          <div>Project &quot;my-saas-app&quot; created.</div>
          <div className="mt-1 text-green-500">Vault Token: hvs.CAES...k2Nz</div>
          <div className="text-xs text-muted-foreground/60 mt-1">Config saved to ~/.envault.yaml</div>
        </div>

        <div className="border-t border-border/30 pt-4">
          <span className="text-primary">$</span>
          <span className="ml-2 text-foreground">envault secret set DATABASE_URL=postgresql://... --env production</span>
        </div>
        <div className="pl-2 text-muted-foreground">
          <div>Secret DATABASE_URL set (version 1) in my-saas-app/production</div>
        </div>

        <div className="border-t border-border/30 pt-4">
          <span className="text-primary">$</span>
          <span className="ml-2 text-foreground">envault env pull --env production</span>
        </div>
        <div className="pl-2 text-muted-foreground">
          <div>Pulled 12 secrets to .env</div>
          <div className="mt-1 text-xs text-muted-foreground/60">DATABASE_URL, STRIPE_SECRET_KEY, JWT_SIGNING_KEY, ...</div>
        </div>

        <div className="border-t border-border/30 pt-4">
          <span className="text-primary">$</span>
          <span className="ml-2 text-foreground">envault onboard alice@team.com --role developer</span>
        </div>
        <div className="pl-2 text-muted-foreground">
          <div>Added alice@team.com as developer to my-saas-app</div>
          <div className="mt-1 text-green-500">Vault Token: hvs.CAES...x9Fp</div>
        </div>

        <div className="flex items-center pt-2">
          <span className="text-primary">$</span>
          <span className="ml-2 inline-block h-4 w-1.5 animate-pulse bg-primary" />
        </div>
      </div>
    </div>
  );
}

/* ─── Stat ─── */
function Stat({ icon: Icon, value, label }: { icon: React.ElementType; value: string; label: string }) {
  return (
    <div className="flex flex-col items-center gap-2 text-center">
      <Icon className="h-5 w-5 text-primary" />
      <div className="text-3xl font-bold tracking-tight">{value}</div>
      <div className="text-sm text-muted-foreground">{label}</div>
    </div>
  );
}

/* ─── Main Page ─── */
export default function LandingPage() {
  return (
    <div className="min-h-screen overflow-hidden">
      <Navbar />

      {/* ─── HERO ─── */}
      <section className="relative pt-32 pb-20 lg:pt-40 lg:pb-32">
        {/* Background grid */}
        <div className="pointer-events-none absolute inset-0 overflow-hidden">
          <div className="absolute inset-0 bg-[linear-gradient(to_right,hsl(var(--border)/0.3)_1px,transparent_1px),linear-gradient(to_bottom,hsl(var(--border)/0.3)_1px,transparent_1px)] bg-[size:64px_64px] [mask-image:radial-gradient(ellipse_80%_50%_at_50%_0%,#000_60%,transparent_100%)]" />
        </div>

        {/* Gradient blobs */}
        <div className="pointer-events-none absolute left-1/4 top-20 h-96 w-96 rounded-full bg-primary/8 blur-3xl" />
        <div className="pointer-events-none absolute right-1/4 bottom-20 h-80 w-80 rounded-full bg-primary/5 blur-3xl" />

        <div className="relative mx-auto max-w-6xl px-6">
          <div className="flex flex-col items-center gap-16 lg:flex-row lg:gap-20">
            {/* Left column — text */}
            <div className="animate-slide-up lg:flex-1">
              <div className="mb-6 inline-flex items-center gap-2 rounded-full border border-primary/20 bg-primary/5 px-4 py-1.5 text-xs font-medium text-primary">
                <ShieldCheck className="h-3.5 w-3.5" />
                Powered by HashiCorp Vault
              </div>

              <h1 className="text-4xl font-extrabold leading-[1.1] tracking-tight sm:text-5xl lg:text-6xl">
                Secrets that stay{" "}
                <span className="bg-gradient-to-r from-primary via-primary to-primary/60 bg-clip-text text-transparent">
                  secret.
                </span>
              </h1>

              <p className="mt-6 max-w-lg text-lg leading-relaxed text-muted-foreground">
                Envault gives your team a single, audited place to store, share, and inject environment variables — backed by HashiCorp Vault, controlled by roles, and observable from a modern dashboard.
              </p>

              <div className="mt-8 flex flex-col gap-3 sm:flex-row sm:items-center">
                <Link href="/login">
                  <Button size="lg" className="gap-2 px-8 shadow-lg shadow-primary/20">
                    Start for Free <ArrowRight className="h-4 w-4" />
                  </Button>
                </Link>
                <a href="#cli">
                  <Button variant="outline" size="lg" className="gap-2 px-8">
                    <Terminal className="h-4 w-4" /> View CLI
                  </Button>
                </a>
              </div>

              <div className="mt-10 flex items-center gap-6 text-sm text-muted-foreground">
                <span className="flex items-center gap-1.5"><Lock className="h-3.5 w-3.5 text-green-500" /> End-to-end encrypted</span>
                <span className="flex items-center gap-1.5"><Users className="h-3.5 w-3.5 text-primary" /> Team-ready</span>
                <span className="flex items-center gap-1.5"><Activity className="h-3.5 w-3.5 text-orange-500" /> Full audit trail</span>
              </div>
            </div>

            {/* Right column — visual (hidden below lg) */}
            <div className="animate-slide-up-delayed lg:flex-1 hero-visual">
              <HeroVisual />
            </div>
          </div>
        </div>
      </section>

      {/* ─── TRUSTED BY / STATS ─── */}
      <section className="border-y border-border/50 bg-muted/20 py-16">
        <div className="mx-auto max-w-4xl px-6">
          <div className="grid grid-cols-2 gap-8 md:grid-cols-4">
            <Stat icon={Shield} value="AES-256" label="Encryption standard" />
            <Stat icon={Zap} value="<50ms" label="Secret retrieval" />
            <Stat icon={Users} value="3" label="RBAC roles" />
            <Stat icon={Layers} value="3" label="Environments" />
          </div>
        </div>
      </section>

      {/* ─── FEATURES ─── */}
      <section id="features" className="py-24 lg:py-32">
        <div className="mx-auto max-w-6xl px-6">
          <div className="mx-auto mb-16 max-w-2xl text-center">
            <h2 className="text-3xl font-bold tracking-tight sm:text-4xl">
              Everything you need for secrets management
            </h2>
            <p className="mt-4 text-lg text-muted-foreground">
              A complete platform that replaces scattered .env files with a secure, audited, team-friendly workflow.
            </p>
          </div>

          <div className="grid gap-5 sm:grid-cols-2 lg:grid-cols-3">
            <FeatureCard
              icon={Lock}
              title="Vault-Backed Storage"
              description="Secret values are stored exclusively in HashiCorp Vault's KV-v2 engine. They never touch your metadata database — ever."
            />
            <FeatureCard
              icon={Users}
              title="Role-Based Access"
              description="Three granular roles — Admin, Developer, and CI — control who can read, write, or manage secrets per project."
            />
            <FeatureCard
              icon={GitBranch}
              title="Multi-Environment"
              description="Development, staging, and production environments are isolated by default. Pull the right secrets for the right context."
            />
            <FeatureCard
              icon={Activity}
              title="Complete Audit Trail"
              description="Every secret read, write, and deletion is logged with user identity, timestamp, and metadata. Full compliance visibility."
            />
            <FeatureCard
              icon={Terminal}
              title="Powerful CLI"
              description="Initialize projects, push/pull .env files, set individual secrets, and onboard teammates — all from your terminal."
            />
            <FeatureCard
              icon={BarChart3}
              title="Prometheus Metrics"
              description="Built-in observability with request counters, latency histograms, and Vault operation tracking out of the box."
            />
            <FeatureCard
              icon={Eye}
              title="Reveal on Demand"
              description="Secret values are masked by default in the dashboard. Reveal them on-demand with an automatic 10-second auto-hide timer."
            />
            <FeatureCard
              icon={KeyRound}
              title="Credential Rotation"
              description="Rotate project Vault tokens with one click or command. All previous tokens are immediately invalidated."
            />
            <FeatureCard
              icon={Clock}
              title="Version Tracking"
              description="Every secret update increments the version counter via Vault's native versioning. See when each secret was last modified."
            />
          </div>
        </div>
      </section>

      {/* ─── HOW IT WORKS ─── */}
      <section id="how-it-works" className="border-y border-border/50 bg-muted/20 py-24 lg:py-32">
        <div className="mx-auto max-w-5xl px-6">
          <div className="mx-auto mb-16 max-w-2xl text-center">
            <h2 className="text-3xl font-bold tracking-tight sm:text-4xl">
              Up and running in minutes
            </h2>
            <p className="mt-4 text-lg text-muted-foreground">
              Four steps from zero to a fully secured secrets pipeline for your team.
            </p>
          </div>

          <div className="grid gap-10 sm:grid-cols-2 lg:grid-cols-4">
            <StepCard
              step={1}
              title="Initialize"
              description="Create a project and receive your Vault token."
              code="envault init my-app"
            />
            <StepCard
              step={2}
              title="Add Secrets"
              description="Push your .env file or set secrets one by one."
              code="envault env push --env prod -f .env"
            />
            <StepCard
              step={3}
              title="Invite Team"
              description="Onboard teammates with scoped roles."
              code="envault onboard dev@co.com --role developer"
            />
            <StepCard
              step={4}
              title="Pull & Deploy"
              description="Pull secrets into any environment or CI pipeline."
              code="envault env pull --env prod"
            />
          </div>
        </div>
      </section>

      {/* ─── SECURITY ─── */}
      <section id="security" className="py-24 lg:py-32">
        <div className="mx-auto max-w-6xl px-6">
          <div className="grid items-center gap-16 lg:grid-cols-2">
            {/* Left — text */}
            <div>
              <div className="mb-4 inline-flex items-center gap-2 rounded-full border border-green-500/20 bg-green-500/5 px-4 py-1.5 text-xs font-medium text-green-600 dark:text-green-400">
                <ShieldCheck className="h-3.5 w-3.5" />
                Security-first architecture
              </div>
              <h2 className="text-3xl font-bold tracking-tight sm:text-4xl">
                Built for teams that can&apos;t afford leaks
              </h2>
              <p className="mt-4 text-muted-foreground leading-relaxed">
                Every architectural decision in Envault prioritizes security. Secret values are completely isolated from metadata, access is scoped by role, and every action leaves an audit trail.
              </p>

              <div className="mt-8 space-y-5">
                {[
                  { icon: Fingerprint, title: "Supabase JWT Auth", desc: "Every request is validated against your Supabase JWKS endpoint. No session cookies to steal." },
                  { icon: Layers, title: "Split Data Model", desc: "PostgreSQL stores metadata only. Vault stores values. A database breach never exposes secrets." },
                  { icon: ShieldCheck, title: "RBAC Enforcement", desc: "Middleware enforces roles on every request. CI tokens can only read — they can never write or manage." },
                  { icon: Activity, title: "Immutable Audit Log", desc: "Who read what, when, from where. Every secret access is recorded with full user context." },
                ].map(({ icon: Icon, title, desc }) => (
                  <div key={title} className="flex gap-4">
                    <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-primary/10">
                      <Icon className="h-5 w-5 text-primary" />
                    </div>
                    <div>
                      <h3 className="text-sm font-semibold">{title}</h3>
                      <p className="text-sm text-muted-foreground">{desc}</p>
                    </div>
                  </div>
                ))}
              </div>
            </div>

            {/* Right — architecture diagram */}
            <div className="relative rounded-2xl border border-border/60 bg-card/80 p-8 shadow-xl">
              <h3 className="mb-6 text-center text-sm font-semibold text-muted-foreground uppercase tracking-wider">Architecture</h3>
              <div className="space-y-4">
                {/* Dashboard layer */}
                <div className="rounded-xl border border-primary/20 bg-primary/5 p-4 text-center">
                  <div className="text-xs font-medium text-primary uppercase tracking-wider mb-1">Dashboard & CLI</div>
                  <div className="text-sm text-muted-foreground">Next.js + Cobra CLI</div>
                </div>
                <div className="flex justify-center">
                  <ChevronRight className="h-5 w-5 rotate-90 text-muted-foreground/40" />
                </div>
                {/* API layer */}
                <div className="rounded-xl border border-border/60 bg-muted/30 p-4 text-center">
                  <div className="text-xs font-medium text-foreground uppercase tracking-wider mb-1">Go API Server</div>
                  <div className="text-sm text-muted-foreground">JWT Auth &middot; RBAC &middot; Rate Limiting</div>
                </div>
                <div className="flex justify-center">
                  <ChevronRight className="h-5 w-5 rotate-90 text-muted-foreground/40" />
                </div>
                {/* Data layer */}
                <div className="grid grid-cols-2 gap-3">
                  <div className="rounded-xl border border-border/60 bg-muted/30 p-4 text-center">
                    <div className="text-xs font-medium text-foreground uppercase tracking-wider mb-1">PostgreSQL</div>
                    <div className="text-xs text-muted-foreground">Users, Projects,<br />Audit Logs</div>
                  </div>
                  <div className="rounded-xl border border-green-500/20 bg-green-500/5 p-4 text-center">
                    <div className="text-xs font-medium text-green-600 dark:text-green-400 uppercase tracking-wider mb-1">Vault KV-v2</div>
                    <div className="text-xs text-muted-foreground">Secret Values<br />(Encrypted)</div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* ─── CLI ─── */}
      <section id="cli" className="border-y border-border/50 bg-muted/20 py-24 lg:py-32">
        <div className="mx-auto max-w-6xl px-6">
          <div className="mx-auto mb-16 max-w-2xl text-center">
            <h2 className="text-3xl font-bold tracking-tight sm:text-4xl">
              Your secrets, one command away
            </h2>
            <p className="mt-4 text-lg text-muted-foreground">
              A full-featured CLI that fits into any workflow — local dev, CI/CD pipelines, or team onboarding.
            </p>
          </div>

          <CLIPreview />

          {/* CLI command reference */}
          <div className="mx-auto mt-12 max-w-2xl">
            <div className="grid grid-cols-2 gap-3 sm:grid-cols-3">
              {[
                { cmd: "init", desc: "Create project" },
                { cmd: "env pull", desc: "Download secrets" },
                { cmd: "env push", desc: "Upload .env file" },
                { cmd: "secret set", desc: "Set a secret" },
                { cmd: "secret get", desc: "Read a secret" },
                { cmd: "onboard", desc: "Invite teammate" },
                { cmd: "rotate", desc: "Rotate tokens" },
                { cmd: "env list", desc: "List all keys" },
                { cmd: "secret delete", desc: "Remove a secret" },
              ].map(({ cmd, desc }) => (
                <div key={cmd} className="rounded-lg border border-border/50 bg-card/50 px-3 py-2.5">
                  <code className="text-xs font-semibold text-primary font-mono">envault {cmd}</code>
                  <p className="mt-0.5 text-xs text-muted-foreground">{desc}</p>
                </div>
              ))}
            </div>
          </div>
        </div>
      </section>

      {/* ─── CTA ─── */}
      <section className="py-24 lg:py-32">
        <div className="mx-auto max-w-4xl px-6 text-center">
          <div className="relative overflow-hidden rounded-3xl border border-border/50 bg-gradient-to-br from-primary/5 via-card to-primary/5 px-8 py-16 shadow-xl sm:px-16">
            {/* Glow */}
            <div className="pointer-events-none absolute inset-0 bg-gradient-to-r from-primary/5 via-transparent to-primary/5 animate-shimmer" />

            <div className="relative">
              <div className="mx-auto mb-6 flex h-16 w-16 items-center justify-center rounded-2xl bg-primary shadow-lg shadow-primary/25">
                <Shield className="h-8 w-8 text-primary-foreground" />
              </div>
              <h2 className="text-3xl font-bold tracking-tight sm:text-4xl">
                Stop sharing secrets in Slack
              </h2>
              <p className="mx-auto mt-4 max-w-lg text-muted-foreground">
                Envault gives your team a secure, audited, role-scoped workflow for managing environment variables. Set up in minutes, not days.
              </p>
              <div className="mt-8 flex flex-col items-center justify-center gap-3 sm:flex-row">
                <Link href="/login">
                  <Button size="lg" className="gap-2 px-8 shadow-lg shadow-primary/20">
                    Get Started Free <ArrowRight className="h-4 w-4" />
                  </Button>
                </Link>
                <a href="https://github.com/bhartiyaanshul/envault" target="_blank" rel="noopener noreferrer">
                  <Button variant="outline" size="lg" className="gap-2 px-8">
                    <svg className="h-4 w-4" viewBox="0 0 24 24" fill="currentColor">
                      <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
                    </svg>
                    View on GitHub
                  </Button>
                </a>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* ─── FOOTER ─── */}
      <footer className="border-t border-border/50 bg-muted/20">
        <div className="mx-auto max-w-6xl px-6 py-12">
          <div className="grid gap-8 sm:grid-cols-2 lg:grid-cols-4">
            {/* Brand */}
            <div className="sm:col-span-2 lg:col-span-1">
              <div className="flex items-center gap-2.5">
                <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary">
                  <Shield className="h-4 w-4 text-primary-foreground" />
                </div>
                <span className="text-lg font-bold">Envault</span>
              </div>
              <p className="mt-3 text-sm text-muted-foreground leading-relaxed">
                Secure secrets management for teams. Built with Go, HashiCorp Vault, and Next.js.
              </p>
            </div>

            {/* Product */}
            <div>
              <h4 className="mb-3 text-sm font-semibold">Product</h4>
              <ul className="space-y-2 text-sm text-muted-foreground">
                <li><a href="#features" className="transition-colors hover:text-foreground">Features</a></li>
                <li><a href="#security" className="transition-colors hover:text-foreground">Security</a></li>
                <li><a href="#cli" className="transition-colors hover:text-foreground">CLI</a></li>
                <li><a href="#how-it-works" className="transition-colors hover:text-foreground">How it Works</a></li>
              </ul>
            </div>

            {/* Developers */}
            <div>
              <h4 className="mb-3 text-sm font-semibold">Developers</h4>
              <ul className="space-y-2 text-sm text-muted-foreground">
                <li><a href="https://github.com/bhartiyaanshul/envault" target="_blank" rel="noopener noreferrer" className="transition-colors hover:text-foreground">GitHub</a></li>
                <li><Link href="/login" className="transition-colors hover:text-foreground">Dashboard</Link></li>
              </ul>
            </div>

            {/* Tech Stack */}
            <div>
              <h4 className="mb-3 text-sm font-semibold">Built With</h4>
              <ul className="space-y-2 text-sm text-muted-foreground">
                <li>Go + Chi Router</li>
                <li>HashiCorp Vault</li>
                <li>Next.js 14</li>
                <li>Supabase Auth</li>
                <li>PostgreSQL</li>
              </ul>
            </div>
          </div>

          <div className="mt-10 flex flex-col items-center justify-between gap-4 border-t border-border/50 pt-8 sm:flex-row">
            <p className="text-xs text-muted-foreground">
              &copy; {new Date().getFullYear()} Envault. All rights reserved.
            </p>
            <p className="text-xs text-muted-foreground">
              Secrets are encrypted and stored in HashiCorp Vault. Values never touch the metadata database.
            </p>
          </div>
        </div>
      </footer>
    </div>
  );
}
