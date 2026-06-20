import { useMutation, useQuery } from "@tanstack/react-query";
import { APIError } from "@/shared/api/client";
import { getAccessToken } from "@/shared/auth/access-token";
import type { CreateFoodProductRequest, FoodProduct, FoodProductDTO } from "../model/food_product";
import { toFoodProduct } from "../model/food_product";

const BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

async function foodFetch<T>(path: string, init?: RequestInit): Promise<T> {
  const token = getAccessToken();
  const res = await fetch(`${BASE_URL}${path}`, {
    ...init,
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...(init?.headers ?? {}),
    },
  });
  if (!res.ok) {
    const body = await res.json().catch(() => undefined);
    throw new APIError(res.status, body);
  }
  return res.json();
}

export function useFoodProductsByNameQuery(name: string) {
  return useQuery({
    queryKey: ["food_products", "name", name],
    enabled: name.trim().length >= 1,
    queryFn: async (): Promise<FoodProduct[]> => {
      const data = await foodFetch<{ food_products: FoodProductDTO[] }>(
        `/food_products?q=${encodeURIComponent(name.trim())}`,
      );
      return (data.food_products ?? []).map(toFoodProduct);
    },
  });
}

export function useFoodProductsByBarcodeQuery(barcode: string, enabled: boolean) {
  return useQuery({
    queryKey: ["food_products", "barcode", barcode],
    enabled: enabled && barcode.trim().length >= 8,
    retry: false,
    queryFn: async (): Promise<FoodProduct[]> => {
      try {
        const data = await foodFetch<{ food_products: FoodProductDTO[] }>(
          `/food_products/barcode/${barcode.trim()}`,
        );
        return (data.food_products ?? []).map(toFoodProduct);
      } catch (err) {
        if (err instanceof APIError && err.status === 404) return [];
        throw err;
      }
    },
  });
}

export function useCreateFoodProductMutation() {
  return useMutation({
    mutationFn: async (body: CreateFoodProductRequest): Promise<string> => {
      const data = await foodFetch<{ food_product_id: string }>("/food_products", {
        method: "POST",
        body: JSON.stringify(body),
      });
      return data.food_product_id;
    },
  });
}
