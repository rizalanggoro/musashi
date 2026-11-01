import {
  createProjectMutation,
  deleteProjectMutation,
  getCostEstimateOptions,
  listProjectsOptions,
  updateProjectMutation,
} from "@/api/@tanstack/react-query.gen";
import type { Project } from "@/api/types.gen";
import { PageHeader } from "@/components/page-header";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Combobox,
  ComboboxChip,
  ComboboxChips,
  ComboboxChipsInput,
  ComboboxContent,
  ComboboxEmpty,
  ComboboxItem,
  ComboboxList,
  ComboboxValue,
  useComboboxAnchor,
} from "@/components/ui/combobox";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Empty,
  EmptyDescription,
  EmptyHeader,
  EmptyTitle,
} from "@/components/ui/empty";
import { Field, FieldGroup, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { useAuth } from "@/lib/auth";
import { formatRupiah } from "@/lib/utils";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { Check, MoreHorizontal, Pencil, Trash2 } from "lucide-react";
import { useState } from "react";

type ProgressKey =
  | "progress_ppt"
  | "progress_drawing"
  | "progress_rab"
  | "progress_vendor"
  | "progress_spb"
  | "progress_pr"
  | "progress_ringi"
  | "progress_po"
  | "progress_construction"
  | "progress_bast";

interface StageOption {
  key: ProgressKey;
  value: string;
  label: string;
}

const PROGRESS_STAGES: StageOption[] = [
  { key: "progress_ppt", value: "progress_ppt", label: "PPT" },
  { key: "progress_drawing", value: "progress_drawing", label: "Drawing" },
  { key: "progress_rab", value: "progress_rab", label: "RAB" },
  { key: "progress_vendor", value: "progress_vendor", label: "Vendor" },
  { key: "progress_spb", value: "progress_spb", label: "SPB" },
  { key: "progress_pr", value: "progress_pr", label: "PR" },
  { key: "progress_ringi", value: "progress_ringi", label: "RINGI" },
  { key: "progress_po", value: "progress_po", label: "PO" },
  {
    key: "progress_construction",
    value: "progress_construction",
    label: "Construction",
  },
  { key: "progress_bast", value: "progress_bast", label: "BAST" },
];

export const Route = createFileRoute("/projects/")({ component: ProjectsPage });

function ProjectsPage() {
  const queryClient = useQueryClient();
  const navigate = useNavigate();
  const { isAdmin, user, openLoginModal } = useAuth();

  const addAnchor = useComboboxAnchor();
  const editAnchor = useComboboxAnchor();

  const [search, setSearch] = useState("");
  const [page] = useState(1);
  const [isAddOpen, setIsAddOpen] = useState(false);
  const [isEditOpen, setIsEditOpen] = useState(false);
  const [selectedProject, setSelectedProject] = useState<Project | null>(null);
  const [editingProject, setEditingProject] = useState<Project | null>(null);

  // Add Form state
  const [name, setName] = useState("");
  const [section, setSection] = useState("");
  const [plant, setPlant] = useState("");
  const [planDate, setPlanDate] = useState("");
  const [budget, setBudget] = useState("");
  const [information, setInformation] = useState("");
  const [addStages, setAddStages] = useState<StageOption[]>([]);
  const [formError, setFormError] = useState("");

  // Edit Form state
  const [editName, setEditName] = useState("");
  const [editSection, setEditSection] = useState("");
  const [editPlant, setEditPlant] = useState("");
  const [editPlanDate, setEditPlanDate] = useState("");
  const [editBudget, setEditBudget] = useState("");
  const [editInformation, setEditInformation] = useState("");
  const [editStages, setEditStages] = useState<StageOption[]>([]);
  const [editFormError, setEditFormError] = useState("");

  const { data, isLoading } = useQuery({
    ...listProjectsOptions({
      query: { q: search || undefined, page, per_page: 25 },
    }),
  });

  const createMut = useMutation(createProjectMutation());
  const updateMut = useMutation(updateProjectMutation());
  const deleteMut = useMutation(deleteProjectMutation());

  const handleOpenAdd = () => {
    if (!user) {
      openLoginModal();
      return;
    }
    setName("");
    setSection("");
    setPlant("");
    setPlanDate("");
    setBudget("");
    setInformation("");
    setAddStages([]);
    setFormError("");
    setIsAddOpen(true);
  };

  const handleOpenEdit = (p: Project) => {
    if (!user) {
      openLoginModal();
      return;
    }
    setEditingProject(p);
    setEditName(p.name);
    setEditSection(p.section || "");
    setEditPlant(p.plant || "");
    setEditPlanDate(p.implementation_plan || "");
    setEditBudget(p.budget || "");
    setEditInformation(p.information || "");
    setEditStages(PROGRESS_STAGES.filter((s) => (p as any)[s.key] === "done"));
    setEditFormError("");
    setIsEditOpen(true);
  };

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    setFormError("");
    if (!name.trim()) {
      setFormError("Nama project wajib diisi.");
      return;
    }
    try {
      const progressPayload = Object.fromEntries(
        PROGRESS_STAGES.map((s) => [
          s.key,
          addStages.some((sel) => sel.key === s.key) ? "done" : "pending",
        ]),
      );

      await createMut.mutateAsync({
        body: {
          name: name.trim(),
          section: section.trim(),
          plant: plant.trim(),
          implementation_plan: planDate || undefined,
          budget: budget ? budget.replace(/\D/g, "") : "0",
          information: information.trim(),
          ...progressPayload,
        },
      });
      await queryClient.invalidateQueries();
      setIsAddOpen(false);
    } catch (err: any) {
      setFormError(err?.message || "Gagal membuat proyek.");
    }
  };

  const handleUpdate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!editingProject) return;
    setEditFormError("");
    if (!editName.trim()) {
      setEditFormError("Nama project wajib diisi.");
      return;
    }
    try {
      const progressPayload = Object.fromEntries(
        PROGRESS_STAGES.map((s) => [
          s.key,
          editStages.some((sel) => sel.key === s.key) ? "done" : "pending",
        ]),
      );

      await updateMut.mutateAsync({
        path: { id: editingProject.id },
        body: {
          name: editName.trim(),
          section: editSection.trim(),
          plant: editPlant.trim(),
          implementation_plan: editPlanDate || undefined,
          budget: editBudget ? editBudget.replace(/\D/g, "") : "0",
          information: editInformation.trim(),
          ...progressPayload,
        },
      });
      await queryClient.invalidateQueries();
      setIsEditOpen(false);
      setEditingProject(null);
    } catch (err: any) {
      setEditFormError(err?.message || "Gagal mengubah proyek.");
    }
  };

  const handleDelete = async (id: string, e: React.MouseEvent) => {
    e.stopPropagation();
    if (!confirm("Yakin ingin menghapus proyek ini?")) return;
    await deleteMut.mutateAsync({ path: { id } });
    await queryClient.invalidateQueries();
  };

  const projects = data?.data ?? [];
  const colSpanCount = isAdmin ? 9 : 8;

  return (
    <div className="flex flex-col gap-4">
      <PageHeader
        title="Projects"
        onSearch={setSearch}
        onAdd={handleOpenAdd}
        addLabel="Project"
      />

      <div className="rounded-md border bg-card overflow-x-auto">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead className="w-12 text-center">No</TableHead>
              <TableHead>Project</TableHead>
              <TableHead>Section</TableHead>
              <TableHead>Plant</TableHead>
              <TableHead>Plan Date</TableHead>
              <TableHead className="text-right">Budget</TableHead>
              <TableHead className="min-w-[240px]">Progress</TableHead>
              <TableHead>Information</TableHead>
              {isAdmin && (
                <TableHead className="w-16 text-center">Action</TableHead>
              )}
            </TableRow>
          </TableHeader>
          <TableBody>
            {isLoading ? (
              <TableRow>
                <TableCell
                  colSpan={colSpanCount}
                  className="h-32 text-center text-muted-foreground"
                >
                  Loading data...
                </TableCell>
              </TableRow>
            ) : projects.length === 0 ? (
              <TableRow>
                <TableCell colSpan={colSpanCount} className="h-48 text-center">
                  <Empty>
                    <EmptyHeader>
                      <EmptyTitle>Belum ada proyek</EmptyTitle>
                      <EmptyDescription>
                        {search
                          ? `Tidak ada proyek yang cocok dengan "${search}".`
                          : "Mulai dengan menambahkan proyek baru."}
                      </EmptyDescription>
                    </EmptyHeader>
                  </Empty>
                </TableCell>
              </TableRow>
            ) : (
              projects.map((p, idx) => {
                const doneStages = PROGRESS_STAGES.filter(
                  (s) => (p as any)[s.key] === "done",
                );

                return (
                  <TableRow
                    key={p.id}
                    onClick={() => setSelectedProject(p)}
                    className="cursor-pointer"
                  >
                    <TableCell className="text-center text-muted-foreground">
                      {(page - 1) * 25 + idx + 1}
                    </TableCell>
                    <TableCell className="font-medium text-foreground whitespace-nowrap">
                      {p.name}
                    </TableCell>
                    <TableCell>{p.section || "-"}</TableCell>
                    <TableCell>{p.plant || "-"}</TableCell>
                    <TableCell className="whitespace-nowrap">
                      {p.implementation_plan || "-"}
                    </TableCell>
                    <TableCell className="text-right font-medium whitespace-nowrap">
                      {formatRupiah(p.budget)}
                    </TableCell>
                    <TableCell>
                      <div className="flex flex-wrap items-center gap-1.5 py-1">
                        {doneStages.length === 0 ? (
                          <Badge
                            variant="outline"
                            className="text-[11px] text-muted-foreground font-normal border-dashed"
                          >
                            Pending
                          </Badge>
                        ) : (
                          doneStages.map((s) => (
                            <Badge
                              key={s.key}
                              variant="secondary"
                              className="text-[11px] font-medium bg-emerald-500/15 text-emerald-600 dark:text-emerald-400 border border-emerald-500/20 gap-1"
                            >
                              <Check className="size-3" />
                              {s.label}
                            </Badge>
                          ))
                        )}
                      </div>
                    </TableCell>
                    <TableCell className="max-w-xs truncate text-muted-foreground">
                      {p.information || "-"}
                    </TableCell>
                    {isAdmin && (
                      <TableCell
                        className="text-center"
                        onClick={(e) => e.stopPropagation()}
                      >
                        <DropdownMenu>
                          <DropdownMenuTrigger
                            render={
                              <Button
                                variant="ghost"
                                size="icon"
                                className="size-8 text-muted-foreground hover:text-foreground"
                              />
                            }
                          >
                            <MoreHorizontal className="size-4" />
                            <span className="sr-only">Menu</span>
                          </DropdownMenuTrigger>
                          <DropdownMenuContent align="end" className="w-32">
                            <DropdownMenuItem onClick={() => handleOpenEdit(p)}>
                              <Pencil className="size-4 mr-2" />
                              <span>Ubah</span>
                            </DropdownMenuItem>
                            <DropdownMenuItem
                              variant="destructive"
                              onClick={(e) => handleDelete(p.id, e)}
                            >
                              <Trash2 className="size-4 mr-2" />
                              <span>Hapus</span>
                            </DropdownMenuItem>
                          </DropdownMenuContent>
                        </DropdownMenu>
                      </TableCell>
                    )}
                  </TableRow>
                );
              })
            )}
          </TableBody>
        </Table>
      </div>

      {/* Modal Add Project */}
      <Dialog open={isAddOpen} onOpenChange={setIsAddOpen}>
        <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>Add New Project</DialogTitle>
            <DialogDescription>
              Masukkan rincian proyek baru dan tahapan monitoring.
            </DialogDescription>
          </DialogHeader>

          {formError && (
            <Alert variant="destructive">
              <AlertDescription>{formError}</AlertDescription>
            </Alert>
          )}

          <form onSubmit={handleCreate} className="flex flex-col gap-4">
            <FieldGroup>
              <Field>
                <FieldLabel htmlFor="proj-name">Project Name</FieldLabel>
                <Input
                  id="proj-name"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  placeholder="Nama Proyek..."
                  required
                />
              </Field>

              <div className="grid grid-cols-2 gap-4">
                <Field>
                  <FieldLabel htmlFor="proj-section">Section</FieldLabel>
                  <Input
                    id="proj-section"
                    value={section}
                    onChange={(e) => setSection(e.target.value)}
                    placeholder="Unit / Bagian..."
                  />
                </Field>
                <Field>
                  <FieldLabel htmlFor="proj-plant">Plant</FieldLabel>
                  <Input
                    id="proj-plant"
                    value={plant}
                    onChange={(e) => setPlant(e.target.value)}
                    placeholder="Lokasi Plant / Pabrik..."
                  />
                </Field>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <Field>
                  <FieldLabel htmlFor="proj-plan">Plan Date</FieldLabel>
                  <Input
                    id="proj-plan"
                    type="date"
                    value={planDate}
                    onChange={(e) => setPlanDate(e.target.value)}
                  />
                </Field>
                <Field>
                  <FieldLabel htmlFor="proj-budget">Budget (Rp)</FieldLabel>
                  <Input
                    id="proj-budget"
                    value={budget}
                    onChange={(e) => setBudget(e.target.value)}
                    placeholder="Contoh: 150000000"
                  />
                </Field>
              </div>

              <Field>
                <FieldLabel htmlFor="proj-info">Information</FieldLabel>
                <Input
                  id="proj-info"
                  value={information}
                  onChange={(e) => setInformation(e.target.value)}
                  placeholder="Keterangan tambahan..."
                />
              </Field>

              <Field>
                <FieldLabel>Progress Stages</FieldLabel>
                <Combobox
                  multiple
                  items={PROGRESS_STAGES}
                  value={addStages}
                  onValueChange={(val: StageOption[]) => setAddStages(val)}
                  isItemEqualToValue={(a, b) => a.key === b.key}
                >
                  <ComboboxValue>
                    {(values: StageOption[]) => (
                      <ComboboxChips ref={addAnchor}>
                        {values?.map((stage) => (
                          <ComboboxChip key={stage.key}>
                            {stage.label}
                          </ComboboxChip>
                        ))}
                        <ComboboxChipsInput
                          placeholder={
                            !values || values.length === 0
                              ? "Pilih tahapan progres..."
                              : ""
                          }
                        />
                      </ComboboxChips>
                    )}
                  </ComboboxValue>
                  <ComboboxContent anchor={addAnchor}>
                    <ComboboxEmpty>Tahapan tidak ditemukan</ComboboxEmpty>
                    <ComboboxList>
                      {(stage: StageOption) => (
                        <ComboboxItem key={stage.key} value={stage}>
                          {stage.label}
                        </ComboboxItem>
                      )}
                    </ComboboxList>
                  </ComboboxContent>
                </Combobox>
              </Field>
            </FieldGroup>

            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => setIsAddOpen(false)}
              >
                Cancel
              </Button>
              <Button type="submit" disabled={createMut.isPending}>
                {createMut.isPending ? "Saving..." : "Save Project"}
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>

      {/* Modal Edit Project */}
      <Dialog open={isEditOpen} onOpenChange={setIsEditOpen}>
        <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
          <DialogHeader>
            <DialogTitle>Edit Project</DialogTitle>
            <DialogDescription>
              Ubah data proyek dan perbarui tahapan progres.
            </DialogDescription>
          </DialogHeader>

          {editFormError && (
            <Alert variant="destructive">
              <AlertDescription>{editFormError}</AlertDescription>
            </Alert>
          )}

          <form onSubmit={handleUpdate} className="flex flex-col gap-4">
            <FieldGroup>
              <Field>
                <FieldLabel htmlFor="edit-proj-name">Project Name</FieldLabel>
                <Input
                  id="edit-proj-name"
                  value={editName}
                  onChange={(e) => setEditName(e.target.value)}
                  placeholder="Nama Proyek..."
                  required
                />
              </Field>

              <div className="grid grid-cols-2 gap-4">
                <Field>
                  <FieldLabel htmlFor="edit-proj-section">Section</FieldLabel>
                  <Input
                    id="edit-proj-section"
                    value={editSection}
                    onChange={(e) => setEditSection(e.target.value)}
                    placeholder="Unit / Bagian..."
                  />
                </Field>
                <Field>
                  <FieldLabel htmlFor="edit-proj-plant">Plant</FieldLabel>
                  <Input
                    id="edit-proj-plant"
                    value={editPlant}
                    onChange={(e) => setEditPlant(e.target.value)}
                    placeholder="Lokasi Plant / Pabrik..."
                  />
                </Field>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <Field>
                  <FieldLabel htmlFor="edit-proj-plan">Plan Date</FieldLabel>
                  <Input
                    id="edit-proj-plan"
                    type="date"
                    value={editPlanDate}
                    onChange={(e) => setEditPlanDate(e.target.value)}
                  />
                </Field>
                <Field>
                  <FieldLabel htmlFor="edit-proj-budget">
                    Budget (Rp)
                  </FieldLabel>
                  <Input
                    id="edit-proj-budget"
                    value={editBudget}
                    onChange={(e) => setEditBudget(e.target.value)}
                    placeholder="Contoh: 150000000"
                  />
                </Field>
              </div>

              <Field>
                <FieldLabel htmlFor="edit-proj-info">Information</FieldLabel>
                <Input
                  id="edit-proj-info"
                  value={editInformation}
                  onChange={(e) => setEditInformation(e.target.value)}
                  placeholder="Keterangan tambahan..."
                />
              </Field>

              <Field>
                <FieldLabel>Progress Stages</FieldLabel>
                <Combobox
                  multiple
                  items={PROGRESS_STAGES}
                  value={editStages}
                  onValueChange={(val: StageOption[]) => setEditStages(val)}
                  isItemEqualToValue={(a, b) => a.key === b.key}
                >
                  <ComboboxValue>
                    {(values: StageOption[]) => (
                      <ComboboxChips ref={editAnchor}>
                        {values?.map((stage) => (
                          <ComboboxChip key={stage.key}>
                            {stage.label}
                          </ComboboxChip>
                        ))}
                        <ComboboxChipsInput
                          placeholder={
                            !values || values.length === 0
                              ? "Pilih tahapan progres..."
                              : ""
                          }
                        />
                      </ComboboxChips>
                    )}
                  </ComboboxValue>
                  <ComboboxContent anchor={editAnchor}>
                    <ComboboxEmpty>Tahapan tidak ditemukan</ComboboxEmpty>
                    <ComboboxList>
                      {(stage: StageOption) => (
                        <ComboboxItem key={stage.key} value={stage}>
                          {stage.label}
                        </ComboboxItem>
                      )}
                    </ComboboxList>
                  </ComboboxContent>
                </Combobox>
              </Field>
            </FieldGroup>

            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={() => {
                  setIsEditOpen(false);
                  setEditingProject(null);
                }}
              >
                Cancel
              </Button>
              <Button type="submit" disabled={updateMut.isPending}>
                {updateMut.isPending ? "Saving..." : "Save Changes"}
              </Button>
            </DialogFooter>
          </form>
        </DialogContent>
      </Dialog>

      {/* Project Detail Modal (Read-Only Detail) */}
      {selectedProject && (
        <ProjectDetailDialog
          project={selectedProject}
          open={!!selectedProject}
          onClose={() => setSelectedProject(null)}
          onNavigate={(path) => {
            setSelectedProject(null);
            navigate({ to: path });
          }}
        />
      )}
    </div>
  );
}

function ProjectDetailDialog({
  project,
  open,
  onClose,
  onNavigate,
}: {
  project: Project;
  open: boolean;
  onClose: () => void;
  onNavigate: (to: string) => void;
}) {
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
              <span className="text-muted-foreground">
                Variance (Budget - RAB)
              </span>
              <span className="font-semibold text-foreground">
                {formatRupiah(cost?.variance ?? "0")}
              </span>
            </div>
            <div className="flex items-center justify-between px-4 py-2.5 bg-muted/30">
              <span className="text-muted-foreground">
                Sisa Anggaran (Budget - SPB)
              </span>
              <span className="font-semibold text-emerald-600 dark:text-emerald-400">
                {formatRupiah(cost?.sisa ?? "0")}
              </span>
            </div>
          </div>

          {/* Project Details */}
          {(project.implementation_plan || project.information) && (
            <div className="flex flex-col gap-1.5 text-xs text-muted-foreground px-1">
              {project.implementation_plan && (
                <div>
                  <span className="font-medium text-foreground">
                    Plan Date:
                  </span>{" "}
                  {project.implementation_plan}
                </div>
              )}
              {project.information && (
                <div>
                  <span className="font-medium text-foreground">
                    Keterangan:
                  </span>{" "}
                  {project.information}
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
