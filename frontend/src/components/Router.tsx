import type { FunctionalComponent } from "preact";
import AsyncRoute from "preact-async-route";
import PreactRouter from "preact-router";
import Home from "../routes/home";

const Router: FunctionalComponent = () => {
  return (
    <PreactRouter>
      <Home path="/" />
      <AsyncRoute path="/games/:id" getComponent={() => import("../routes/game").then(module => module.default)} />
      <NotFound default />
    </PreactRouter>
  );
};

const NotFound: FunctionalComponent = () => {
  return (
    <div className="flex-1 flex flex-col justify-center container px-4">
      <p>404 error</p>
      <h2 className="text-2xl font-bold mt-2 mb-4 md:text-4xl">We couldn&#39;t find that page</h2>
      <p className="text-gray-500 dark:text-gray-400">
        Sorry, the page you&#39;re looking for doesn&#39;t exist or has been moved.
      </p>
      <div className="flex items-center gap-8 mt-4">
        <button
          onClick={() => history.back()}
          className="text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200 border border-gray-500 hover:border-gray-700 dark:border-gray-400 hover:dark:border-gray-200 px-4 py-2  transition-colors duration-300 ease-in-out"
        >
          Go back
        </button>
        <a
          href="/"
          className="bg-indigo-500 hover:bg-indigo-600 text-white px-6 py-2 font-semibold transition-colors duration-300 ease-in-out"
        >
          Take me home
        </a>
      </div>
    </div>
  );
};

export default Router;
