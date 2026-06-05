import type { ExerciseDTO } from "@/shared/api/client";

export type Exercise = {
  id: string;
  name: string;
  createdAt: string;
  updatedAt: string;
};

export function toExercise(dto: ExerciseDTO): Exercise {
  return {
    id: dto.id ?? "",
    name: dto.name ?? "",
    createdAt: dto.created_at ?? "",
    updatedAt: dto.updated_at ?? "",
  };
}
