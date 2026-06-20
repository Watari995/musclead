"use client";

import { useEffect, useRef, useState } from "react";
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
            className={`px-3 py-1 rounded-full text-xs border transition-colors ${
              mode === m
                ? "bg-[var(--color-ink)] text-[var(--color-bg)] border-[var(--color-ink)]"
                : "border-[var(--color-line)] text-[var(--color-ink-muted)] hover:border-[var(--color-ink)]"
            }`}
          >
            {m === "name" ? "名前で検索" : "バーコード"}
          </button>
        ))}
      </div>

      {mode === "name" ? (
        <TextInput
          placeholder="例: おにぎり、プロテインバー…"
          value={nameQuery}
          onChange={handleNameChange}
        />
      ) : (
        <div className="flex gap-2">
          <TextInput
            placeholder="バーコードの数字を入力 (8〜13桁)"
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
            className="shrink-0 px-4 py-2 rounded-md bg-[var(--color-ink)] text-[var(--color-bg)] text-sm font-medium disabled:opacity-40 transition-opacity"
          >
            検索
          </button>
        </div>
      )}

      {isLoading && (
        <p className="text-xs text-[var(--color-ink-muted)]">検索中…</p>
      )}

      {isError && (
        <ErrorText>検索に失敗しました。もう一度お試しください。</ErrorText>
      )}

      {showResults && (
        <ul className="border border-[var(--color-line)] rounded-lg overflow-hidden bg-[var(--color-surface)] divide-y divide-[var(--color-line)]">
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
            見つかりませんでした
          </p>
          <button
            type="button"
            onClick={() =>
              onNotFound(mode === "barcode" ? barcodeInput.trim() : undefined)
            }
            className="text-xs text-[var(--color-ink)] underline"
          >
            登録する →
          </button>
        </div>
      )}
    </div>
  );
}
