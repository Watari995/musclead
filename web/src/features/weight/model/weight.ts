export type Weight = {
  id: string;
  weightKg: string;
  bodyFatPercentage: string | null;
  skeletalMuscleKg: string | null;
  measuredAt: string;
  createdAt: string;
  updatedAt: string;
};

export function toLocalInput(d: Date): string {
  const pad = (n: number) => String(n).padStart(2, "0");
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`;
}
