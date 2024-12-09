import { useFetch } from "@/composables/useFetch";
import type { EncoderResponse } from "~/types/encoders";

export const getEncoders = async () => useFetch<EncoderResponse>("/encoders");
