import { getCostEstimateOptions } from "@/api/@tanstack/react-query.gen";
import type { Project } from "@/api/types.gen";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { formatRupiah } from "@/lib/utils";
import { useQuery } from "@tanstack/react-query";
import { Check } from "lucide-react";
import { PROGRESS_STAGES } from "./constants";

export interface ProjectDetailDialogProps {
  project: Project;
  open: boolean;
  onClose: () => void;
  onNavigate: (to: string) => void;
}

export function ProjectDetailDialog({
  project,
  open,
  onClose,
  onNavigate,
}: ProjectDetailDialogProps) {
  const { data: cost } = useQuery({
    ...getCostEstimateOptions({ path: { id: project.id } }),
    enabled: open,
  });

  return (
    <Dialog open={open} onOpenChange={(val) => !val && onClose()}>
      <DialogContent className="max-w-2xl">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2 text-xl">
            <span>{project.name}</span>
            <Badge variant="outline" className="font-normal">
              {project.section || "General"} · {project.plant || "Site"}
            </Badge>
          </DialogTitle>
          <DialogDescription>
            Ringkasan estimasi biaya dan tahapan progres proyek.
          </DialogDescription>
        </DialogHeader>

        {/* Financial Breakdown (Vertical Stack) */}
        <div className="flex flex-col gap-5 py-3">
          <div className="flex flex-col divide-y divide-border/60 rounded-lg border bg-muted/20 text-sm">
            <div className="flex items-center justify-between px-4 py-2.5">
              <span className="text-muted-foreground">Budget</span>
              <span className="font-semibold text-foreground">
                {formatRupiah(cost?.budget ?? project.budget)}
              </span>
            </div>
            <div className="flex items-center justify-between px-4 py-2.5">
              <span className="text-muted-foreground">RAB (Rencana)</span>
              <span className="font-medium text-foreground">
                {formatRupiah(cost?.rab ?? "0")}
              </span>
            </div>
            <div className="flex items-center justify-between px-4 py-2.5">
              <span className="text-muted-foreground">SPB (Terbayar)</span>
              <span className="font-medium text-foreground">
                {formatRupiah(cost?.spb ?? "0")}
              </span>
            </div>
            <div className="flex items-center justify-between px-4 py-2.5 bg-muted/30">
              <span className="text-muted-foreground">Variance (Budget - RAB)</span>
              <span className="font-semibold text-foreground">
                {formatRupiah(cost?.variance ?? "0")}
              </span>
            </div>
            <div className="flex items-center justify-between px-4 py-2.5 bg-muted/30">
              <span className="text-muted-foreground">Sisa Anggaran (Budget - SPB)</span>
              <span className="font-semibold text-emerald-600 dark:text-emerald-400">
                {formatRupiah(cost?.sisa ?? "0")}
              </span>
            </div>
          </div>

          {/* Project Details */}
          {(project.implementation_plan || project.note) && (
            <div className="flex flex-col gap-1.5 text-xs text-muted-foreground px-1">
              {project.implementation_plan && (
                <div>
                  <span className="font-medium text-foreground">Plan Date:</span>{" "}
                  {project.implementation_plan}
                </div>
              )}
              {project.note && (
                <div>
                  <span className="font-medium text-foreground">Note:</span>{" "}
                  {project.note}
                </div>
              )}
            </div>
          )}

          {/* Read-Only Progress Stages Status */}
          <div className="flex flex-col gap-2">
            <span className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              Progress Stages
            </span>
            <div className="flex flex-wrap items-center gap-2 pt-1">
              {PROGRESS_STAGES.map((s) => {
                const isDone = (project as any)[s.key] === "done";
                return isDone ? (
                  <Badge
                    key={s.key}
                    variant="secondary"
                    className="h-auto flex items-center gap-1.5 px-2.5 py-1 text-xs font-medium bg-emerald-500/15 text-emerald-600 dark:text-emerald-400 border border-emerald-500/30"
                  >
                    <Check className="size-3.5 stroke-[2.5]" />
                    <span>{s.label}</span>
                  </Badge>
                ) : (
                  <Badge
                    key={s.key}
                    variant="outline"
                    className="h-auto flex items-center gap-1.5 px-2.5 py-1 text-xs text-muted-foreground border-dashed bg-muted/20 font-normal"
                  >
                    <span className="size-1.5 rounded-full bg-muted-foreground/40" />
                    <span>{s.label}</span>
                  </Badge>
                );
              })}
            </div>
          </div>
        </div>

        <DialogFooter className="flex-row items-center justify-between sm:justify-between">
          <div className="flex gap-2">
            <Button
              size="sm"
              variant="outline"
              onClick={() => onNavigate(`/schedule?projectId=${project.id}`)}
            >
              Schedule
            </Button>
            <Button
              size="sm"
              variant="outline"
              onClick={() => onNavigate(`/volume?projectId=${project.id}`)}
            >
              Volume Take-off
            </Button>
            <Button
              size="sm"
              variant="default"
              onClick={() => onNavigate(`/rab?projectId=${project.id}`)}
            >
              RAB
            </Button>
          </div>
          <Button variant="ghost" size="sm" onClick={onClose}>
            Close
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

