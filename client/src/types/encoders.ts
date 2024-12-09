// API Response Types
export interface AVOptionEnum {
  id: string;
  option: string;
  description: string;
}

export interface AVOption {
  name: string;
  description: string;
  type: string;
  options: AVOptionEnum[];
}

export interface Encoder {
  type: string;
  name: string;
  description: string;
  formats: string[];
  sampleRates: number[];
  options: AVOption[];
}

// API Request Types (chosen settings)
export interface AVOptionProfile {
  name: string;
  value: string;
}

export interface AudioEncoderProfile {
  name: string;
  sampleRate: number;
  options: AVOptionProfile[];
}

export interface EncoderProfile {
  name: string;
  format: string;
  options: AVOptionProfile[];
  audioEncoder: AudioEncoderProfile;
}


export interface EncoderResponse {
  defaultEncoder: EncoderProfile;
  videoEncoders: Encoder[];
  audioEncoders: Encoder[];
}
