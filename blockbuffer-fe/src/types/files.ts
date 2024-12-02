export interface File {
  id: string;
  filePath: string;
  status: string;
  progress: number;
}

export enum MessageTypes {
  CREATE_FILE = 'create_file',
  UPDATE_FILE = 'update_file',
  DELETE_FILE = 'delete_file',
  REFRESH_FILES = 'refresh_files',
}

export interface FileMessage {
  type: MessageTypes;
  data: File[];
}
