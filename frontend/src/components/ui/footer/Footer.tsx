import "./Footer.css";
import type { FC } from "react";

const Footer: FC = () => {
  return (
    <footer className="footer">
      <div className="footer-title">Наши контакты</div>
      <div className="footer-email">email@example.com</div>
      <div className="footer-phone">+7-918-989-88-99</div>
    </footer>
  );
};

export default Footer;