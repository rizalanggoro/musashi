import { updateProjectMutation } from "@/api/@tanstack/react-query.gen";
import type { Project } from "@/api/types.gen";
import { Alert, AlertDescription } from "@/components/ui/alert";
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
import { Field, FieldGroup, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import {
  PROJECT_PLANTS,
  PROJECT_SECTIONS,
  PROGRESS_STAGES,
  type ProjectPlant,
  type ProjectSection,
  type StageOption,
} from "./constants";

export interface EditProjectDialogProps {
  project: Project | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess?: () => void;
}

function EditProjectForm({
  project,
  onClose,
  onSuccess,
}: {
  project: Project;
  onClose: () => void;
  onSuccess?: () => void;
}) {
  const queryClient = useQueryClient();
  const anchor = useComboboxAnchor();

  const [name, setName] = useState(project.name || "");
  const [section, setSection] = useState<ProjectSection>(() => {
    if (PROJECT_SECTIONS.includes(project.section as ProjectSection)) {
      return project.section as ProjectSection;
    }
    return "GA Building";
  });
  const [plant, setPlant] = useState<ProjectPlant>(() => {
    if (PROJECT_PLANTS.includes(project.plant as ProjectPlant)) {
      return project.plant as ProjectPlant;
    }
    return "Cikarang";
  });
  const [planDate, setPlanDate] = useState(project.implementation_plan || "");
  const [budget, setBudget] = useState(project.budget || "");
  const [note, setNote] = useState(project.note || "");
  const [stages, setStages] = useState<StageOption[]>(() =>
    PROGRESS_STAGES.filter((s) => (project as any)[s.key] === "done")
  );
  const [formError, setFormError] = useState("");

  const updateMut = useMutation(updateProjectMutation());

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setFormError("");
    if (!name.trim()) {
      setFormError("Nama project wajib diisi.");
      return;
    }
    if (!PROJECT_SECTIONS.includes(section)) {
      setFormError("Section harus berupa salah satu dari: GA Building, EHS, atau Wellfare.");
      return;
    }
    if (!PROJECT_PLANTS.includes(plant)) {
      setFormError("Plant harus berupa salah satu dari: Cikarang atau Karawang.");
      return;
    }

    try {
      const progressPayload = Object.fromEntries(
        PROGRESS_STAGES.map((s) => [
          s.key,
          stages.some((sel) => sel.key === s.key) ? "done" : "pending",
        ])
      );

      await updateMut.mutateAsync({
        path: { id: project.id },
        body: {
          name: name.trim(),
          section,
          plant,
          implementation_plan: planDate || undefined,
          budget: budget ? budget.replace(/\D/g, "") : "0",
          note: note.trim(),
          ...progressPayload,
        },
      });

      await queryClient.invalidateQueries();
      onClose();
      onSuccess?.();
    } catch (err: any) {
      setFormError(err?.message || "Gagal mengubah proyek.");
    }
  };

  return (
    <DialogContent className="sm:max-w-2xl max-h-[90vh] overflow-y-auto">
      <DialogHeader>
        <DialogTitle>Edit Project</DialogTitle>
        <DialogDescription>
          Ubah data proyek dan perbarui tahapan progres.
        </DialogDescription>
      </DialogHeader>

      {formError && (
        <Alert variant="destructive">
          <AlertDescription>{formError}</AlertDescription>
        </Alert>
      )}

      <form onSubmit={handleSubmit} className="flex flex-col gap-4">
        <FieldGroup>
          <Field>
            <FieldLabel htmlFor="edit-proj-name">Project Name</FieldLabel>
            <Input
              id="edit-proj-name"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Nama Proyek..."
              required
            />
          </Field>

          <div className="grid grid-cols-2 gap-4">
            <Field>
              <FieldLabel htmlFor="edit-proj-section">Section</FieldLabel>
              <Select
                items={PROJECT_SECTIONS.map((sec) => ({ value: sec, label: sec }))}
                value={section}
                onValueChange={(val) => {
                  if (val && PROJECT_SECTIONS.includes(val as ProjectSection)) {
                    setSection(val as ProjectSection);
                  }
                }}
              >
                <SelectTrigger id="edit-proj-section" className="w-full">
                  <SelectValue placeholder="Pilih Section..." />
                </SelectTrigger>
                <SelectContent>
                  <SelectGroup>
                    {PROJECT_SECTIONS.map((sec) => (
                      <SelectItem key={sec} value={sec}>
                        {sec}
                      </SelectItem>
                    ))}
                  </SelectGroup>
                </SelectContent>
              </Select>
            </Field>
            <Field>
              <FieldLabel htmlFor="edit-proj-plant">Plant</FieldLabel>
              <Select
                items={PROJECT_PLANTS.map((pl) => ({ value: pl, label: pl }))}
                value={plant}
                onValueChange={(val) => {
                  if (val && PROJECT_PLANTS.includes(val as ProjectPlant)) {
                    setPlant(val as ProjectPlant);
                  }
                }}
              >
                <SelectTrigger id="edit-proj-plant" className="w-full">
                  <SelectValue placeholder="Pilih Plant..." />
                </SelectTrigger>
                <SelectContent>
                  <SelectGroup>
                    {PROJECT_PLANTS.map((pl) => (
                      <SelectItem key={pl} value={pl}>
                        {pl}
                      </SelectItem>
                    ))}
                  </SelectGroup>
                </SelectContent>
              </Select>
            </Field>
          </div>

          <div className="grid grid-cols-2 gap-4">
            <Field>
              <FieldLabel htmlFor="edit-proj-plan">Plan Date</FieldLabel>
              <Input
                id="edit-proj-plan"
                type="date"
                value={planDate}
                onChange={(e) => setPlanDate(e.target.value)}
              />
            </Field>
            <Field>
              <FieldLabel htmlFor="edit-proj-budget">Budget (Rp)</FieldLabel>
              <Input
                id="edit-proj-budget"
                value={budget}
                onChange={(e) => setBudget(e.target.value)}
                placeholder="Contoh: 150000000"
              />
            </Field>
          </div>

          <Field>
            <FieldLabel htmlFor="edit-proj-note">Note</FieldLabel>
            <Input
              id="edit-proj-note"
              value={note}
              onChange={(e) => setNote(e.target.value)}
              placeholder="Catatan tambahan..."
            />
          </Field>

          <Field>
            <FieldLabel>Progress Stages</FieldLabel>
            <Combobox
              multiple
              items={PROGRESS_STAGES}
              value={stages}
              onValueChange={(val: StageOption[]) => setStages(val)}
              isItemEqualToValue={(a, b) => a.key === b.key}
            >
              <ComboboxValue>
                {(values: StageOption[]) => (
                  <ComboboxChips ref={anchor}>
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
              <ComboboxContent anchor={anchor}>
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
            onClick={onClose}
          >
            Cancel
          </Button>
          <Button type="submit" disabled={updateMut.isPending}>
            {updateMut.isPending ? "Saving..." : "Save Changes"}
          </Button>
        </DialogFooter>
      </form>
    </DialogContent>
  );
}

export function EditProjectDialog({
  project,
  open,
  onOpenChange,
  onSuccess,
}: EditProjectDialogProps) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      {open && project && (
        <EditProjectForm
          key={project.id}
          project={project}
          onClose={() => onOpenChange(false)}
          onSuccess={onSuccess}
        />
      )}
    </Dialog>
  );
}
