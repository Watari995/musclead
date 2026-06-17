"use client";

// 新規トレーニングの下書きを localStorage に自動退避するフック。
// アプリをキルしても端末に残り、 開き直すと「復元しますか?」で戻せる。
// 提出までの一時退避なので、 提出成功時に clear() で消す。

import { useEffect, useRef, useState } from "react";
import {
  type TrainingDraft,
  NEW_TRAINING_DRAFT_KEY,
  createInitialTraining,
  deserializeDraft,
  serializeDraft,
} from "./training-draft";

const AUTOSAVE_DELAY_MS = 500;

type UseNewTrainingDraft = {
  draft: TrainingDraft;
  setDraft: (draft: TrainingDraft) => void;
  /** 端末に残っていた下書きがあるか(復元バナー表示の判定) */
  restorable: boolean;
  /** 残っていた下書きをフォームへ反映する */
  restore: () => void;
  /** 残っていた下書きを破棄して新規から始める */
  discard: () => void;
  /** 下書きを消す(提出成功時など)。 以後の自動退避も止める。 */
  clear: () => void;
};

export function useNewTrainingDraft(): UseNewTrainingDraft {
  const [draft, setDraft] = useState<TrainingDraft>(createInitialTraining);
  const [restorable, setRestorable] = useState<TrainingDraft | null>(null);
  // 復元の判断が済むまで自動退避しない(残っている下書きの誤上書き防止)
  const [armed, setArmed] = useState(false);
  const stoppedRef = useRef(false);

  // マウント時に既存下書きを読む(自動適用はせず、 ユーザーに復元を委ねる)。
  // localStorage は SSR に無いため mount 後に読む。 外部ストアの初回同期なので
  // set-state-in-effect は許容(ThemePicker と同じ方針)。
  useEffect(() => {
    const saved = deserializeDraft(
      localStorage.getItem(NEW_TRAINING_DRAFT_KEY),
    );
    // eslint-disable-next-line react-hooks/set-state-in-effect
    setRestorable(saved);
    if (!saved) setArmed(true);
  }, []);

  // 自動退避(debounce)
  useEffect(() => {
    if (!armed || stoppedRef.current) return;
    const id = setTimeout(() => {
      if (stoppedRef.current) return;
      localStorage.setItem(NEW_TRAINING_DRAFT_KEY, serializeDraft(draft));
    }, AUTOSAVE_DELAY_MS);
    return () => clearTimeout(id);
  }, [draft, armed]);

  const restore = () => {
    if (restorable) setDraft(restorable);
    setRestorable(null);
    setArmed(true);
  };

  const discard = () => {
    localStorage.removeItem(NEW_TRAINING_DRAFT_KEY);
    setRestorable(null);
    setArmed(true);
  };

  const clear = () => {
    stoppedRef.current = true;
    localStorage.removeItem(NEW_TRAINING_DRAFT_KEY);
  };

  return {
    draft,
    setDraft,
    restorable: restorable !== null,
    restore,
    discard,
    clear,
  };
}
