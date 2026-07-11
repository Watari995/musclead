"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { useTranslations } from "next-intl";
import type { TrainingDTO } from "@/shared/api/client";
import { useAccessToken } from "@/shared/auth/access-token";
import { useExercisesQuery } from "@/features/training/api/exercises";
import {
  useRoutinesQuery,
  useStartTrainingFromRoutineMutation,
} from "@/features/training/api/routines";
import {
  useDeleteTrainingMutation,
  useTrainingsQuery,
} from "@/features/training/api/trainings";
import { formatDateTime } from "@/features/training/model/training-draft";
import { Button, Card, ErrorText, Popover, SectionTitle } from "@/shared/ui";

export default function TrainingsPage() {
  const router = useRouter();
  const { token, ready } = useAccessToken();

  useEffect(() => {
    if (ready && !token) router.replace("/login");
  }, [ready, token, router]);

  const query = useTrainingsQuery(Boolean(token));
  const exercisesQuery = useExercisesQuery(Boolean(token));
  const routinesQuery = useRoutinesQuery(Boolean(token));
  const startFromRoutine = useStartTrainingFromRoutineMutation();
  const exerciseNameByID = new Map<string, string>();
  for (const ex of exercisesQuery.data ?? []) {
    exerciseNameByID.set(ex.id, ex.name);
  }

  const t = useTranslations("trainings");
  const tr = useTranslations("routines");
  const tc = useTranslations("common");
  const del = useDeleteTrainingMutation();

  if (!ready || !token) return null;

  const routines = routinesQuery.data ?? [];
  const hasRoutines = routines.length > 0;

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <SectionTitle>{t("title")}</SectionTitle>
        <div className="flex items-center gap-2 flex-wrap">
          <Popover
            align="end"
            trigger={({
              onClick,
              "aria-expanded": ariaExpanded,
              "aria-controls": ariaControls,
            }) => (
              <Button
                type="button"
                variant="ghost"
                onClick={onClick}
                disabled={!hasRoutines || startFromRoutine.isPending}
                title={
                  hasRoutines ? undefined : t("pleaseRegisterRoutines")
                }
                aria-expanded={ariaExpanded}
                aria-controls={ariaControls}
              >
                {startFromRoutine.isPending
                  ? t("starting")
                  : t("startFromRoutine")}
              </Button>
            )}
          >
            <ul className="py-1 min-w-[14rem] max-h-72 overflow-auto">
              {routines.map((r) => (
                <li key={r.id}>
                  <button
                    type="button"
                    onClick={() =>
                      startFromRoutine.mutate(r, {
                        onSuccess: (data) => {
                          const id = data.training_id ?? "";
                          if (id) router.push(`/trainings/${id}/edit`);
                        },
                      })
                    }
                    className="w-full text-left px-4 py-2 text-sm hover:bg-[var(--color-surface-alt)]"
                  >
                    <div className="font-medium truncate">{r.name}</div>
                    <div className="text-xs text-[var(--color-ink-muted)]">
                      {tr("exerciseCount", { count: (r.routine_exercises ?? []).length })}
                    </div>
                  </button>
                </li>
              ))}
            </ul>
          </Popover>
          <Link href="/trainings/new">
            <Button>{t("newRecord")}</Button>
          </Link>
        </div>
      </div>

      {startFromRoutine.isError && (
        <ErrorText>{(startFromRoutine.error as Error).message}</ErrorText>
      )}

      {query.isLoading && (
        <p className="text-sm text-[var(--color-ink-muted)]">{tc("loading")}</p>
      )}
      {query.isError && (
        <ErrorText>{(query.error as Error).message}</ErrorText>
      )}

      {query.data && query.data.length === 0 && (
        <Card className="p-8 text-center text-sm text-[var(--color-ink-muted)]">
          {t("noTrainings")}
        </Card>
      )}

      {query.data && query.data.length > 0 && (
        <ul className="space-y-3">
          {query.data.map((t) => (
            <TrainingCard
              key={t.id}
              training={t}
              exerciseNameByID={exerciseNameByID}
              onDelete={() => {
                if (confirm(t("deleteConfirm"))) {
                  del.mutate(t.id ?? "");
                }
              }}
              deleting={del.isPending}
            />
          ))}
        </ul>
      )}
    </div>
  );
}

function TrainingCard({
  training,
  exerciseNameByID,
  onDelete,
  deleting,
}: {
  training: TrainingDTO;
  exerciseNameByID: Map<string, string>;
  onDelete: () => void;
  deleting: boolean;
}) {
  const tc = useTranslations("common");
  const totalSets = (training.exercises ?? []).reduce(
    (n, ex) => n + (ex.sets?.length ?? 0),
    0,
  );
  const exerciseNames = (training.exercises ?? [])
    .map((ex) => (ex.exercise_id ? exerciseNameByID.get(ex.exercise_id) : undefined))
    .filter((n): n is string => Boolean(n))
    .slice(0, 3);

  return (
    <li className="bg-[var(--color-surface)] border border-[var(--color-line)] rounded-lg p-4 flex items-start justify-between gap-4">
      <Link
        href={`/trainings/${training.id}`}
        className="flex-1 min-w-0 space-y-1 hover:opacity-70 transition-opacity"
      >
        <div className="flex items-center justify-between">
          <span className="text-sm font-bold tracking-tight">
            {formatDateTime(training.started_at)}
          </span>
          <span className="text-xs text-[var(--color-ink-muted)]">
            {(training.exercises ?? []).length} 種目 / {totalSets} セット
          </span>
        </div>
        {exerciseNames.length > 0 && (
          <p className="text-sm text-[var(--color-ink)] line-clamp-1">
            {exerciseNames.join(" / ")}
            {(training.exercises ?? []).length > 3 && " …"}
          </p>
        )}
        {training.memo && (
          <p className="text-xs text-[var(--color-ink-muted)] line-clamp-1">
            {training.memo}
          </p>
        )}
      </Link>
      <div className="flex flex-col gap-1 shrink-0">
        <Link
          href={`/trainings/${training.id}/edit`}
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
