export type FoodProduct = {
  id: string;
  barcode: string | null;
  name: string;
  calories: number;
  proteinG: string | null;
  fatG: string | null;
  carbohydrateG: string | null;
  registerSource: "open_food_facts" | "user";
};

export type FoodProductDTO = {
  id: string;
  barcode?: string;
  name: string;
  calories: number;
  protein_g?: string;
  fat_g?: string;
  carbohydrate_g?: string;
  register_source: string;
};

export type CreateFoodProductRequest = {
  barcode?: string;
  name: string;
  calories: number;
  protein_g?: string;
  fat_g?: string;
  carbohydrate_g?: string;
};

export function toFoodProduct(dto: FoodProductDTO): FoodProduct {
  return {
    id: dto.id,
    barcode: dto.barcode ?? null,
    name: dto.name,
    calories: dto.calories,
    proteinG: dto.protein_g ?? null,
    fatG: dto.fat_g ?? null,
    carbohydrateG: dto.carbohydrate_g ?? null,
    registerSource: dto.register_source as "open_food_facts" | "user",
  };
}
