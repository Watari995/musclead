"use client";

import { useState } from "react";
import {
  useRecordWeightMutation,
  type RecordWeightRequest,
} from "@/features/weight/api/weights";
import { toLocalInput } from "@/features/weight/model/weight";
import { Button, Card, ErrorText, Label, TextInput } from "@/shared/ui";

type FormState = {
  weight_kg: string;
  body_fat_percentage: string;
  skeletal_muscle_kg: string;
  measured_at: string;
};

export function RecordWeightForm() {
  const [form, setForm] = useState<FormState>(initialForm);
  const [success, setSuccess] = useState(false);
  const recordMutation = useRecordWeightMutation();

  return (
    <Card className="p-4 sm:p-5">
      <form
        className="space-y-4"
        onSubmit={(e) => {
          e.preventDefault();
          setSuccess(false);
          const body: RecordWeightRequest = {
            weight_kg: form.weight_kg,
            measured_at: new Date(form.measured_at).toISOString(),
          };
          if (form.body_fat_percentage.trim() !== "") {
            body.body_fat_percentage = form.body_fat_percentage;
          }
          if (form.skeletal_muscle_kg.trim() !== "") {
            body.skeletal_muscle_kg = form.skeletal_muscle_kg;
          }
          recordMutation.mutate(body, {
            onSuccess: () => {
              setSuccess(true);
              setForm(initialForm());
            },
          });
        }}
      >
        <Label label="体重 (kg)">
          <TextInput
            type="number"
            step="0.01"
            min={0}
            value={form.weight_kg}
            onChange={(e) => setForm({ ...form, weight_kg: e.target.value })}
            required
          />
        </Label>
        <Label label="日時">
          <TextInput
            type="datetime-local"
            value={form.measured_at}
            onChange={(e) => setForm({ ...form, measured_at: e.target.value })}
            required
          />
        </Label>
        <Label label="体脂肪率 (%) ※任意">
          <TextInput
            type="number"
            step="0.01"
            min={0}
            max={100}
            value={form.body_fat_percentage}
            onChange={(e) =>
              setForm({ ...form, body_fat_percentage: e.target.value })
            }
          />
        </Label>
        <Label label="骨格筋量 (kg) ※任意">
          <TextInput
            type="number"
            step="0.01"
            min={0}
            value={form.skeletal_muscle_kg}
            onChange={(e) =>
              setForm({ ...form, skeletal_muscle_kg: e.target.value })
            }
          />
        </Label>
        {recordMutation.isError && (
          <ErrorText>{(recordMutation.error as Error).message}</ErrorText>
        )}
        {success && (
          <p className="text-sm text-[var(--color-ink-muted)]">
            記録しました
          </p>
        )}
        <Button type="submit" fullWidth disabled={recordMutation.isPending}>
          {recordMutation.isPending ? "記録中…" : "記録する"}
        </Button>
      </form>
    </Card>
  );
}

function initialForm(): FormState {
  return {
    weight_kg: "",
    body_fat_percentage: "",
    skeletal_muscle_kg: "",
    measured_at: toLocalInput(new Date()),
  };
}
