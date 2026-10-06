// app/layout.tsx
import { ClientConfigProvider } from "@/components/providers/ClientConfigProvider/ClientConfigProvider";
import ScrollResetProvider from "@/components/providers/ScrollResetProvier/ScrollResetProvider";
import { loadConfig } from "@/lib/config";
import type { Metadata, Viewport } from "next";
import { Oswald } from "next/font/google";
import "./global.css"; // ← глобальные стили подключаются ЗДЕСЬ
// import { Providers } from "./providers";

// устанавливаем шрифт Google Oswald
const oswald = Oswald({
  weight: ["200", "300", "400", "500", "600", "700"],
  subsets: ["latin", "cyrillic"],
  variable: "--font-oswald", // доступен в CSS как var(--font-oswald)
  display: "swap",
});

// TODO: edit title and description
export const metadata: Metadata = {
  title: {
    default: "КиберНавигатор | Обучение кибербезопасности для молодежи",
    template: "%s | КиберНавигатор  ",
  },
  description:
    "КиберНавигатор — платформа повышения киберграмотности молодежи и студентов РФ. Изучайте статьи по безопасности и тренируйтесь в симуляторе мошенника.",
  keywords: [
    "кибербезопасность",
    "киберграмотность",
    "симулятор мошенника",
    "защита персональных данных",
    "молодежь",
    "студенты",
  ],
  verification: {
    google: "cCIWxjjJPbOikaMFZujIwog-MUdnMhcdHG4pVtZ5480"
  },
  openGraph: {
    title: "КиберНавигатор — защита от онлайн-угроз",
    description:
      "Интерактивная платформа для молодежи. Читай статьи и учись распознавать уловки социальной инженерии в симуляторе мессенджера.",
    locale: "ru_RU",
    type: "website",
  },
};

export const viewport: Viewport = {
  themeColor: "#000000",
  width: "device-width",
  initialScale: 1,
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const { externalUrl } = loadConfig();

  return (
    <html lang="ru" className={oswald.variable}>
      <body className={oswald.className}>
        <ClientConfigProvider clientApiUrl={externalUrl}>
          {children}
        </ClientConfigProvider>
        <ScrollResetProvider />
      </body>
    </html>
  );
}
