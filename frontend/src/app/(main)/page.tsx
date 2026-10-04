import type { FC } from "react";

import Hero from "@/app/(main)/_components/Hero/Hero";
import Stats from "@/app/(main)/_components/Stats/Stats";
import Feature from "./_components/Feature/Feature";

const Home: FC = () => {
  return (
    <>
      <Hero />
      <Stats />
      <Feature />
    </>
  );
};

export default Home;
