import {
  deleteUnitPriceMutation,
  listUnitPricesOptions,
} from "@/api/@tanstack/react-query.gen";
import type { UnitPriceDto } from "@/api/types.gen";
import { PageHeader } from "@/components/page-header";
import {
  UnitPriceCategoryFilter,
  UnitPriceDialog,
  UnitPriceTable,
} from "@/components/unit-prices";
import { useAuth } from "@/lib/auth";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState } from "react";

export const Route = createFileRoute("/unit-prices/")({
  component: UnitPricesPage,
  validateSearch: (search: Record<string, unknown>): { q?: string; category?: string } => ({
    q: (search.q as string) || undefined,
    category: (search.category as string) || undefined,
  }),
});

function UnitPricesPage() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { q, category } = Route.useSearch();
  const { isAdmin, user, openLoginModal } = useAuth();

  const currentCategory = category || "ALL";
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const [editingItem, setEditingItem] = useState<UnitPriceDto | null>(null);

  const { data: items = [], isLoading } = useQuery({
    ...listUnitPricesOptions({
      query: {
        q: q || undefined,
        category: currentCategory === "ALL" ? undefined : (currentCategory as any),
      },
    }),
  });

  const deleteMut = useMutation(deleteUnitPriceMutation());

  const handleOpenAdd = () => {
    if (!user) {
      openLoginModal();
      return;
    }
    setEditingItem(null);
    setIsDialogOpen(true);
  };

  const handleOpenEdit = (item: UnitPriceDto) => {
    if (!user) {
      openLoginModal();
      return;
    }
    setEditingItem(item);
    setIsDialogOpen(true);
  };

  const handleDelete = async (item: UnitPriceDto) => {
    if (!confirm(`Hapus harga satuan ${item.code}?`)) return;
    try {
      await deleteMut.mutateAsync({ path: { id: item.id } });
      await queryClient.invalidateQueries();
    } catch (err: any) {
      alert(
        err?.message ||
          "Gagal menghapus: Item mungkin sedang digunakan pada AHSP."
      );
    }
  };

  return (
    <div className="flex flex-col gap-4 min-w-0">
      <PageHeader
        title="Unit Prices"
        searchValue={q || ""}
        onSearch={(val) => {
          navigate({
            search: ((prev: any) => ({ ...prev, q: val || undefined })) as any,
            replace: true,
          });
        }}
        onAdd={handleOpenAdd}
        addLabel="Unit Price"
      />

      <UnitPriceCategoryFilter
        value={currentCategory}
        onChange={(val) => {
          navigate({
            search: ((prev: any) => ({
              ...prev,
              category: val === "ALL" ? undefined : val,
            })) as any,
            replace: true,
          });
        }}
      />

      <UnitPriceTable
        items={items}
        isLoading={isLoading}
        isAdmin={isAdmin}
        search={q || ""}
        onEdit={handleOpenEdit}
        onDelete={handleDelete}
      />

      <UnitPriceDialog
        open={isDialogOpen}
        onOpenChange={(open) => {
          setIsDialogOpen(open);
          if (!open) setEditingItem(null);
        }}
        item={editingItem}
      />
    </div>
  );
}
