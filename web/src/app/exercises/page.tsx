"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { useTranslations } from "next-intl";
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
import { useAccessToken } from "@/shared/auth/access-token";
import {
  ExerciseInUseError,
  useDeleteExerciseMutation,
  useExercisesQuery,
  useReorderExercisesMutation,
} from "@/features/training/api/exercises";
import type { Exercise } from "@/features/training/model/exercise";
import { Button, Card, ErrorText, SectionTitle } from "@/shared/ui";

export default function ExercisesPage() {
  const router = useRouter();
  const { token, ready } = useAccessToken();
  const t = useTranslations("exercises");
  const tc = useTranslations("common");

  useEffect(() => {
    if (ready && !token) router.replace("/login");
  }, [ready, token, router]);

  const query = useExercisesQuery(Boolean(token));
  const del = useDeleteExerciseMutation();
  const reorder = useReorderExercisesMutation();

  const sensors = useSensors(
    // ハンドルから 6px 以上動かして初めてドラッグ開始。 タップ(編集/削除)と誤検知しないため。
    useSensor(PointerSensor, { activationConstraint: { distance: 6 } }),
    useSensor(KeyboardSensor, { coordinateGetter: sortableKeyboardCoordinates }),
  );

  if (!ready || !token) return null;

  const exercises = query.data ?? [];

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event;
    if (!over || active.id === over.id) return;
    const from = exercises.findIndex((e) => e.id === active.id);
    const to = exercises.findIndex((e) => e.id === over.id);
    if (from === -1 || to === -1) return;
    reorder.mutate(arrayMove(exercises, from, to));
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <SectionTitle>{t("title")}</SectionTitle>
        <div className="flex items-center gap-2">
          <Link
            href="/exercises/progress"
            className="text-sm text-[var(--color-ink-muted)] hover:text-[var(--color-ink)] px-3 py-2"
          >
            {t("progressGraphLink")}
          </Link>
          <Link href="/exercises/new">
            <Button>{t("newExercise")}</Button>
          </Link>
        </div>
      </div>

      {query.isLoading && (
        <p className="text-sm text-[var(--color-ink-muted)]">{tc("loading")}</p>
      )}
      {query.isError && (
        <ErrorText>{(query.error as Error).message}</ErrorText>
      )}
      {del.isError && (
        <ErrorText>
          {del.error instanceof ExerciseInUseError
            ? del.error.message
            : (del.error as Error).message}
        </ErrorText>
      )}
      {reorder.isError && (
        <ErrorText>{tc("reorderFailed")}</ErrorText>
      )}

      {query.data && exercises.length === 0 && (
        <Card className="p-8 text-center text-sm text-[var(--color-ink-muted)]">
          {t("noExercisesYet")}
        </Card>
      )}

      {exercises.length > 0 && (
        <DndContext
          sensors={sensors}
          collisionDetection={closestCenter}
          onDragEnd={handleDragEnd}
        >
          <SortableContext
            items={exercises.map((e) => e.id)}
            strategy={verticalListSortingStrategy}
          >
            <ul className="space-y-2">
              {exercises.map((ex) => (
                <SortableExerciseRow
                  key={ex.id}
                  exercise={ex}
                  onDelete={() => {
                    if (confirm(t("deleteConfirm", { name: ex.name }))) {
                      del.mutate(ex.id);
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

function SortableExerciseRow({
  exercise,
  onDelete,
  deleting,
}: {
  exercise: Exercise;
  onDelete: () => void;
  deleting: boolean;
}) {
  const tc = useTranslations("common");
  const t = useTranslations("exercises");
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id: exercise.id });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
  };

  return (
    <li
      ref={setNodeRef}
      style={style}
      className={`bg-[var(--color-surface)] rough p-4 flex items-center justify-between gap-2 ${
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
        href={`/exercises/${exercise.id}/edit`}
        className="flex-1 min-w-0 hover:opacity-70 transition-opacity"
      >
        <p className="text-sm font-bold tracking-tight">{exercise.name}</p>
      </Link>
      <div className="flex gap-1 shrink-0">
        <Link
          href={`/exercises/${exercise.id}/edit`}
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
