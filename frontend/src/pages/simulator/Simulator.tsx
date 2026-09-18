import type { FC } from "react";
import Footer from "../../components/Footer/Footer";
import SimulatorMainPage from "./SimulatorMainPage/SimulatorMainPage";

const Simulator: FC = () => {
  return (
    <>
      <SimulatorMainPage />
      <Footer />
    </>
  );
};

export default Simulator;
