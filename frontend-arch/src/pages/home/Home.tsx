import type { FC } from "react";
import Footer from "../../components/Footer/Footer";
import Header from "../../components/Header/Header";

import Feature from "./sections/Feature/Feature";
import Hero from "./sections/Hero/Hero";
import Stats from "./sections/Stats/Stats";

const Home: FC = () => {
  return (
    <>
      <Header />
      <Hero />
      <Stats />
      <Feature />
      <Footer />
    </>
  );
};

export default Home;
