"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import {
  DndContext,
  KeyboardSensor,
  PointerSensor,
  closestCenter,
  useSensor,
  useSensors,
  type DragEndEvent,
} from "@dnd-kit/core";
import {
  SortableContext,
  arrayMove,
  sortableKeyboardCoordinates,
  useSortable,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { useTranslations } from "next-intl";
import type { RoutineDTO } from "@/shared/api/client";
import { useAccessToken } from "@/shared/auth/access-token";
import {
  useDeleteRoutineMutation,
  useReorderRoutinesMutation,
  useRoutinesQuery,
} from "@/features/training/api/routines";
import { Button, Card, ErrorText, SectionTitle } from "@/shared/ui";

export default function RoutinesPage() {
  const router = useRouter();
  const { token, ready } = useAccessToken();
  const t = useTranslations("routines");
  const tc = useTranslations("common");

  useEffect(() => {
    if (ready && !token) router.replace("/login");
  }, [ready, token, router]);

  const query = useRoutinesQuery(Boolean(token));
  const del = useDeleteRoutineMutation();
  const reorder = useReorderRoutinesMutation();

  const sensors = useSensors(
    // ハンドルから 6px 以上動かして初めてドラッグ開始。 タップ(編集/削除)と誤検知しないため。
    useSensor(PointerSensor, { activationConstraint: { distance: 6 } }),
    useSensor(KeyboardSensor, { coordinateGetter: sortableKeyboardCoordinates }),
  );

  if (!ready || !token) return null;

  const routines = query.data ?? [];

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event;
    if (!over || active.id === over.id) return;
    const from = routines.findIndex((r) => r.id === active.id);
    const to = routines.findIndex((r) => r.id === over.id);
    if (from === -1 || to === -1) return;
    reorder.mutate(arrayMove(routines, from, to));
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <SectionTitle>{t("title")}</SectionTitle>
        <Link href="/routines/new">
          <Button>{t("newRoutine")}</Button>
        </Link>
      </div>

      {query.isLoading && (
        <p className="text-sm text-[var(--color-ink-muted)]">{tc("loading")}</p>
      )}
      {query.isError && (
        <ErrorText>{(query.error as Error).message}</ErrorText>
      )}
      {del.isError && <ErrorText>{(del.error as Error).message}</ErrorText>}
      {reorder.isError && (
        <ErrorText>{tc("reorderFailed")}</ErrorText>
      )}

      {query.data && routines.length === 0 && (
        <Card className="p-8 text-center text-sm text-[var(--color-ink-muted)]">
          {t("noRoutinesYet")}
        </Card>
      )}

      {routines.length > 0 && (
        <DndContext
          sensors={sensors}
          collisionDetection={closestCenter}
          onDragEnd={handleDragEnd}
        >
          <SortableContext
            items={routines.map((r) => r.id ?? "")}
            strategy={verticalListSortingStrategy}
          >
            <ul className="space-y-3">
              {routines.map((r) => (
                <SortableRoutineCard
                  key={r.id}
                  routine={r}
                  onDelete={() => {
                    if (confirm(t("deleteConfirm", { name: r.name ?? "" }))) {
                      del.mutate(r.id ?? "");
                    }
                  }}
                  deleting={del.isPending}
                />
              ))}
            </ul>
          </SortableContext>
        </DndContext>
      )}
    </div>
  );
}

function SortableRoutineCard({
  routine,
  onDelete,
  deleting,
}: {
  routine: RoutineDTO;
  onDelete: () => void;
  deleting: boolean;
}) {
  const t = useTranslations("routines");
  const tc = useTranslations("common");
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id: routine.id ?? "" });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
  };

  const exerciseNames = (routine.routine_exercises ?? [])
    .map((e) => e.exercise_name)
    .filter(Boolean)
    .slice(0, 3);

  return (
    <li
      ref={setNodeRef}
      style={style}
      className={`bg-[var(--color-surface)] rough p-4 flex items-start justify-between gap-2 ${
        isDragging ? "z-10 shadow-lg opacity-90" : ""
      }`}
    >
      <button
        type="button"
        {...attributes}
        {...listeners}
        aria-label={tc("dragToReorder")}
        className="shrink-0 text-[var(--color-ink-muted)] hover:text-[var(--color-ink)] touch-none cursor-grab active:cursor-grabbing px-1 h-8 inline-flex items-center"
      >
        <GripIcon />
      </button>
      <Link
        href={`/routines/${routine.id}`}
        className="flex-1 min-w-0 space-y-1 hover:opacity-70 transition-opacity"
      >
        <div className="flex items-center justify-between">
          <span className="text-sm font-bold tracking-tight">
            {routine.name}
          </span>
          <span className="text-xs text-[var(--color-ink-muted)]">
            {t("exerciseCount", { count: (routine.routine_exercises ?? []).length })}
          </span>
        </div>
        {exerciseNames.length > 0 && (
          <p className="text-sm text-[var(--color-ink)] line-clamp-1">
            {exerciseNames.join(" / ")}
            {(routine.routine_exercises ?? []).length > 3 && " …"}
          </p>
        )}
      </Link>
      <div className="flex flex-col gap-1 shrink-0">
        <Link
          href={`/routines/${routine.id}/edit`}
          className="text-xs text-[var(--color-ink-muted)] hover:text-[var(--color-ink)] px-2 h-8 inline-flex items-center"
        >
          {tc("edit")}
        </Link>
        <button
          type="button"
          onClick={onDelete}
          disabled={deleting}
          className="text-xs text-[var(--color-ink-muted)] hover:text-[var(--color-accent)] disabled:opacity-50 px-2 h-8"
        >
          {tc("delete")}
        </button>
      </div>
    </li>
  );
}

function GripIcon() {
  return (
    <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor" aria-hidden="true">
      <circle cx="5" cy="3" r="1.4" />
      <circle cx="11" cy="3" r="1.4" />
      <circle cx="5" cy="8" r="1.4" />
      <circle cx="11" cy="8" r="1.4" />
      <circle cx="5" cy="13" r="1.4" />
      <circle cx="11" cy="13" r="1.4" />
    </svg>
  );
}
