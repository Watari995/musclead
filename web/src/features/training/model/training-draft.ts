// training フォームの素材。 階層構造(training → exercises → sets)を
// React の useState で扱いやすい immutable な形にまとめる。
// API DTO への双方向変換もここに置く(handler/view から呼ばれる)。

import type {
  RecordTrainingRequest,
  TrainingDTO,
  TrainingExerciseDTO,
  TrainingSetDTO,
} from "@/shared/api/client";

// ─── Draft 型(フォーム保持用) ─────────────────────────

export type SetDraft = {
  /** React key 用のローカル ID(永続化されない) */
  key: string;
  setNumber: number;
  weightKg: string;
  reps: number;
  /** null = 種目のデフォルトを使う */
  restSeconds: number | null;
  memo: string;
};

export type ExerciseDraft = {
  key: string;
  exerciseID: string;
  displayOrder: number;
  restSeconds: number | null;
  memo: string;
  sets: SetDraft[];
};

export type TrainingDraft = {
  /** datetime-local 形式(YYYY-MM-DDTHH:mm) */
  startedAt: string;
  /** datetime-local 形式 or 空文字 */
  endedAt: string;
  memo: string;
  exercises: ExerciseDraft[];
};

// ─── キー生成(React key 用、 サーバー ID と無関係) ───

let keySeed = 0;
function nextKey(prefix: string): string {
  keySeed += 1;
  return `${prefix}-${keySeed}-${Math.random().toString(36).slice(2, 8)}`;
}

// ─── 初期状態 ──────────────────────────────────────

export function createInitialSet(setNumber: number): SetDraft {
  return {
    key: nextKey("set"),
    setNumber,
    weightKg: "",
    reps: 0,
    restSeconds: null,
    memo: "",
  };
}

export function createInitialExercise(displayOrder: number): ExerciseDraft {
  return {
    key: nextKey("ex"),
    exerciseID: "",
    displayOrder,
    restSeconds: null,
    memo: "",
    sets: [createInitialSet(1)],
  };
}

export function createInitialTraining(): TrainingDraft {
  return {
    startedAt: toLocalInput(new Date()),
    endedAt: "",
    memo: "",
    exercises: [createInitialExercise(1)],
  };
}

// ─── Exercise 操作 ────────────────────────────────

export function addExercise(draft: TrainingDraft): TrainingDraft {
  const nextOrder = draft.exercises.length + 1;
  return {
    ...draft,
    exercises: [...draft.exercises, createInitialExercise(nextOrder)],
  };
}

export function removeExercise(
  draft: TrainingDraft,
  index: number,
): TrainingDraft {
  const exercises = draft.exercises
    .filter((_, i) => i !== index)
    .map((ex, i) => ({ ...ex, displayOrder: i + 1 }));
  return { ...draft, exercises };
}

export function moveExercise(
  draft: TrainingDraft,
  from: number,
  to: number,
): TrainingDraft {
  if (from === to) return draft;
  const next = [...draft.exercises];
  const [moved] = next.splice(from, 1);
  next.splice(to, 0, moved);
  const exercises = next.map((ex, i) => ({ ...ex, displayOrder: i + 1 }));
  return { ...draft, exercises };
}

export function updateExercise(
  draft: TrainingDraft,
  index: number,
  patch: Partial<Omit<ExerciseDraft, "key" | "sets" | "displayOrder">>,
): TrainingDraft {
  const exercises = draft.exercises.map((ex, i) =>
    i === index ? { ...ex, ...patch } : ex,
  );
  return { ...draft, exercises };
}

// ─── Set 操作 ─────────────────────────────────────

export function addSet(
  draft: TrainingDraft,
  exerciseIndex: number,
): TrainingDraft {
  const exercises = draft.exercises.map((ex, i) => {
    if (i !== exerciseIndex) return ex;
    const nextNumber = ex.sets.length + 1;
    return { ...ex, sets: [...ex.sets, createInitialSet(nextNumber)] };
  });
  return { ...draft, exercises };
}

export function removeSet(
  draft: TrainingDraft,
  exerciseIndex: number,
  setIndex: number,
): TrainingDraft {
  const exercises = draft.exercises.map((ex, i) => {
    if (i !== exerciseIndex) return ex;
    const sets = ex.sets
      .filter((_, j) => j !== setIndex)
      .map((s, j) => ({ ...s, setNumber: j + 1 }));
    return { ...ex, sets };
  });
  return { ...draft, exercises };
}

export function updateSet(
  draft: TrainingDraft,
  exerciseIndex: number,
  setIndex: number,
  patch: Partial<Omit<SetDraft, "key" | "setNumber">>,
): TrainingDraft {
  const exercises = draft.exercises.map((ex, i) => {
    if (i !== exerciseIndex) return ex;
    const sets = ex.sets.map((s, j) => (j === setIndex ? { ...s, ...patch } : s));
    return { ...ex, sets };
  });
  return { ...draft, exercises };
}

// ─── Training 本体の更新 ──────────────────────────

export function updateTraining(
  draft: TrainingDraft,
  patch: Partial<Pick<TrainingDraft, "startedAt" | "endedAt" | "memo">>,
): TrainingDraft {
  return { ...draft, ...patch };
}

// ─── ISO / datetime-local 変換 ────────────────────

/** Date → "YYYY-MM-DDTHH:mm"(datetime-local 用) */
export function toLocalInput(d: Date): string {
  const pad = (n: number) => String(n).padStart(2, "0");
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`;
}

/** ISO 8601 → datetime-local(ローカルタイムゾーンで丸める) */
export function isoToLocalInput(iso: string): string {
  return toLocalInput(new Date(iso));
}

/** datetime-local("YYYY-MM-DDTHH:mm")→ ISO 8601 */
export function localInputToISO(local: string): string {
  if (!local) return "";
  return new Date(local).toISOString();
}

// ─── DTO ↔ Draft 変換 ────────────────────────────

/** Draft → 送信用 Request payload */
export function toRecordRequest(
  draft: TrainingDraft,
): RecordTrainingRequest {
  return {
    started_at: localInputToISO(draft.startedAt),
    ended_at: draft.endedAt ? localInputToISO(draft.endedAt) : undefined,
    memo: draft.memo ? draft.memo : undefined,
    exercises: draft.exercises.map((ex) => ({
      exercise_id: ex.exerciseID,
      display_order: ex.displayOrder,
      rest_seconds: ex.restSeconds ?? undefined,
      memo: ex.memo ? ex.memo : undefined,
      sets: ex.sets.map((s) => ({
        set_number: s.setNumber,
        // 未入力(空文字)は 0 として送る(VO が空文字を弾くため)
        weight_kg: s.weightKg || "0",
        reps: s.reps,
        rest_seconds: s.restSeconds ?? undefined,
        memo: s.memo ? s.memo : undefined,
      })),
    })),
  };
}

/** TrainingDTO(取得結果) → Draft(編集画面用) */
export function fromTrainingDTO(dto: TrainingDTO): TrainingDraft {
  return {
    startedAt: dto.started_at ? isoToLocalInput(dto.started_at) : toLocalInput(new Date()),
    endedAt: dto.ended_at ? isoToLocalInput(dto.ended_at) : "",
    memo: dto.memo ?? "",
    exercises: (dto.exercises ?? []).map(exerciseFromDTO),
  };
}

function exerciseFromDTO(dto: TrainingExerciseDTO): ExerciseDraft {
  return {
    key: nextKey("ex"),
    exerciseID: dto.exercise_id ?? "",
    displayOrder: dto.display_order ?? 1,
    restSeconds: dto.rest_seconds ?? null,
    memo: dto.memo ?? "",
    sets: (dto.sets ?? []).map(setFromDTO),
  };
}

function setFromDTO(dto: TrainingSetDTO): SetDraft {
  return {
    key: nextKey("set"),
    setNumber: dto.set_number ?? 1,
    weightKg: dto.weight_kg ?? "0",
    reps: dto.reps ?? 0,
    restSeconds: dto.rest_seconds ?? null,
    memo: dto.memo ?? "",
  };
}

// ─── 表示ヘルパー ─────────────────────────────────

export function formatDateTime(iso: string | undefined): string {
  if (!iso) return "";
  return new Date(iso).toLocaleString("ja-JP", {
    dateStyle: "short",
    timeStyle: "short",
  });
}

/**
 * セットの「実際の休憩時間」 を解決する。
 * セット自身が NULL なら親(種目)のデフォルト、 それも NULL なら null。
 */
export function resolveRestSeconds(
  setRest: number | null | undefined,
  exerciseDefaultRest: number | null | undefined,
): number | null {
  if (setRest !== null && setRest !== undefined) return setRest;
  if (exerciseDefaultRest !== null && exerciseDefaultRest !== undefined) {
    return exerciseDefaultRest;
  }
  return null;
}
