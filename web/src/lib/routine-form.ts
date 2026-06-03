// routine フォームの素材。 Training と違い 1 階層(exercise 順)だけなので軽量。

import type {
  RoutineDTO,
  RoutineExerciseDTO,
  UpsertRoutineRequest,
} from "@/api/client";

export type RoutineExerciseDraft = {
  /** React key 用のローカル ID(永続化されない) */
  key: string;
  exerciseID: string;
  displayOrder: number;
};

export type RoutineDraft = {
  name: string;
  exercises: RoutineExerciseDraft[];
};

let keySeed = 0;
function nextKey(prefix: string): string {
  keySeed += 1;
  return `${prefix}-${keySeed}-${Math.random().toString(36).slice(2, 8)}`;
}

export function createInitialExercise(
  displayOrder: number,
): RoutineExerciseDraft {
  return { key: nextKey("rex"), exerciseID: "", displayOrder };
}

export function createInitialRoutine(): RoutineDraft {
  return { name: "", exercises: [createInitialExercise(1)] };
}

export function addExercise(draft: RoutineDraft): RoutineDraft {
  const nextOrder = draft.exercises.length + 1;
  return {
    ...draft,
    exercises: [...draft.exercises, createInitialExercise(nextOrder)],
  };
}

export function removeExercise(
  draft: RoutineDraft,
  index: number,
): RoutineDraft {
  const exercises = draft.exercises
    .filter((_, i) => i !== index)
    .map((ex, i) => ({ ...ex, displayOrder: i + 1 }));
  return { ...draft, exercises };
}

export function moveExercise(
  draft: RoutineDraft,
  from: number,
  to: number,
): RoutineDraft {
  if (from === to) return draft;
  const next = [...draft.exercises];
  const [moved] = next.splice(from, 1);
  next.splice(to, 0, moved);
  const exercises = next.map((ex, i) => ({ ...ex, displayOrder: i + 1 }));
  return { ...draft, exercises };
}

export function setExerciseID(
  draft: RoutineDraft,
  index: number,
  exerciseID: string,
): RoutineDraft {
  const exercises = draft.exercises.map((ex, i) =>
    i === index ? { ...ex, exerciseID } : ex,
  );
  return { ...draft, exercises };
}

export function setName(draft: RoutineDraft, name: string): RoutineDraft {
  return { ...draft, name };
}

/** Draft → 送信用 Request payload */
export function toUpsertRequest(draft: RoutineDraft): UpsertRoutineRequest {
  return {
    name: draft.name,
    exercises: draft.exercises.map((ex) => ({
      exercise_id: ex.exerciseID,
      display_order: ex.displayOrder,
    })),
  };
}

/** RoutineDTO(取得結果) → Draft(編集画面用) */
export function fromRoutineDTO(dto: RoutineDTO): RoutineDraft {
  return {
    name: dto.name ?? "",
    exercises: (dto.routine_exercises ?? []).map(exerciseFromDTO),
  };
}

function exerciseFromDTO(dto: RoutineExerciseDTO): RoutineExerciseDraft {
  return {
    key: nextKey("rex"),
    exerciseID: dto.exercise_id ?? "",
    displayOrder: dto.display_order ?? 1,
  };
}
