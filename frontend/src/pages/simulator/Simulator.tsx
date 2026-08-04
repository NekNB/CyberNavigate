import type { FC } from "react";
import Footer from "../../components/Footer/Footer";
import Header from "../../components/Header/Header";
import SimulatorMainPage from "./SimulatorMainPage/SimulatorMainPage";

const Simulator: FC = () => {
  return (
    <>
      <Header />
      <SimulatorMainPage />
      <Footer />
    </>
  );
};

export default Simulator;
