import {FunctionalComponent} from "preact";

const Header: FunctionalComponent = () => {
    return <header className="py-8">
        <div className="container px-4">
            <h1 className="font-extrabold text-3xl flex items-center mb-4">🏆 Trophies.gg</h1>
        </div>
    </header>;
};

export default Header;
