import { forwardRef, type ComponentPropsWithoutRef } from "react";
import { TextInput } from "./TextInput";

type NumberFieldProps = Omit<
  ComponentPropsWithoutRef<"input">,
  "value" | "onChange" | "type"
> & {
  /** 数値。未入力は undefined で表す(0 を初期値に使わない) */
  value: number | undefined;
  /** 未入力(空欄)になったときは undefined を返す */
  onChange: (value: number | undefined) => void;
  /** 必須: フィールドごとに例示値や単位を渡す(共通の既定値は持たせない) */
  placeholder: string;
};

/**
 * 数値入力。 空欄は常に "" として表示し、 0 を出さない
 * (フォーカスのたびに 0 を消す手間をなくす)。
 * 詳細は .cursor/rules/10-web-design-system.mdc「数値入力」を参照。
 */
export const NumberField = forwardRef<HTMLInputElement, NumberFieldProps>(
  function NumberField({ value, onChange, ...props }, ref) {
    return (
      <TextInput
        ref={ref}
        {...props}
        type="number"
        value={value ?? ""}
        onChange={(e) => {
          const raw = e.target.value;
          onChange(raw === "" ? undefined : Number(raw));
        }}
      />
    );
  },
);
