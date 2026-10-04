"use client";

import { IScenario } from "@/types/simulator";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { FC } from "react";
import styles from "./SideBar.module.css";
interface Props {
  scenariosList: IScenario[];
}

const SideBar: FC<Props> = ({ scenariosList: scenariosList }) => {
  const path = usePathname();
  return (
    <div className={styles.scenariosList}>
      {scenariosList.map((item) => {
        const id = item.id;
        const title = item.title;
        const isActive = path.includes(id);

        return (
          <Link
            key={id}
            className={`${styles.scenarioTitle} ${isActive ? styles.active : ""}`}
            href={`/simulator/${item.id}`}
          >
            {title}
          </Link>
        );
      })}
    </div>
  );
};

export default SideBar;
