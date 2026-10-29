import { Button } from "@/components/ui/button";
import {
  InputGroup,
  InputGroupAddon,
  InputGroupInput,
} from "@/components/ui/input-group";
import { Plus, Search } from "lucide-react";

type Props = {
  title: string;
  onSearch?: (value: string) => void;
  onAdd?: () => void;
  addLabel?: string;
};

export function PageHeader({
  title,
  onSearch,
  onAdd,
  addLabel = "Add",
}: Props) {
  return (
    <>
      <div className="mb-6 flex flex-wrap items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-foreground">
            {title}
          </h1>
        </div>

        <div className="flex items-center gap-3">
          {onSearch && (
            <InputGroup className="w-64">
              <InputGroupAddon>
                <Search data-icon="inline-start" />
              </InputGroupAddon>
              <InputGroupInput
                placeholder="Search..."
                onChange={(e) => onSearch?.(e.target.value)}
              />
            </InputGroup>
          )}
          {onAdd && (
            <Button onClick={onAdd}>
              <Plus data-icon="inline-start" />
              {addLabel}
            </Button>
          )}
        </div>
      </div>
    </>
  );
}
