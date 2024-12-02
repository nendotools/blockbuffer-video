export interface File {
  id: string;
  filePath: string;
  status: string;
  progress: number;
}

export interface FileMessage {
  type: 'update_file' | 'refresh_files';
  data: File[];
}
