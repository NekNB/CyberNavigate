import {
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState,
  type FC,
} from "react";
import { useParams } from "react-router";

import {
  CreateSession,
  GetResults,
  GetStep,
  SendAnswer,
} from "../../app/api/Simulator/Simulator";
import type {
  IChat,
  IChatAnswer,
  IChatFile,
  IChatMessage,
} from "../../types/messenger";
import {
  type IAction,
  type IMessage,
  type IResults,
  type ISMS,
  type IStep,
} from "../../types/simulator";
import Final from "./Final/Final";
import Loader, { type ILoaderProps } from "./Loader/Loader";
import Messenger from "./Messenger/Messenger";

const delay = () => new Promise((resolve) => setTimeout(resolve, 2_000));

const withDelay = async (f: Promise<void>) => {
  await f;
  await delay();
};

const Scenario: FC = () => {
  const { id: scenarioId } = useParams<{ id: string }>();
  const [isLoading, setIsLoading] = useState(true);
  const [step, setStep] = useState<IStep>();
  const [chats, setChats] = useState<Map<string, IChat>>(new Map());
  const [stepCompleted, setStepCompleted] = useState(false);
  const [isFrozen, setIsFrozen] = useState(false);
  const [isGameFinished, setIsGameFinished] = useState(false);
  const [isGetResults, setIsGetResults] = useState(false);
  const [results, setResults] = useState<IResults>({
    errors: [],
    gameDuration: 0,
    trustGraph: [],
  });

  const [sms, setSms] = useState<ISMS[] | undefined>();

  const unReadCount = useMemo((): number => {
    let count = 0;
    chats.forEach((chat) => {
      count += chat.unRead || 0;
    });
    return count;
  }, [chats]);
  const prevUnReadCount = useRef(unReadCount);

  useEffect(() => {
    const handleBeforeUnload = (event: BeforeUnloadEvent) => {
      event.preventDefault();
    };
    window.addEventListener("beforeunload", handleBeforeUnload);
    return () => {
      window.removeEventListener("beforeunload", handleBeforeUnload);
    };
  }, []);

  useEffect(() => {
    const getResults = async () => {
      setResults(await GetResults());
      setIsGameFinished(true);
    };
    if (isGetResults) {
      getResults();
    }
    console.log(isGetResults);
  }, [isGetResults]);

  const getStep = useCallback(async () => {
    const response = await GetStep();
    switch (response.code) {
      case 200:
        setStep(response.step as IStep);
        console.log("Получен ответ 200");
        break;
      case 204:
        setStep(undefined);
        setResults(await GetResults());
        setIsGameFinished(true);
        console.log("Получен ответ 204");
        break;
    }
    setStepCompleted(false);
  }, []);

  useEffect(() => {
    const doAction = (action: IAction) => {
      switch (action.type) {
        case "sms":
          setSms((prevSms) => {
            const newSms = action.message as ISMS;
            if (prevSms) {
              return [newSms, ...prevSms];
            } else {
              return [newSms];
            }
          });
          break;
        case "message":
          const msg = action.message as IMessage;

          const newMsg: IChatMessage = {
            messageId: msg.id,
            isInput: true,
            answers: msg.answers?.map((answer): IChatAnswer => {
              return { answerId: answer.id, text: answer.text };
            }),
            files: msg.files?.map((file): IChatFile => {
              return {
                ...file,
                fileId: file.id,
              };
            }),
            text: msg.text,
          };
          setChats((prevChats) => {
            const chat = prevChats?.get(msg.senderId);
            const newChats = new Map<string, IChat>(prevChats);

            if (chat) {
              const updatedChat: IChat = {
                ...chat,
                messages: [...chat.messages, newMsg],
                unRead: (chat.unRead || 0) + 1,
              };
              newChats.set(updatedChat.senderId, updatedChat);
            } else {
              newChats.set(msg.senderId, {
                messages: [newMsg],
                senderId: msg.senderId,
                senderName: msg.senderName,
                unRead: 1,
              });
            }

            return newChats;
          });
      }
    };

    if (step) {
      console.log("Step изменился");
      const actions = step.actions;
      const timeoutsIds = actions.map((action) =>
        setTimeout(() => doAction(action), action.delay * 1000),
      );

      // Вычисляем максимальный delay и ставим задержку
      const slowestAction = actions.reduce((max, current) => {
        return current.delay > max.delay ? current : max;
      }, actions[0]);
      timeoutsIds.push(
        setTimeout(
          () => setStepCompleted(true),
          slowestAction.delay * 1000 + 10,
        ),
      );

      return () => {
        timeoutsIds.forEach(clearTimeout);
      };
    } else {
      console.log("Step init");
    }
  }, [step]);

  useEffect(() => {
    console.log(prevUnReadCount.current, unReadCount, stepCompleted);
    if (
      (prevUnReadCount.current > 0 || true) &&
      stepCompleted &&
      unReadCount === 0
    ) {
      console.log("Все сообщения прочитаны, запрашиваем следующий шаг");
      getStep();
    }

    // Обновляем реф для следующего рендера
    prevUnReadCount.current = unReadCount;
  }, [unReadCount, stepCompleted]);

  const loaderTasks: ILoaderProps["actions"] = useMemo(
    () => [
      {
        action: () => withDelay(CreateSession(scenarioId as string)),
        placeholder: "Создаем сессию",
      },
      {
        action: () => withDelay(getStep()),
        placeholder: "Готовимся запустить игру",
      },
    ],
    [scenarioId, getStep],
  );

  const handleSetIsLoading = useCallback(async () => {
    await delay();
    setIsLoading(false);
  }, []);

  const readMessage = useCallback(
    (chatId: string) => {
      setChats((prevChats) => {
        const chat = prevChats.get(chatId)!;
        let newUnread = chat.unRead;

        const lastMessage = chat.messages.at(-1);
        const hasAnswers =
          lastMessage?.answers && lastMessage.answers.length > 0;
        newUnread = hasAnswers ? 1 : 0;

        // Если чата нет или он уже полностью прочитан (unRead === 0), ничего не меняем
        if (!chat || chat.unRead === newUnread) {
          return prevChats; // Возвращаем старый объект -> ре-рендера не будет
        }

        return new Map<string, IChat>([
          ...prevChats,
          [chatId, { ...chat, unRead: newUnread }],
        ]);
      });
    },
    [chats],
  );

  const sendAnswer = useCallback(
    async (senderId: string, answer: IChatAnswer) => {
      await SendAnswer(answer.answerId);

      const chat = chats.get(senderId)!;

      const newMessage: IChatMessage = {
        messageId: crypto.randomUUID(),
        isInput: false,
        text: answer.text,
      };

      const updatedChat: IChat = {
        ...chat,
        messages: [...chat.messages, newMessage],
        unRead: 0,
      };

      setChats((prevChats) => {
        return new Map<string, IChat>([...prevChats, [senderId, updatedChat]]);
      });
    },
    [chats],
  );

  return isLoading ? (
    <Loader actions={loaderTasks} setIsLoading={handleSetIsLoading} />
  ) : (
    <>
      <Messenger
        sms={sms}
        setIsGetResults={(isGetResults) => setIsGetResults(isGetResults)}
        isFrozen={isFrozen}
        sendAnswer={sendAnswer}
        chats={chats as Map<string, IChat>}
        setIsFinished={setIsGameFinished}
        readMessage={readMessage}
      />

      {isGameFinished && (
        <Final
          results={results}

          onClose={() => {
            setIsGameFinished(false);
            setIsFrozen(true);
          }}
        />
      )}
    </>
  );
};

export default Scenario;
