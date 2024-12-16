import { useFetch } from "@/composables/useFetch";
import type { PresetsResponse } from "~/types/presets";

export const getPresets = async () => useFetch<PresetsResponse>("/presets");
