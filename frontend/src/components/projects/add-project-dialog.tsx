import { createProjectMutation } from "@/api/@tanstack/react-query.gen";
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

export interface AddProjectDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onSuccess?: () => void;
}

function AddProjectForm({
  onClose,
  onSuccess,
}: {
  onClose: () => void;
  onSuccess?: () => void;
}) {
  const queryClient = useQueryClient();
  const anchor = useComboboxAnchor();

  const [name, setName] = useState("");
  const [section, setSection] = useState<ProjectSection>("GA Building");
  const [plant, setPlant] = useState<ProjectPlant>("Cikarang");
  const [planDate, setPlanDate] = useState("");
  const [budget, setBudget] = useState("");
  const [note, setNote] = useState("");
  const [stages, setStages] = useState<StageOption[]>([]);
  const [formError, setFormError] = useState("");

  const createMut = useMutation(createProjectMutation());

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

      await createMut.mutateAsync({
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
      setFormError(err?.message || "Gagal membuat proyek.");
    }
  };

  return (
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

      <form onSubmit={handleSubmit} className="flex flex-col gap-4">
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
              <Select
                value={section}
                onValueChange={(val) => {
                  if (val && PROJECT_SECTIONS.includes(val as ProjectSection)) {
                    setSection(val as ProjectSection);
                  }
                }}
              >
                <SelectTrigger id="proj-section" className="w-full">
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
              <FieldLabel htmlFor="proj-plant">Plant</FieldLabel>
              <Select
                value={plant}
                onValueChange={(val) => {
                  if (val && PROJECT_PLANTS.includes(val as ProjectPlant)) {
                    setPlant(val as ProjectPlant);
                  }
                }}
              >
                <SelectTrigger id="proj-plant" className="w-full">
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
            <FieldLabel htmlFor="proj-note">Note</FieldLabel>
            <Input
              id="proj-note"
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
          <Button type="submit" disabled={createMut.isPending}>
            {createMut.isPending ? "Saving..." : "Save Project"}
          </Button>
        </DialogFooter>
      </form>
    </DialogContent>
  );
}

export function AddProjectDialog({
  open,
  onOpenChange,
  onSuccess,
}: AddProjectDialogProps) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      {open && (
        <AddProjectForm
          onClose={() => onOpenChange(false)}
          onSuccess={onSuccess}
        />
      )}
    </Dialog>
  );
}
