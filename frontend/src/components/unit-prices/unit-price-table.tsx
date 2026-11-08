import type { UnitPriceDto } from "@/api/types.gen";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
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
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { formatRupiah } from "@/lib/utils";
import { MoreHorizontal, Pencil, Trash2 } from "lucide-react";

export interface UnitPriceTableProps {
  items: UnitPriceDto[];
  isLoading?: boolean;
  isAdmin?: boolean;
  search?: string;
  onEdit?: (item: UnitPriceDto) => void;
  onDelete?: (item: UnitPriceDto) => void;
}

export function UnitPriceTable({
  items,
  isLoading = false,
  isAdmin = false,
  search = "",
  onEdit,
  onDelete,
}: UnitPriceTableProps) {
  const colSpan = isAdmin ? 8 : 7;

  return (
    <div className="rounded-md border bg-card overflow-hidden">
      <Table className="min-w-[750px]">
        <TableHeader>
          <TableRow>
            <TableHead className="w-12 text-center">No</TableHead>
            <TableHead className="w-28">Code</TableHead>
            <TableHead className="w-28">Category</TableHead>
            <TableHead>Work Item</TableHead>
            <TableHead className="w-20 text-center">Unit</TableHead>
            <TableHead className="w-36 text-right">Unit Price</TableHead>
            <TableHead>Note</TableHead>
            {isAdmin && <TableHead className="w-20 text-center">Action</TableHead>}
          </TableRow>
        </TableHeader>
        <TableBody>
          {isLoading ? (
            <TableRow>
              <TableCell colSpan={colSpan} className="h-32 text-center text-muted-foreground">
                Loading data...
              </TableCell>
            </TableRow>
          ) : items.length === 0 ? (
            <TableRow>
              <TableCell colSpan={colSpan} className="h-48 text-center">
                <Empty>
                  <EmptyHeader>
                    <EmptyTitle>Tidak ada harga satuan</EmptyTitle>
                    <EmptyDescription>
                      {search
                        ? `Tidak ada item yang cocok dengan "${search}".`
                        : "Mulai tambahkan daftar harga satuan baru."}
                    </EmptyDescription>
                  </EmptyHeader>
                </Empty>
              </TableCell>
            </TableRow>
          ) : (
            items.map((item, idx) => (
              <TableRow key={item.id}>
                <TableCell className="text-center text-muted-foreground">{idx + 1}</TableCell>
                <TableCell className="font-mono text-xs font-semibold">{item.code}</TableCell>
                <TableCell>
                  <Badge variant="outline" className="text-[10px] font-normal uppercase">
                    {item.category}
                  </Badge>
                </TableCell>
                <TableCell className="font-medium text-foreground">{item.work_item}</TableCell>
                <TableCell className="text-center text-muted-foreground">{item.unit}</TableCell>
                <TableCell className="text-right font-medium">{formatRupiah(item.unit_price)}</TableCell>
                <TableCell className="text-muted-foreground">{item.note || "-"}</TableCell>
                {isAdmin && (
                  <TableCell className="text-center">
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
                        {onEdit && (
                          <DropdownMenuItem onClick={() => onEdit(item)}>
                            <Pencil className="size-4 mr-2" />
                            <span>Ubah</span>
                          </DropdownMenuItem>
                        )}
                        {onDelete && (
                          <DropdownMenuItem
                            variant="destructive"
                            onClick={() => onDelete(item)}
                          >
                            <Trash2 className="size-4 mr-2" />
                            <span>Hapus</span>
                          </DropdownMenuItem>
                        )}
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </TableCell>
                )}
              </TableRow>
            ))
          )}
        </TableBody>
      </Table>
    </div>
  );
}

