import PreactRouter from "preact-router";
import type {FunctionalComponent} from "preact";
import Home from "../routes/home";

const Router: FunctionalComponent = () => {
    return <PreactRouter>
        <Home path="/" />
    </PreactRouter>;
};

export default Router;
