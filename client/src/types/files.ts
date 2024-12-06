export interface File {
  id: string;
  filePath: string;
  status: FileStatuses;
  progress: number;
  duration: number; // in seconds
}

export enum MessageTypes {
  CREATE_FILE = 'create_file',
  UPDATE_FILE = 'update_file',
  DELETE_FILE = 'delete_file',
  REFRESH_FILES = 'refresh_files',
}

export enum FileStatuses {
  NEW = 'new',
  QUEUED = 'queued',
  PROCESSING = 'processing',
  COMPLETED = 'completed',
  COMPLETEDELETED = 'completed-deleted',
  CANCELLED = 'cancelled',
  REJECTED = 'rejected',
  FAILED = 'failed',
  DELETED = 'deleted'
}

export interface FileMessage {
  type: MessageTypes;
  data: File[];
}
