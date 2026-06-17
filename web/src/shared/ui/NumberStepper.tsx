import { NumberField } from "./NumberField";

type NumberStepperProps = {
  /** 現在値。 未入力は undefined */
  value: number | undefined;
  onChange: (value: number | undefined) => void;
  /** ±ボタン1回の増減幅(小数桁の丸めにも使う) */
  step: number;
  /** 必須: 例示値や単位を示すプレースホルダ */
  placeholder: string;
  min?: number;
  max?: number;
  /** 未入力時に ± を押したら開始する値(既定 0) */
  fallback?: number;
  /** ボタンの aria-label 用ラベル(例: 体重) */
  label?: string;
};

const stepBtn =
  "h-11 w-11 shrink-0 rounded-md border border-[var(--color-line)] bg-[var(--color-surface)] text-[var(--color-ink)] text-xl leading-none inline-flex items-center justify-center transition-colors hover:bg-[var(--color-surface-alt)] disabled:opacity-50 disabled:cursor-not-allowed";

/** step の小数桁で丸める(0.1 + 0.2 の浮動小数点誤差を防ぐ) */
function roundToStep(n: number, step: number): number {
  const decimals = (String(step).split(".")[1] ?? "").length;
  return Number(n.toFixed(decimals));
}

/**
 * 数値入力に −/＋ ボタンを付けたステッパー。
 * 体重のように前回からほぼ変わらない値を、 キーボード無しで微調整するためのもの。
 */
export function NumberStepper({
  value,
  onChange,
  step,
  placeholder,
  min,
  max,
  fallback = 0,
  label,
}: NumberStepperProps) {
  const current = value ?? fallback;

  const nudge = (delta: number) => {
    let next = roundToStep(current + delta, step);
    if (min !== undefined && next < min) next = min;
    if (max !== undefined && next > max) next = max;
    onChange(next);
  };

  return (
    <div className="flex items-stretch gap-2">
      <button
        type="button"
        className={stepBtn}
        onClick={() => nudge(-step)}
        disabled={min !== undefined && current <= min}
        aria-label={label ? `${label}を減らす` : "減らす"}
      >
        −
      </button>
      <NumberField
        value={value}
        onChange={onChange}
        step={step}
        min={min}
        max={max}
        placeholder={placeholder}
        className="text-center"
        aria-label={label}
      />
      <button
        type="button"
        className={stepBtn}
        onClick={() => nudge(step)}
        disabled={max !== undefined && current >= max}
        aria-label={label ? `${label}を増やす` : "増やす"}
      >
        ＋
      </button>
    </div>
  );
}
