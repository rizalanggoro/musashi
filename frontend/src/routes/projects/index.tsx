import {
  deleteProjectMutation,
  listProjectsOptions,
} from "@/api/@tanstack/react-query.gen";
import type { Project } from "@/api/types.gen";
import { PageHeader } from "@/components/page-header";
import {
  AddProjectDialog,
  EditProjectDialog,
  PROGRESS_STAGES,
  ProjectDetailDialog,
} from "@/components/projects";
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
import { useAuth } from "@/lib/auth";
import { formatRupiah } from "@/lib/utils";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { Check, MoreHorizontal, Pencil, Trash2 } from "lucide-react";
import { useState } from "react";

export const Route = createFileRoute("/projects/")({ component: ProjectsPage });

function ProjectsPage() {
  const queryClient = useQueryClient();
  const navigate = useNavigate();
  const { isAdmin, user, openLoginModal } = useAuth();

  const [search, setSearch] = useState("");
  const [page] = useState(1);
  const [isAddOpen, setIsAddOpen] = useState(false);
  const [editingProject, setEditingProject] = useState<Project | null>(null);
  const [selectedProject, setSelectedProject] = useState<Project | null>(null);

  const { data, isLoading } = useQuery({
    ...listProjectsOptions({
      query: { q: search || undefined, page, per_page: 25 },
    }),
  });

  const deleteMut = useMutation(deleteProjectMutation());

  const handleOpenAdd = () => {
    if (!user) {
      openLoginModal();
      return;
    }
    setIsAddOpen(true);
  };

  const handleOpenEdit = (p: Project) => {
    if (!user) {
      openLoginModal();
      return;
    }
    setEditingProject(p);
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
              <TableHead>Note</TableHead>
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
                <TableCell
                  colSpan={colSpanCount}
                  className="h-48 text-center"
                >
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
                  (s) => (p as any)[s.key] === "done"
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
                      {p.note || "-"}
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

      <AddProjectDialog
        open={isAddOpen}
        onOpenChange={setIsAddOpen}
      />

      <EditProjectDialog
        project={editingProject}
        open={Boolean(editingProject)}
        onOpenChange={(open) => {
          if (!open) setEditingProject(null);
        }}
      />

      {selectedProject && (
        <ProjectDetailDialog
          project={selectedProject}
          open={Boolean(selectedProject)}
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
