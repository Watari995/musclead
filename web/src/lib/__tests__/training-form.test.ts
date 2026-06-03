import { describe, expect, it } from "vitest";
import {
  addExercise,
  addSet,
  createInitialExercise,
  createInitialSet,
  createInitialTraining,
  fromTrainingDTO,
  isoToLocalInput,
  localInputToISO,
  moveExercise,
  removeExercise,
  removeSet,
  resolveRestSeconds,
  toRecordRequest,
  updateExercise,
  updateSet,
  updateTraining,
} from "../training-form";

describe("createInitialTraining", () => {
  it("starts with a single exercise that has a single set", () => {
    const draft = createInitialTraining();
    expect(draft.exercises).toHaveLength(1);
    expect(draft.exercises[0].sets).toHaveLength(1);
    expect(draft.exercises[0].displayOrder).toBe(1);
    expect(draft.exercises[0].sets[0].setNumber).toBe(1);
  });

  it("gives each exercise a unique React key", () => {
    const a = createInitialExercise(1);
    const b = createInitialExercise(2);
    expect(a.key).not.toEqual(b.key);
  });

  it("gives each set a unique React key", () => {
    const a = createInitialSet(1);
    const b = createInitialSet(2);
    expect(a.key).not.toEqual(b.key);
  });
});

describe("exercise operations", () => {
  it("addExercise appends with next display order", () => {
    const draft = createInitialTraining();
    const after = addExercise(draft);
    expect(after.exercises).toHaveLength(2);
    expect(after.exercises[1].displayOrder).toBe(2);
  });

  it("removeExercise drops by index and re-numbers the remaining order", () => {
    const draft = addExercise(addExercise(createInitialTraining())); // 3 exercises
    const after = removeExercise(draft, 1);
    expect(after.exercises).toHaveLength(2);
    expect(after.exercises.map((e) => e.displayOrder)).toEqual([1, 2]);
  });

  it("moveExercise rotates and re-numbers", () => {
    let draft = createInitialTraining();
    draft = updateExercise(draft, 0, { exerciseID: "A" });
    draft = addExercise(draft);
    draft = updateExercise(draft, 1, { exerciseID: "B" });
    draft = addExercise(draft);
    draft = updateExercise(draft, 2, { exerciseID: "C" });

    const moved = moveExercise(draft, 0, 2);
    expect(moved.exercises.map((e) => e.exerciseID)).toEqual(["B", "C", "A"]);
    expect(moved.exercises.map((e) => e.displayOrder)).toEqual([1, 2, 3]);
  });

  it("updateExercise patches only the targeted exercise", () => {
    const draft = addExercise(createInitialTraining());
    const after = updateExercise(draft, 1, { exerciseID: "ex-uuid-squat" });
    expect(after.exercises[0].exerciseID).toBe("");
    expect(after.exercises[1].exerciseID).toBe("ex-uuid-squat");
  });
});

describe("set operations", () => {
  it("addSet appends with next set number under the right exercise", () => {
    const draft = addExercise(createInitialTraining());
    const after = addSet(draft, 1);
    expect(after.exercises[0].sets).toHaveLength(1);
    expect(after.exercises[1].sets).toHaveLength(2);
    expect(after.exercises[1].sets[1].setNumber).toBe(2);
  });

  it("removeSet drops by index and re-numbers", () => {
    let draft = addSet(addSet(createInitialTraining(), 0), 0); // 3 sets
    draft = updateSet(draft, 0, 0, { weightKg: "60.00" });
    draft = updateSet(draft, 0, 1, { weightKg: "65.00" });
    draft = updateSet(draft, 0, 2, { weightKg: "70.00" });

    const after = removeSet(draft, 0, 1);
    expect(after.exercises[0].sets).toHaveLength(2);
    expect(after.exercises[0].sets.map((s) => s.setNumber)).toEqual([1, 2]);
    expect(after.exercises[0].sets.map((s) => s.weightKg)).toEqual([
      "60.00",
      "70.00",
    ]);
  });

  it("updateSet patches only the targeted set", () => {
    const draft = addSet(createInitialTraining(), 0);
    const after = updateSet(draft, 0, 1, { weightKg: "80.00", reps: 5 });
    expect(after.exercises[0].sets[0].weightKg).toBe("0");
    expect(after.exercises[0].sets[1].weightKg).toBe("80.00");
    expect(after.exercises[0].sets[1].reps).toBe(5);
  });
});

describe("updateTraining", () => {
  it("patches the top-level fields only", () => {
    const draft = createInitialTraining();
    const after = updateTraining(draft, { memo: "morning push day" });
    expect(after.memo).toBe("morning push day");
    expect(after.exercises).toBe(draft.exercises);
  });
});

describe("datetime helpers", () => {
  it("localInputToISO produces an ISO timestamp", () => {
    const iso = localInputToISO("2026-06-02T18:00");
    expect(iso).toMatch(/^2026-06-02T\d{2}:00:00\.000Z$/);
  });

  it("isoToLocalInput round-trips through datetime-local", () => {
    const original = localInputToISO("2026-06-02T18:30");
    const back = isoToLocalInput(original);
    expect(back).toBe("2026-06-02T18:30");
  });

  it("localInputToISO returns empty string for empty input", () => {
    expect(localInputToISO("")).toBe("");
  });
});

describe("toRecordRequest", () => {
  it("maps draft into the API request shape and drops empty optional fields", () => {
    let draft = createInitialTraining();
    draft = updateTraining(draft, { startedAt: "2026-06-02T18:00" });
    draft = updateExercise(draft, 0, {
      exerciseID: "ex-uuid-bench",
      restSeconds: 90,
    });
    draft = updateSet(draft, 0, 0, {
      weightKg: "60.00",
      reps: 10,
    });

    const req = toRecordRequest(draft);
    expect(req.started_at).toMatch(/^2026-06-02T/);
    expect(req.ended_at).toBeUndefined();
    expect(req.memo).toBeUndefined();
    expect(req.exercises).toHaveLength(1);
    expect(req.exercises![0]).toMatchObject({
      exercise_id: "ex-uuid-bench",
      display_order: 1,
      rest_seconds: 90,
    });
    expect(req.exercises![0].sets).toEqual([
      {
        set_number: 1,
        weight_kg: "60.00",
        reps: 10,
      },
    ]);
  });

  it("preserves non-null memo / ended_at when populated", () => {
    let draft = createInitialTraining();
    draft = updateTraining(draft, {
      startedAt: "2026-06-02T18:00",
      endedAt: "2026-06-02T19:00",
      memo: "good day",
    });

    const req = toRecordRequest(draft);
    expect(req.memo).toBe("good day");
    expect(req.ended_at).toMatch(/^2026-06-02T/);
  });
});

describe("fromTrainingDTO", () => {
  it("rehydrates DTO into a draft retaining hierarchy and memos", () => {
    const draft = fromTrainingDTO({
      id: "tid",
      user_id: "uid",
      started_at: "2026-06-02T09:00:00Z",
      ended_at: "2026-06-02T10:00:00Z",
      memo: "rest",
      created_at: "2026-06-02T09:00:00Z",
      updated_at: "2026-06-02T10:00:00Z",
      exercises: [
        {
          id: "ex1",
          exercise_id: "ex-uuid-squat",
          display_order: 1,
          rest_seconds: 120,
          memo: undefined,
          sets: [
            {
              id: "s1",
              set_number: 1,
              weight_kg: "80.00",
              reps: 5,
              rest_seconds: undefined,
              memo: undefined,
            },
          ],
        },
      ],
    });

    expect(draft.memo).toBe("rest");
    expect(draft.exercises).toHaveLength(1);
    expect(draft.exercises[0].exerciseID).toBe("ex-uuid-squat");
    expect(draft.exercises[0].sets).toHaveLength(1);
    expect(draft.exercises[0].sets[0].weightKg).toBe("80.00");
  });

  it("falls back to safe defaults when fields are missing", () => {
    const draft = fromTrainingDTO({});
    expect(draft.exercises).toEqual([]);
    expect(draft.memo).toBe("");
    expect(draft.endedAt).toBe("");
  });
});

describe("resolveRestSeconds", () => {
  it("prefers the set value when present", () => {
    expect(resolveRestSeconds(45, 90)).toBe(45);
  });

  it("falls back to the exercise default when the set is null", () => {
    expect(resolveRestSeconds(null, 90)).toBe(90);
  });

  it("returns null when neither is set", () => {
    expect(resolveRestSeconds(null, null)).toBeNull();
    expect(resolveRestSeconds(undefined, undefined)).toBeNull();
  });
});
