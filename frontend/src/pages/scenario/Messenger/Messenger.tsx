import { useEffect, useState, type FC } from "react";
import FileDownloader from "../../../types/fileDownloader";
import type { IChat, IChatAnswer } from "../../../types/messenger";

import type { ISMS } from "../../../types/simulator";
import Chat from "./Chat/Chat";
import styles from "./Messenger.module.css";
import SMS from "./Sms/Sms";

interface MessengerProps {
  chats: Map<string, IChat>;
  setIsFinished: React.Dispatch<React.SetStateAction<boolean>>;
  sendAnswer: (senderId: string, answer: IChatAnswer) => Promise<void>;
  readMessage: (chatId: string) => void;
  isFrozen: boolean;
  setIsGetResults: (isGetResults: boolean) => void;
  sms?: ISMS[];
}

const Messenger: FC<MessengerProps> = ({
  chats,
  setIsFinished,
  setIsGetResults,
  readMessage,
  sendAnswer,
  isFrozen,
  sms,
}) => {
  const [openedChat, setOpenedChat] = useState<string>();
  const fileDownloader = new FileDownloader();
  const [openedSms, setOpenedSms] = useState(false);

  const handleChoiceChat = (chatId: string) => {
    setOpenedChat(chatId);
  };

  useEffect(() => {
    if (sms) {
      setOpenedSms(true);
    }
  }, [sms]);

  useEffect(() => {
    if (openedChat) {
      readMessage(openedChat);
    }
  }, [chats, openedChat, readMessage]);
  return (
    <div className={styles.messenger}>
      {sms && openedSms && (
        <SMS sms={sms} closeSms={() => setOpenedSms(false)} />
      )}
      <header className={styles.header}>
        <img
          src="/assets/messenger.webp"
          className={styles.logo}
          onClick={() => setOpenedSms(true)}
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
            setIsGetResults={setIsGetResults}
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
