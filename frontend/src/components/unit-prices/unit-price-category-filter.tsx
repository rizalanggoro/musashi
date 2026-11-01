import { ToggleGroup, ToggleGroupItem } from "@/components/ui/toggle-group";
import { UNIT_PRICE_CATEGORIES } from "./constants";

export interface UnitPriceCategoryFilterProps {
  value: string;
  onChange: (value: string) => void;
  categories?: readonly string[];
  className?: string;
}

export function UnitPriceCategoryFilter({
  value,
  onChange,
  categories = UNIT_PRICE_CATEGORIES,
  className,
}: UnitPriceCategoryFilterProps) {
  return (
    <ToggleGroup
      value={[value]}
      onValueChange={(val) => {
        if (val && val.length > 0) {
          onChange(val[0]);
        }
      }}
      spacing={2}
      className={className || "flex flex-wrap"}
    >
      {categories.map((cat) => (
        <ToggleGroupItem key={cat} value={cat} className="h-8 px-3 text-xs">
          {cat}
        </ToggleGroupItem>
      ))}
    </ToggleGroup>
  );
}

