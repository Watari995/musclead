"use client";

import { useState } from "react";
import { usePreferencesQuery, useUpdatePreferencesMutation } from "@/features/user/api/user";
import { useAccessToken } from "@/shared/auth/access-token";
import { ErrorText } from "@/shared/ui";

const PRESET_COLORS = [
  "#4A90E2", "#7ED321", "#FF6B6B", "#F5A623",
  "#BD10E0", "#50E3C2", "#B8E986", "#9013FE",
];

type ColorFieldProps = {
  label: string;
  value: string;
  onChange: (color: string) => void;
  disabled: boolean;
};

function ColorField({ label, value, onChange, disabled }: ColorFieldProps) {
  return (
    <div className="flex items-center justify-between">
      <div className="flex items-center gap-3">
        <span
          className="w-5 h-5 rounded-full border border-[var(--color-line)]"
          style={{ background: value }}
        />
        <span className="text-sm">{label}</span>
      </div>
      <div className="flex gap-1.5 flex-wrap justify-end">
        {PRESET_COLORS.map((c) => (
          <button
            key={c}
            type="button"
            disabled={disabled}
            onClick={() => onChange(c)}
            className={`w-6 h-6 rounded-full border-2 transition-all ${
              value === c ? "border-[var(--color-ink)] scale-110" : "border-transparent"
            }`}
            style={{ background: c }}
            aria-label={c}
          />
        ))}
        <input
          type="color"
          value={value}
          onChange={(e) => onChange(e.target.value)}
          disabled={disabled}
          className="w-6 h-6 rounded cursor-pointer border border-[var(--color-line)] p-0"
          title="カスタム色"
        />
      </div>
    </div>
  );
}

export function CalendarColorPicker() {
  const { token, ready } = useAccessToken();
  const { data: prefs } = usePreferencesQuery(ready && Boolean(token));
  const mutation = useUpdatePreferencesMutation();

  const [trainingColor, setTrainingColor] = useState<string | undefined>();
  const [mealColor, setMealColor] = useState<string | undefined>();
  const [weightColor, setWeightColor] = useState<string | undefined>();

  const effectiveTraining = trainingColor ?? prefs?.training_color ?? "#4A90E2";
  const effectiveMeal = mealColor ?? prefs?.meal_color ?? "#7ED321";
  const effectiveWeight = weightColor ?? prefs?.weight_color ?? "#FF6B6B";

  const handleChange = (field: "training_color" | "meal_color" | "weight_color", value: string) => {
    if (field === "training_color") setTrainingColor(value);
    if (field === "meal_color") setMealColor(value);
    if (field === "weight_color") setWeightColor(value);
    mutation.mutate({ [field]: value });
  };

  return (
    <div className="space-y-4">
      <ColorField
        label="トレーニング"
        value={effectiveTraining}
        onChange={(v) => handleChange("training_color", v)}
        disabled={mutation.isPending}
      />
      <ColorField
        label="食事"
        value={effectiveMeal}
        onChange={(v) => handleChange("meal_color", v)}
        disabled={mutation.isPending}
      />
      <ColorField
        label="体重"
        value={effectiveWeight}
        onChange={(v) => handleChange("weight_color", v)}
        disabled={mutation.isPending}
      />
      {mutation.isError && (
        <ErrorText>{(mutation.error as Error).message}</ErrorText>
      )}
    </div>
  );
}
