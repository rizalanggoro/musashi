import {
  createUnitPriceMutation,
  updateUnitPriceMutation,
} from "@/api/@tanstack/react-query.gen";
import type { UnitPriceDto } from "@/api/types.gen";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
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
import { UNIT_PRICE_FORM_CATEGORIES, type UnitPriceFormCategory } from "./constants";

export interface UnitPriceDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  item?: UnitPriceDto | null;
  onSuccess?: () => void;
}

function UnitPriceForm({
  item,
  onClose,
  onSuccess,
}: {
  item: UnitPriceDto | null;
  onClose: () => void;
  onSuccess?: () => void;
}) {
  const queryClient = useQueryClient();

  const [code, setCode] = useState(item?.code || "");
  const [category, setCategory] = useState<UnitPriceFormCategory>(
    (item?.category as UnitPriceFormCategory) || "MATERIALS"
  );
  const [workItem, setWorkItem] = useState(item?.work_item || "");
  const [unit, setUnit] = useState(item?.unit || "");
  const [unitPrice, setUnitPrice] = useState(item?.unit_price || "");
  const [note, setNote] = useState(item?.note || "");
  const [errorMessage, setErrorMessage] = useState("");

  const createMut = useMutation(createUnitPriceMutation());
  const updateMut = useMutation(updateUnitPriceMutation());

  const isPending = item ? updateMut.isPending : createMut.isPending;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setErrorMessage("");

    if (!workItem.trim() || !unit.trim()) {
      setErrorMessage("Item dan unit wajib diisi.");
      return;
    }

    if (item && !code.trim()) {
      setErrorMessage("Code wajib diisi.");
      return;
    }

    try {
      if (item) {
        await updateMut.mutateAsync({
          path: { id: item.id },
          body: {
            code: code.trim(),
            category,
            work_item: workItem.trim(),
            unit: unit.trim(),
            unit_price: unitPrice.replace(/\D/g, "") || "0",
            note: note.trim(),
          },
        });
      } else {
        await createMut.mutateAsync({
          body: {
            code: code.trim(),
            category,
            work_item: workItem.trim(),
            unit: unit.trim(),
            unit_price: unitPrice.replace(/\D/g, "") || "0",
            note: note.trim(),
          },
        });
      }

      await queryClient.invalidateQueries();
      onSuccess?.();
      onClose();
    } catch (err: any) {
      setErrorMessage(
        err?.message ||
          (item ? "Gagal mengubah harga satuan." : "Gagal menyimpan harga satuan.")
      );
    }
  };

  return (
    <>
      <DialogHeader>
        <DialogTitle>{item ? "Edit Unit Price" : "Add Unit Price"}</DialogTitle>
        <DialogDescription>
          {item
            ? "Perbarui data harga satuan."
            : "Tambahkan harga satuan baru ke master database."}
        </DialogDescription>
      </DialogHeader>

      {errorMessage && (
        <Alert variant="destructive">
          <AlertDescription>{errorMessage}</AlertDescription>
        </Alert>
      )}

      <form onSubmit={handleSubmit} className="flex flex-col gap-4">
        <FieldGroup>
          <div className="grid grid-cols-2 gap-4">
            <Field>
              <FieldLabel htmlFor="up-code">
                Code {!item && <span className="text-xs text-muted-foreground font-normal">(opsional)</span>}
              </FieldLabel>
              <Input
                id="up-code"
                value={code}
                onChange={(e) => setCode(e.target.value)}
                placeholder={item ? "MTR-001" : "Otomatis jika kosong"}
                required={Boolean(item)}
              />
            </Field>
            <Field>
              <FieldLabel htmlFor="up-category">Category</FieldLabel>
              <Select
                value={category}
                onValueChange={(val) => setCategory(val as UnitPriceFormCategory)}
              >
                <SelectTrigger id="up-category">
                  <SelectValue placeholder="Pilih Kategori" />
                </SelectTrigger>
                <SelectContent>
                  <SelectGroup>
                    {UNIT_PRICE_FORM_CATEGORIES.map((cat) => (
                      <SelectItem key={cat} value={cat}>
                        {cat}
                      </SelectItem>
                    ))}
                  </SelectGroup>
                </SelectContent>
              </Select>
            </Field>
          </div>

          <Field>
            <FieldLabel htmlFor="up-workitem">Work Item</FieldLabel>
            <Input
              id="up-workitem"
              value={workItem}
              onChange={(e) => setWorkItem(e.target.value)}
              placeholder="Nama item pekerjaan / material..."
              required
            />
          </Field>

          <div className="grid grid-cols-2 gap-4">
            <Field>
              <FieldLabel htmlFor="up-unit">Unit</FieldLabel>
              <Input
                id="up-unit"
                value={unit}
                onChange={(e) => setUnit(e.target.value)}
                placeholder="m3, kg, OH, dll."
                required
              />
            </Field>
            <Field>
              <FieldLabel htmlFor="up-price">Unit Price (Rp)</FieldLabel>
              <Input
                id="up-price"
                value={unitPrice}
                onChange={(e) => setUnitPrice(e.target.value)}
                placeholder="150000"
                required
              />
            </Field>
          </div>

          <Field>
            <FieldLabel htmlFor="up-note">Note</FieldLabel>
            <Input
              id="up-note"
              value={note}
              onChange={(e) => setNote(e.target.value)}
              placeholder="Catatan spesifikasi..."
            />
          </Field>
        </FieldGroup>

        <DialogFooter>
          <Button type="button" variant="outline" onClick={onClose}>
            Cancel
          </Button>
          <Button type="submit" disabled={isPending}>
            {isPending ? "Saving..." : "Save"}
          </Button>
        </DialogFooter>
      </form>
    </>
  );
}

export function UnitPriceDialog({
  open,
  onOpenChange,
  item,
  onSuccess,
}: UnitPriceDialogProps) {
  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-md">
        {open && (
          <UnitPriceForm
            item={item ?? null}
            onClose={() => onOpenChange(false)}
            onSuccess={onSuccess}
          />
        )}
      </DialogContent>
    </Dialog>
  );
}

export function AddUnitPriceDialog({
  open,
  onOpenChange,
  onSuccess,
}: Omit<UnitPriceDialogProps, "item">) {
  return (
    <UnitPriceDialog
      open={open}
      onOpenChange={onOpenChange}
      item={null}
      onSuccess={onSuccess}
    />
  );
}

export function EditUnitPriceDialog({
  open,
  onOpenChange,
  item,
  onSuccess,
}: UnitPriceDialogProps) {
  return (
    <UnitPriceDialog
      open={open}
      onOpenChange={onOpenChange}
      item={item}
      onSuccess={onSuccess}
    />
  );
}

