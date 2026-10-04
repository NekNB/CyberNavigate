import { Route, Routes } from "react-router";
import Article from "../pages/article/Article";
import Home from "../pages/home/Home";
import Scenario from "../pages/scenario/Scenario";
import Simulator from "../pages/simulator/Simulator";

function App() {
  return (
    <Routes>
      <Route path="/" element={<Home />} />
      <Route path="/articles/*" element={<Article />} />
      <Route path="/simulator" element={<Simulator />} />
      <Route path="/scenario/:id" element={<Scenario />} />
      <Route path="*" element={"Not Found"} />
    </Routes>
  );
}

export default App;
