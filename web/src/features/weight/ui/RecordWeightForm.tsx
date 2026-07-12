"use client";

import { useState } from "react";
import { useTranslations } from "next-intl";
import {
  useRecordWeightMutation,
  useWeightsQuery,
  type RecordWeightRequest,
} from "@/features/weight/api/weights";
import { toLocalInput } from "@/features/weight/model/weight";
import {
  Button,
  Card,
  ErrorText,
  Label,
  NumberStepper,
  TextInput,
} from "@/shared/ui";

type FormState = {
  weightKg: number | undefined;
  bodyFatPercentage: number | undefined;
  skeletalMuscleKg: number | undefined;
  measuredAt: string;
};

export function RecordWeightForm() {
  const t = useTranslations("weights");
  const tCommon = useTranslations("common");
  const [form, setForm] = useState<FormState>(emptyForm);
  const [success, setSuccess] = useState(false);
  const recordMutation = useRecordWeightMutation();

  // 前回値プリフィル: 最新記録(weights[0]、 measured_at DESC なので先頭) を
  // 一度だけ初期値に流し込む。 体重は前日からほぼ変わらないので入力が ± だけで済む。
  // effect ではなくレンダー中に state を調整する React 公式パターン。 prefilled を
  // ガードにして一度きりにし、 背景 refetch で入力中の値が上書きされるのを防ぐ。
  const { data: weights } = useWeightsQuery();
  const [prefilled, setPrefilled] = useState(false);
  const last = weights?.[0];
  if (!prefilled && last) {
    setPrefilled(true);
    setForm((f) => ({
      ...f,
      weightKg: toNum(last.weightKg),
      bodyFatPercentage: toNum(last.bodyFatPercentage),
      skeletalMuscleKg: toNum(last.skeletalMuscleKg),
    }));
  }

  return (
    <Card className="p-4 sm:p-5">
      <form
        className="space-y-4"
        onSubmit={(e) => {
          e.preventDefault();
          setSuccess(false);
          if (form.weightKg === undefined) return;
          const body: RecordWeightRequest = {
            weight_kg: String(form.weightKg),
            measured_at: new Date(form.measuredAt).toISOString(),
          };
          if (form.bodyFatPercentage !== undefined) {
            body.body_fat_percentage = String(form.bodyFatPercentage);
          }
          if (form.skeletalMuscleKg !== undefined) {
            body.skeletal_muscle_kg = String(form.skeletalMuscleKg);
          }
          recordMutation.mutate(body, {
            onSuccess: () => {
              setSuccess(true);
              // 入力値は次回の前回値になるので保持し、 日時だけ現在に更新する
              setForm((f) => ({ ...f, measuredAt: toLocalInput(new Date()) }));
            },
          });
        }}
      >
        <Label label={t("weightKg")}>
          <NumberStepper
            value={form.weightKg}
            onChange={(v) => setForm({ ...form, weightKg: v })}
            step={0.1}
            min={0}
            placeholder={t("exampleWeight")}
            label={t("typeWeight")}
          />
        </Label>
        <Label label={t("bodyFat")}>
          <NumberStepper
            value={form.bodyFatPercentage}
            onChange={(v) => setForm({ ...form, bodyFatPercentage: v })}
            step={0.1}
            min={0}
            max={100}
            placeholder={t("exampleBodyFat")}
            label={t("typeBodyFat")}
          />
        </Label>
        <Label label={t("muscleMass")}>
          <NumberStepper
            value={form.skeletalMuscleKg}
            onChange={(v) => setForm({ ...form, skeletalMuscleKg: v })}
            step={0.1}
            min={0}
            placeholder={t("exampleMuscle")}
            label={t("typeMuscle")}
          />
        </Label>
        <Label label={tCommon("dateTime")}>
          <TextInput
            type="datetime-local"
            value={form.measuredAt}
            onChange={(e) => setForm({ ...form, measuredAt: e.target.value })}
            required
          />
        </Label>
        {recordMutation.isError && (
          <ErrorText>{(recordMutation.error as Error).message}</ErrorText>
        )}
        {success && (
          <p className="text-sm text-[var(--color-ink-muted)]">{tCommon("recorded")}</p>
        )}
        <Button
          type="submit"
          fullWidth
          disabled={recordMutation.isPending || form.weightKg === undefined}
        >
          {recordMutation.isPending ? tCommon("recording") : tCommon("record")}
        </Button>
      </form>
    </Card>
  );
}

function emptyForm(): FormState {
  return {
    weightKg: undefined,
    bodyFatPercentage: undefined,
    skeletalMuscleKg: undefined,
    measuredAt: toLocalInput(new Date()),
  };
}

/** API の文字列値(decimal) を number | undefined に変換 */
function toNum(v: string | null | undefined): number | undefined {
  if (v === null || v === undefined || v === "") return undefined;
  const n = Number(v);
  return Number.isNaN(n) ? undefined : n;
}
