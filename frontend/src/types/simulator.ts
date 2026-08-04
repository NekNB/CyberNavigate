export const EScenarioDifficulty = {
  EASY: "easy",
  MIDDLE: "middle",
  HARD: "hard",
} as const;

export type EScenarioDifficulty =
  (typeof EScenarioDifficulty)[keyof typeof EScenarioDifficulty];

export interface IScenario {
  id: string;
  difficulty: EScenarioDifficulty;
  title: string;
  description: string;
  articleIds?: string[];
}

export const EActionType = {
  MESSAGE: "message",
  SMS: "sms",
} as const;

export type EActionType = (typeof EActionType)[keyof typeof EActionType];

export interface IStep {
  id: string;
  actions: IAction[];
}

export interface IAction {
  id: string;
  type: EActionType;
  message: IMessage | ISMS;
  delay: number;
}

export interface IMessage {
  id: string;
  senderId: string;
  senderName: string;
  text?: string;
  files?: IFile[];
  answers?: IAnswer[];
}
export interface ISMS {
  id: string;
  senderName: string;
  text: string;
}

export interface IFile {
  id: string;
  filename: string;
  isSafe: boolean;
  error?: string;
  size: number;
}

export interface IAnswer {
  id: string;
  text: string;
  error?: string;
  addTrust: number;
}

export interface IFile {
  id: string;
  filename: string;
  isSafe: boolean;
  error?: string;
  size: number;
}

export interface IResults {
  gameDuration: number;
  errors?: string[];
  trustGraph: number[];
}
