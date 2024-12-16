export type AVOption = Record<string, string>

export interface VideoPreset {
  codec: string;
  format: string;
  options: AVOption[];
}

export interface AudioPreset {
  codec: string;
  sampleRate: number;
  options: AVOption[];
}

export interface Preset {
  name: string;
  description: string;
  extension: string;
  video: VideoPreset;
  audio: AudioPreset;
}

export interface PresetsResponse {
  presets: Preset[];
}
