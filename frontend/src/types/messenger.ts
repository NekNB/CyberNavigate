export interface IChat {
  senderId: string;
  senderName: string;
  unRead: number;
  messages: IChatMessage[];
}

export interface IChatMessage {
  messageId: string;
  text?: string;
  isInput: boolean;
  files?: IChatFile[];
  answers?: IChatAnswer[];
}

export interface IChatFile {
  fileId: string;
  filename: string;
  isSafe: boolean;
  error?: string;
  size: number;
}

export interface IChatAnswer {
  answerId: string;
  text: string;
}
