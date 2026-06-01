"use client";

import { useSyncExternalStore } from "react";

let accessToken: string | null = null;
let initialized = false;
const listeners = new Set<() => void>();

function emit(): void {
  listeners.forEach((l) => l());
}

export function getAccessToken(): string | null {
  return accessToken;
}

export function setAccessToken(token: string): void {
  accessToken = token;
  initialized = true;
  emit();
}

export function clearAccessToken(): void {
  accessToken = null;
  initialized = true;
  emit();
}

export function markInitialized(): void {
  initialized = true;
  emit();
}

export function isInitialized(): boolean {
  return initialized;
}

function subscribe(onChange: () => void): () => void {
  listeners.add(onChange);
  return () => {
    listeners.delete(onChange);
  };
}

export function useAccessToken(): { token: string | null; ready: boolean } {
  const token = useSyncExternalStore(
    subscribe,
    () => accessToken,
    () => null,
  );
  const ready = useSyncExternalStore(
    subscribe,
    () => initialized,
    () => false,
  );
  return { token, ready };
}
