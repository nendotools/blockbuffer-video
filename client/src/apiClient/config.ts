import { useFetch } from "@/composables/useFetch";
import type { Config } from "~/types/config";

export const getConfig = async () => useFetch<Config>("/config");
export const updateConfig = async (config: Record<string, any>) =>
  useFetch("/config", { method: "POST", body: { ...config } });
