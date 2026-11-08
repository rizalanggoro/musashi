import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { cn } from "cn";
import { UNIT_PRICE_CATEGORIES } from "./constants";

export interface UnitPriceCategoryFilterProps {
  value: string;
  onChange: (value: string) => void;
  categories?: readonly string[];
  className?: string;
  showLabel?: boolean;
}

export function UnitPriceCategoryFilter({
  value,
  onChange,
  categories = UNIT_PRICE_CATEGORIES,
  className,
  showLabel = true,
}: UnitPriceCategoryFilterProps) {
  return (
    <div className={cn("flex items-center gap-2", className)}>
      {showLabel && (
        <span className="text-sm font-medium text-muted-foreground whitespace-nowrap">
          Category:
        </span>
      )}
      <Select
        items={categories.map((cat) => ({ value: cat, label: cat }))}
        value={value}
        onValueChange={(val) => {
          if (val) {
            onChange(val);
          }
        }}
      >
        <SelectTrigger className="w-44">
          <SelectValue placeholder="Category" />
        </SelectTrigger>
        <SelectContent>
          <SelectGroup>
            {categories.map((cat) => (
              <SelectItem key={cat} value={cat}>
                {cat}
              </SelectItem>
            ))}
          </SelectGroup>
        </SelectContent>
      </Select>
    </div>
  );
}
