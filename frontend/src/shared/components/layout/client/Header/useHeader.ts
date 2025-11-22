"use client";
import { useQueryClient } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { signOut, useSession } from "next-auth/react";
import { useCallback } from "react";

export function useHeader() {
  const { data: session } = useSession();
  const router = useRouter();
  const queryClient = useQueryClient();

  const handleSignOut = useCallback(async () => {
    await signOut({ redirect: false });
    // キャッシュをすべてクリア
    queryClient.clear();
    router.push("/login");
  }, [router, queryClient]);

  return {
    userName: session?.user?.name || undefined,
    userEmail: session?.user?.email || undefined,
    userImage: session?.user?.image || undefined,
    handleSignOut,
  };
}
