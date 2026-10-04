"use client";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { FC, useState } from "react";
import styles from "./MobileSideBar.module.css";

interface MobileSideBarProps {
  elements: { id: string; title: string }[];
  drawerTitle: string;
  openButtonTitle: string;
  baseUrl: string;
}

const MobileSideBar: FC<MobileSideBarProps> = ({
  elements,
  drawerTitle,
  openButtonTitle,
  baseUrl,
}) => {
  const path = usePathname();

  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const closeMenu = () => {
    setIsMenuOpen(false);
  };

  return (
    <>
      <button
        className={styles.mobileMenuTrigger}
        onClick={() => setIsMenuOpen(true)}
      >
        ☰ {openButtonTitle}
      </button>
      <aside
        className={`${styles.elementsMenu} ${isMenuOpen ? styles.drawerActive : ""}`}
      >
        <div className={styles.drawerHeader}>
          <span className={styles.drawerTitle}>{drawerTitle}</span>
          <button
            className={styles.closeBtn}
            onClick={() => setIsMenuOpen(false)}
            aria-label="Закрыть"
          >
            ✕
          </button>
        </div>

        <div className={styles.mobileArticlesList}>
          {elements.map((element) => (
            <Link
              key={element.id}
              className={`${styles.mobileArticleItem} ${
                path.includes(element.id) ? styles.mobileArticleActive : ""
              }`}
              onClick={() => setIsMenuOpen(false)}
              href={`${baseUrl}/${element.id}`}
            >
              {element.title}
            </Link>
          ))}
        </div>
      </aside>

      {isMenuOpen && (
        <div className={styles.overlay} onClick={() => setIsMenuOpen(false)} />
      )}
    </>
  );
};
export default MobileSideBar;
