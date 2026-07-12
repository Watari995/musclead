"use client";

import { useEffect, useRef, useState } from "react";
import { useTranslations } from "next-intl";
import { ErrorText, TextInput } from "@/shared/ui";
import {
  useFoodProductsByBarcodeQuery,
  useFoodProductsByNameQuery,
} from "../api/food_products";
import type { FoodProduct } from "../model/food_product";

type Props = {
  onSelect: (food: FoodProduct) => void;
  onNotFound: (barcode?: string) => void;
};

type Mode = "name" | "barcode";

export function FoodSearchSection({ onSelect, onNotFound }: Props) {
  const t = useTranslations("food");
  const [mode, setMode] = useState<Mode>("name");
  const [nameQuery, setNameQuery] = useState("");
  const [barcodeInput, setBarcodeInput] = useState("");
  const [barcodeEnabled, setBarcodeEnabled] = useState(false);
  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const nameResult = useFoodProductsByNameQuery(nameQuery);
  const barcodeResult = useFoodProductsByBarcodeQuery(barcodeInput, barcodeEnabled);

  const handleNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const val = e.target.value;
    setNameQuery(val);
    if (debounceRef.current) clearTimeout(debounceRef.current);
    debounceRef.current = setTimeout(() => {}, 0);
  };

  useEffect(() => {
    return () => {
      if (debounceRef.current) clearTimeout(debounceRef.current);
    };
  }, []);

  const handleModeChange = (m: Mode) => {
    setMode(m);
    setNameQuery("");
    setBarcodeInput("");
    setBarcodeEnabled(false);
  };

  const handleBarcodeSearch = () => {
    setBarcodeEnabled(true);
  };

  const handleSelect = (food: FoodProduct) => {
    setNameQuery("");
    setBarcodeInput("");
    setBarcodeEnabled(false);
    onSelect(food);
  };

  const activeResults =
    mode === "name" ? (nameResult.data ?? []) : (barcodeResult.data ?? []);
  const isLoading =
    mode === "name" ? nameResult.isLoading : barcodeResult.isFetching;
  const isError =
    mode === "name" ? nameResult.isError : barcodeResult.isError;
  const isFetched =
    mode === "name" ? nameResult.isFetched : barcodeResult.isFetched;

  const showResults = activeResults.length > 0;
  const showEmpty =
    isFetched &&
    !isLoading &&
    activeResults.length === 0 &&
    (mode === "name" ? nameQuery.trim().length > 0 : barcodeEnabled);

  return (
    <div className="space-y-3">
      <div className="flex gap-2">
        {(["name", "barcode"] as const).map((m) => (
          <button
            key={m}
            type="button"
            onClick={() => handleModeChange(m)}
            className={`rough rough-pill px-3 py-1 text-xs transition-colors ${
              mode === m
                ? "bg-[var(--color-ink)] text-[var(--color-surface)]"
                : "text-[var(--color-ink-muted)]"
            }`}
          >
            {m === "name" ? t("searchByName") : t("barcode")}
          </button>
        ))}
      </div>

      {mode === "name" ? (
        <TextInput
          placeholder={t("namePlaceholder")}
          value={nameQuery}
          onChange={handleNameChange}
        />
      ) : (
        <div className="flex gap-2">
          <TextInput
            placeholder={t("barcodePlaceholder")}
            value={barcodeInput}
            onChange={(e) => {
              setBarcodeInput(e.target.value);
              setBarcodeEnabled(false);
            }}
            type="tel"
          />
          <button
            type="button"
            onClick={handleBarcodeSearch}
            disabled={barcodeInput.trim().length < 8}
            className="rough shrink-0 px-4 py-2 bg-[var(--color-ink)] text-[var(--color-surface)] text-sm font-medium disabled:opacity-40 transition-opacity"
          >
            {t("search")}
          </button>
        </div>
      )}

      {isLoading && (
        <p className="text-xs text-[var(--color-ink-muted)]">{t("searching")}</p>
      )}

      {isError && (
        <ErrorText>{t("searchFailed")}</ErrorText>
      )}

      {showResults && (
        <ul className="rough overflow-hidden bg-[var(--color-surface)] divide-y divide-[var(--color-line)]">
          {activeResults.map((food) => (
            <li key={food.id}>
              <button
                type="button"
                onClick={() => handleSelect(food)}
                className="w-full text-left px-4 py-3 hover:bg-[var(--color-surface-alt)] transition-colors"
              >
                <p className="text-sm font-semibold">{food.name}</p>
                <p className="text-xs text-[var(--color-ink-muted)] mt-0.5">
                  {food.calories} kcal
                  {food.proteinG && ` · P ${food.proteinG}g`}
                  {food.fatG && ` · F ${food.fatG}g`}
                  {food.carbohydrateG && ` · C ${food.carbohydrateG}g`}
                </p>
              </button>
            </li>
          ))}
        </ul>
      )}

      {showEmpty && (
        <div className="flex items-center justify-between py-1">
          <p className="text-xs text-[var(--color-ink-muted)]">
            {t("notFound")}
          </p>
          <button
            type="button"
            onClick={() =>
              onNotFound(mode === "barcode" ? barcodeInput.trim() : undefined)
            }
            className="text-xs text-[var(--color-ink)] underline"
          >
            {t("registerLink")}
          </button>
        </div>
      )}
    </div>
  );
}
