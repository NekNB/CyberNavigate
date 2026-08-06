import type { FC } from "react";
import type { ISMS } from "../../../../types/simulator";
import styles from "./Sms.module.css";
export interface SMSProps {
  sms: ISMS[];
  closeSms: () => void;
}

const SMS: FC<SMSProps> = ({ sms, closeSms }) => {
  return (
    <div className={styles.overlay} onClick={closeSms}>
      <div className={styles.smsWrapper} onClick={(e) => e.stopPropagation()}>
        {sms.map((sms) => {
          return (
            <div key={sms.id} className={styles.sms}>
              <h4 className={styles.senderName}>{sms.senderName}</h4>
              <p className={styles.text}> {sms.text}</p>
            </div>
          );
        })}
      </div>
    </div>
  );
};

export default SMS;
