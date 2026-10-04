import Footer from "@/app/(main)/_components/Footer/Footer";
import Header from "@/app/(main)/_components/Header/Header";

const MainLayout = ({ children }: { children: React.ReactNode }) => {
  return (
    <>
      <Header />
      {children}
      <Footer />
    </>
  );
};

export default MainLayout;
