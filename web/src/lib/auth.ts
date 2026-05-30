"use client";

import { useSyncExternalStore } from "react";

export const USER_ID_STORAGE_KEY = "musclead.userId";
const EVENT_NAME = "musclead:userId";

export function getStoredUserId(): string | null {
  if (typeof window === "undefined") return null;
  return window.localStorage.getItem(USER_ID_STORAGE_KEY);
}

export function setStoredUserId(userId: string): void {
  window.localStorage.setItem(USER_ID_STORAGE_KEY, userId);
  window.dispatchEvent(new Event(EVENT_NAME));
}

export function clearStoredUserId(): void {
  window.localStorage.removeItem(USER_ID_STORAGE_KEY);
  window.dispatchEvent(new Event(EVENT_NAME));
}

function subscribe(onChange: () => void): () => void {
  window.addEventListener(EVENT_NAME, onChange);
  window.addEventListener("storage", onChange);
  return () => {
    window.removeEventListener(EVENT_NAME, onChange);
    window.removeEventListener("storage", onChange);
  };
}

export function useUserId(): { userId: string | null; ready: boolean } {
  const userId = useSyncExternalStore(
    subscribe,
    getStoredUserId,
    () => null, // SSR snapshot
  );
  const ready = useSyncExternalStore(
    subscribe,
    () => true,
    () => false, // false on SSR, true after hydration
  );
  return { userId, ready };
}
