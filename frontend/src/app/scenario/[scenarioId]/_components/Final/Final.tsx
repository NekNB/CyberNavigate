import { useMemo, type FC } from "react";

import type { IResults } from "@/types/simulator";
import { useRouter } from "next/navigation";
import styles from "./Final.module.css";

interface FinalProps {
  onClose: () => void;
  results: IResults;
}

const Final: FC<FinalProps> = ({ results }) => {
  const router = useRouter();

  const gameDuration = useMemo(() => {
    if (results.gameDuration === 0 || results.gameDuration === 1) {
      return `1 минуту`;
    }

    const durationString = results.gameDuration.toString();
    if ((durationString.at(-1) as string) in ["2", "3", "4"]) {
      return `${durationString} минуты`;
    } else {
      return `${durationString} минут`;
    }
  }, [results]);

  return (
    <div className={styles.wrapper}>
      <div className={styles.finalWindow} onClick={(e) => e.stopPropagation()}>
        <h3 className={styles.title}>Игра окончена</h3>

        <p className={styles.duration}>⏱ Игра длилась: {gameDuration}</p>

        {results.errors && results.errors.length > 0 && (
          <div className={styles.errors}>
            <h3>Ошибки:</h3>
            <ul>
              {results.errors.map((error, index) => (
                <li key={index}>{error}</li>
              ))}
            </ul>
          </div>
        )}
        <button
          className={styles.backToMenu}
          onClick={() => router.push("/simulator")}
        >
          Вернуться в меню
        </button>
      </div>
    </div>
  );
};

export default Final;
