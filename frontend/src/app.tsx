import type { FunctionalComponent } from "preact";
import { SWRConfig } from "swr";
import { Footer } from "./components/Footer";
import Header from "./components/Header";
import Router from "./components/Router";

const Layout: FunctionalComponent = () => {
  return (
    <div className="min-h-screen flex flex-col">
      <Header />
      <main class="min-y-8 flex-1 flex flex-col">
        <Router />
      </main>
      <Footer />
    </div>
  );
};

const fetcher = (url: RequestInfo | URL, init?: RequestInit) => fetch(`/api${url}`, init).then(res => res.json());

export const App: FunctionalComponent = () => {
  return (
    <SWRConfig value={{ fetcher }}>
      <Layout />
    </SWRConfig>
  );
};
