import useSWR from "swr";
import type { FunctionalComponent} from "preact";
import Link from "./Link";

type ApiResponse = {
    version?: string;
    date: string;
    commit?: string;
};

const BuildInfo: FunctionalComponent = () => {
    const { data, error } = useSWR<ApiResponse>("/version");

    if (error || !data) return null;

    return <div className="text-gray-700 dark:text-slate-200 text-sm text-right font-mono">
        {data.version && <div className="space-x-2">
            <Link href={`https://github.com/lukecarr/trophies/releases/tag/${data.version}`} newTab>
                {data.version}
            </Link>
            {data.commit && <Link href={`https://github.com/lukecarr/trophies/commit/${data.commit}`} newTab>
                ({data.commit.slice(0, 7)})
            </Link>}
        </div>}
        <p className="text-xs">Built at: {data.date}</p>
    </div>
}

export const Footer: FunctionalComponent = () => {
    return <footer className="py-8">
        <div className="container px-4 flex justify-between items-center">
            <p><Link href={`https://github.com/lukecarr/trophies`} newTab>Trophies.gg</Link> is made with ❤️ by <Link href={`https://github.com/lukecarr`} newTab>Luke Carr</Link></p>
            <BuildInfo />
        </div>
    </footer>
}