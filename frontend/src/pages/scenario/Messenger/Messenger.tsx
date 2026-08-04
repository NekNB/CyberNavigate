import { useEffect, useState, type FC } from "react";
import FileDownloader from "../../../types/fileDownloader";
import type { IChat, IChatAnswer } from "../../../types/messenger";

import Chat from "./Chat/Chat";
import styles from "./Messenger.module.css";

interface MessengerProps {
  chats: Map<string, IChat>;
  setIsFinished: React.Dispatch<React.SetStateAction<boolean>>;
  sendAnswer: (senderId: string, answer: IChatAnswer) => Promise<void>;
  readMessage: (chatId: string) => void;
  isFrozen: boolean;
}

const Messenger: FC<MessengerProps> = ({
  chats,
  setIsFinished,
  readMessage,
  sendAnswer,
  isFrozen,
}) => {
  const [openedChat, setOpenedChat] = useState<string>();
  const fileDownloader = new FileDownloader();

  const handleChoiceChat = (chatId: string) => {
    setOpenedChat(chatId);
  };

  useEffect(() => {
    if (openedChat) {
      readMessage(openedChat);
    }
  }, [chats, openedChat, readMessage]);
  return (
    <div className={styles.messenger}>
      <header className={styles.header}>
        <img
          src="/assets/messenger.png"
          className={styles.logo}
          onClick={isFrozen ? () => setIsFinished(true) : () => {}}
        />
      </header>
      <nav className={styles.chatList}>
        {[...chats].map(([chatId, chatData]) => {
          return (
            <div
              key={chatId}
              onClick={() => handleChoiceChat(chatId)}
              className={styles.chatCard}
            >
              <h3 className={styles.chatName}>{chatData.senderName}</h3>
              {(chatData.unRead as number) > 0 && (
                <div className={styles.unReadCount}>{chatData.unRead}</div>
              )}

              {chatData.messages.at(-1)?.text && (
                <p className={styles.lastMessage}>
                  {chatData.messages.at(-1)?.text}
                </p>
              )}
            </div>
          );
        })}
      </nav>
      <main className={styles.chat}>
        {openedChat && (
          <Chat
            isFrozen={isFrozen}
            setIsFinished={setIsFinished}
            chat={chats.get(openedChat)!}
            fileDownloader={fileDownloader}
            sendAnswer={sendAnswer}
          />
        )}
      </main>
    </div>
  );
};

export default Messenger;
