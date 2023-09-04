import type { FunctionalComponent } from "preact";
import AsyncRoute from "preact-async-route";
import PreactRouter from "preact-router";
import Home from "../routes/home";

const Router: FunctionalComponent = () => {
  return (
    <PreactRouter>
      <Home path="/" />
      <AsyncRoute path="/games/:game" getComponent={() => import("../routes/game").then(module => module.default)} />
    </PreactRouter>
  );
};

export default Router;
