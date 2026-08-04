import { useEffect, useState, type FC } from "react";
import type FileDownloader from "../../../../types/fileDownloader";
import type { IChat, IChatAnswer } from "../../../../types/messenger";

import styles from "./Chat.module.css";
import downloadImg from "/assets/download.svg";
import fileImg from "/assets/file.svg";
interface ChatProps {
  chat: IChat;
  fileDownloader: FileDownloader;
  sendAnswer: (senderId: string, answer: IChatAnswer) => void;
  isFrozen: boolean;
  setIsFinished: React.Dispatch<React.SetStateAction<boolean>>;
}

const Chat: FC<ChatProps> = ({
  chat,
  fileDownloader,
  isFrozen,
  sendAnswer,
  setIsFinished,
}) => {
  const handleAnswerOnClick = async (answer: IChatAnswer) => {
    sendAnswer(chat.senderId, answer);
  };

  const answers = chat.messages.at(-1)?.answers;
  return (
    <div className={`${isFrozen ? styles.frozen : ""} ${styles.chatWindow}`}>
      <div className={styles.chatName}>{chat.senderName}</div>

      <div className={styles.chat}>
        {chat.messages.map((message) => {
          return (
            <Message
              key={message.messageId}
              isInput={message.isInput}
              text={message.text}
              files={message.files}
              fileDownloader={fileDownloader}
              setIsFinished={setIsFinished}
            />
          );
        })}
      </div>

      {(answers?.length as number) > 0 && (
        <div className={styles.answers}>
          {chat.messages.at(-1)?.answers?.map((answer) => {
            return (
              <div
                key={answer.answerId}
                className={styles.answer}
                onClick={() => handleAnswerOnClick(answer)}
              >
                {answer.text}
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
};

export default Chat;

interface MessageProps {
  isInput: boolean;
  text?: string;
  files?: { fileId: string; filename: string; size: number }[];
  fileDownloader: FileDownloader;
  setIsFinished: React.Dispatch<React.SetStateAction<boolean>>;
}
const Message: FC<MessageProps> = ({
  isInput,
  text,
  files,
  fileDownloader,
  setIsFinished,
}) => {
  return (
    <div
      className={`${styles.message} ${isInput ? styles.inputMessage : styles.outputMessage}`}
    >
      <p className={styles.messageText}>{text}</p>
      {files?.length
        ? files.map((file) => {
            return (
              <File
                key={`${file.filename}`}
                fileId={file.fileId}
                filename={file.filename}
                size={file.size}
                fileDownloader={fileDownloader}
                setIsFinished={setIsFinished}
              />
            );
          })
        : null}
    </div>
  );
};

interface FileProps {
  fileId: string;
  filename: string;
  size: number;
  fileDownloader: FileDownloader;
  setIsFinished: React.Dispatch<React.SetStateAction<boolean>>;
}
// Константы для SVG круга
const RADIUS = 24; // Радиус круга (подгоните под размер вашей картинки)
const CIRCUMFERENCE = 2 * Math.PI * RADIUS; // Длина окружности

const File: FC<FileProps> = ({
  fileId,
  filename,
  size,
  fileDownloader,
  setIsFinished,
}) => {
  const [isDownloaded, setIsDownloaded] = useState(
    fileDownloader.isDownloaded(fileId),
  );
  const [isLoading, setIsLoading] = useState(false);
  const [isGameFall, setIsGameFall] = useState(false);

  useEffect(() => {
    if (isGameFall && isDownloaded && !isLoading) {
      console.log("Игра окончена");
      setIsFinished(true);
    }
  }, [isGameFall, isDownloaded, isLoading]);

  const handleDownloadClick = () => {
    if (!isDownloaded && !isLoading) {
      setIsLoading(true);
      fileDownloader
        .downloadFile(fileId)
        .then((isSafe) => {
          setIsGameFall(!isSafe);
        })
        .catch((error) => {
          console.error("Ошибка при скачивании файла:", error);
        });

      requestAnimationFrame(() => {
        setIsDownloaded(true);
      });

      const durationMs = animationDuration * 1000;
      setTimeout(() => {
        setIsLoading(false);
      }, durationMs);
    }
  };

  const animationDuration = size / 100;

  return (
    <div className={styles.file}>
      <div className={styles.fileImgWrapper}>
        <img
          className={styles.fileImg}
          src={!isLoading && isDownloaded ? fileImg : downloadImg}
          onClick={handleDownloadClick}
        />

        {/* SVG для рисования круга */}
        {isLoading && (
          <svg
            className={styles.progressRing}
            width="100%"
            height="100%"
            viewBox="0 0 56 56"
          >
            <circle
              className={styles.progressCircle}
              cx="28"
              cy="28"
              r={RADIUS}
              fill="transparent"
              strokeWidth="3" // Толщина границы
              strokeDasharray={CIRCUMFERENCE}
              strokeDashoffset={isDownloaded ? 0 : CIRCUMFERENCE}
              style={{
                transition: `stroke-dashoffset ${animationDuration}s linear`,
              }}
            />
          </svg>
        )}
      </div>

      <p className={styles.filename}>{filename}</p>
      <p className={styles.fileSize}>{size} MB</p>
    </div>
  );
};
