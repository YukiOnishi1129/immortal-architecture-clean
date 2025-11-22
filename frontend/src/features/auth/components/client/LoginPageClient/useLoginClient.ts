"use client";

import { signIn } from "next-auth/react";
import { useCallback } from "react";

export function useLoginClient() {
  const handleGoogleLogin = useCallback(() => {
    signIn("google", {
      callbackUrl: "/notes",
    });
  }, []);

  return {
    handleGoogleLogin,
  };
}
