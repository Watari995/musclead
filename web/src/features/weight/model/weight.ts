export type Weight = {
  id: string;
  weightKg: string;
  bodyFatPercentage: string | null;
  skeletalMuscleKg: string | null;
  measuredAt: string;
  createdAt: string;
  updatedAt: string;
};

export type WeightDTO = {
  id: string;
  weight_kg: string;
  body_fat_percentage?: string | null;
  skeletal_muscle_kg?: string | null;
  measured_at: string;
  created_at: string;
  updated_at: string;
};

export function toWeight(dto: WeightDTO): Weight {
  return {
    id: dto.id,
    weightKg: dto.weight_kg,
    bodyFatPercentage: dto.body_fat_percentage ?? null,
    skeletalMuscleKg: dto.skeletal_muscle_kg ?? null,
    measuredAt: dto.measured_at,
    createdAt: dto.created_at,
    updatedAt: dto.updated_at,
  };
}

export function toLocalInput(d: Date): string {
  const pad = (n: number) => String(n).padStart(2, "0");
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`;
}

export function formatWeightDateTime(iso: string): string {
  if (!iso) return "";
  return new Date(iso).toLocaleString("ja-JP", {
    dateStyle: "short",
    timeStyle: "short",
  });
}
