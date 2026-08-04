import React from "react";
import styles from "./ProgressBar.module.css"; // Подключим стили ниже

interface ProgressBarProps {
  progress: number;
}

const ProgressBar: React.FC<ProgressBarProps> = ({ progress }) => {
  return (
    <div className={styles.progressBarContainer}>
      {/* Динамическая ширина в процентах */}
      <div
        className={styles.progressBarFill}
        style={{ width: `${progress}%` }}
      />
      <span className={styles.progressText}>{Math.round(progress)}%</span>
    </div>
  );
};

export default ProgressBar;
