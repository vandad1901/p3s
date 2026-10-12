import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";
import type { MouseEvent } from "react";
import { useSearchParams } from "react-router-dom";

type PageItem = number | "ellipsis";

const range = (from: number, to: number) =>
  Array.from({ length: to - from + 1 }, (_, i) => from + i);

function getPageItems(page: number, total: number, siblings = 1): PageItem[] {
  const maxVisible = siblings * 2 + 5;
  if (total <= maxVisible) return range(1, total);

  const left = Math.max(page - siblings, 1);
  const right = Math.min(page + siblings, total);
  const showLeftEllipsis = left > 3;
  const showRightEllipsis = right < total - 2;
  const edgeCount = 3 + 2 * siblings;

  if (!showLeftEllipsis && showRightEllipsis) {
    return [...range(1, edgeCount), "ellipsis", total];
  }
  if (showLeftEllipsis && !showRightEllipsis) {
    return [1, "ellipsis", ...range(total - edgeCount + 1, total)];
  }
  return [1, "ellipsis", ...range(left, right), "ellipsis", total];
}

const disabledLinkProps = {
  "aria-disabled": true,
  tabIndex: -1,
  className: "pointer-events-none opacity-50",
  onClick: (e: MouseEvent) => e.preventDefault(),
} as const;

export function PostsPagination({
  page,
  totalPages,
  label = "Pagination",
}: {
  page: number;
  totalPages: number;
  label?: string;
}) {
  const [searchParams] = useSearchParams();

  if (totalPages <= 1) return null;

  const href = (p: number) => {
    const next = new URLSearchParams(searchParams);
    next.set("page", String(p));
    return `?${next}`;
  };

  return (
    <Pagination
      aria-label={label}
      className="py-2"
    >
      <PaginationContent>
        <PaginationItem>
          <PaginationPrevious
            to={href(Math.max(page - 1, 1))}
            {...(page > 1 ? {} : disabledLinkProps)}
          />
        </PaginationItem>

        {getPageItems(page, totalPages).map((item) => (
          <PaginationItem key={item}>
            {item === "ellipsis" ? (
              <PaginationEllipsis />
            ) : (
              <PaginationLink
                to={href(item)}
                isActive={item === page}
                aria-label={`Page ${item}`}
              >
                {item}
              </PaginationLink>
            )}
          </PaginationItem>
        ))}

        <PaginationItem>
          <PaginationNext
            to={href(Math.min(page + 1, totalPages))}
            {...(page < totalPages ? {} : disabledLinkProps)}
          />
        </PaginationItem>
      </PaginationContent>
    </Pagination>
  );
}
