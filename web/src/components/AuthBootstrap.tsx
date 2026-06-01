"use client";

import { useEffect, useRef } from "react";
import { useRouter } from "next/navigation";
import { refreshRequest } from "@/api/auth";
import {
  isInitialized,
  markInitialized,
  setAccessToken,
} from "@/lib/access-token";

const AUTH_EXPIRED_EVENT = "musclead:auth-expired";

export function AuthBootstrap() {
  const router = useRouter();
  const bootstrapped = useRef(false);

  useEffect(() => {
    let cancelled = false;
    if (!bootstrapped.current) {
      bootstrapped.current = true;
      (async () => {
        if (isInitialized()) return;
        const tokens = await refreshRequest();
        if (cancelled) return;
        if (tokens) {
          setAccessToken(tokens.access_token);
        } else {
          markInitialized();
        }
      })();
    }
    const onExpired = () => {
      router.replace("/login");
    };
    window.addEventListener(AUTH_EXPIRED_EVENT, onExpired);
    return () => {
      cancelled = true;
      window.removeEventListener(AUTH_EXPIRED_EVENT, onExpired);
    };
  }, [router]);

  return null;
}
