"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  apiClient,
  type UpsertMealTemplateRequest,
  type ReorderMealTemplatesRequest,
} from "@/shared/api/client";
import { toMealTemplate, type MealTemplate } from "../model/meal_template";

export const MEAL_TEMPLATES_QUERY_KEY = ["meal_templates"] as const;

export function useMealTemplatesQuery(enabled: boolean = true) {
  return useQuery({
    queryKey: MEAL_TEMPLATES_QUERY_KEY,
    enabled,
    queryFn: async (): Promise<MealTemplate[]> => {
      const { data, error, response } = await apiClient.GET(
        "/meal_templates",
        { params: { query: { limit: 100, offset: 0 } } },
      );
      if (error) {
        throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      }
      return (data?.meal_templates ?? []).map(toMealTemplate);
    },
  });
}

export function useCreateMealTemplateMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (body: UpsertMealTemplateRequest) => {
      const { data, error, response } = await apiClient.POST(
        "/meal_templates",
        { body },
      );
      if (error) {
        throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      }
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: MEAL_TEMPLATES_QUERY_KEY });
    },
  });
}

export function useUpdateMealTemplateMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async ({
      id,
      body,
    }: {
      id: string;
      body: UpsertMealTemplateRequest;
    }) => {
      const { data, error, response } = await apiClient.PUT(
        "/meal_templates/{id}",
        { params: { path: { id } }, body },
      );
      if (error) {
        throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      }
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: MEAL_TEMPLATES_QUERY_KEY });
    },
  });
}

export function useDeleteMealTemplateMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (id: string) => {
      const { error, response } = await apiClient.DELETE(
        "/meal_templates/{id}",
        { params: { path: { id } } },
      );
      if (error) {
        throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: MEAL_TEMPLATES_QUERY_KEY });
    },
  });
}

export function useReorderMealTemplatesMutation() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (body: ReorderMealTemplatesRequest) => {
      const { error, response } = await apiClient.POST(
        "/meal_templates/reorder",
        { body },
      );
      if (error) {
        throw new Error(error.error?.message ?? `HTTP ${response.status}`);
      }
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: MEAL_TEMPLATES_QUERY_KEY });
    },
  });
}
