import type {FunctionalComponent} from "preact";
import useSWR from "swr";
import {useMemo, useState} from "preact/compat";
import Fuse from "fuse.js";

type Game = {
    id: number;
    name: string;
    description?: string;
    iconURL?: string;
    psnID: string;
    platform: string;
};

type GameCount = {
    id: number;
    rarity: string;
    count: number;
}

const rarityBg: {
    [rarity: string]: string
} = {
    'bronze': 'bg-orange-400',
    'silver': 'bg-gray-400',
    'gold': 'bg-yellow-400',
}

const raritySort: {
    [rarity: string]: number
} = {
    'bronze': 0,
    'silver': 1,
    'gold': 2,
}

const GameCard = ({ game, counts }: { game: Game, counts: Omit<GameCount, 'id'>[] }) => {
    const name = game.name.trim();

    return <div className="bg-gray-100 dark:bg-slate-800 flex relative">
        <span className="absolute right-0 top-0 z-10 text-xs font-semibold px-2 py-1 bg-gray-200 text-gray-700 dark:bg-slate-700 dark:text-slate-400">{game.platform}</span>
        <div className="h-32 w-32 bg-cover bg-center flex" style={{backgroundImage: `url(${game.iconURL})`}}>
            <div className="backdrop-blur-lg flex items-center justify-center">
                <img className="w-full" src={game.iconURL} alt={name} />
            </div>
        </div>
        <div className="p-4 flex-1 flex flex-col justify-center space-y-4 overflow-hidden">
            <h2 title={name} className="font-bold text-xl overflow-hidden whitespace-nowrap text-ellipsis">
                {name}
            </h2>
            <div className="grid gap-6 lg:gap-3" style={{gridTemplateColumns: `repeat(${new Set(counts.map(x => x.rarity)).size}, minmax(0, 1fr))`}}>
                {counts.sort((a, b) => {
                    const diff = raritySort[a.rarity] - raritySort[b.rarity];
                    if (diff !== 0) {
                        return diff;
                    }
                    return a.count - b.count;
                }).map(({ rarity }) => <div key={rarity} className="h-2 flex bg-gray-200 dark:bg-slate-900">
                    <div style={{width: "70%"}} className={rarityBg[rarity]} />
                </div>)}
            </div>

        </div>
    </div>;
};

const Home: FunctionalComponent = () => {
    const { data } = useSWR<Game[]>('/games');
    const { data: counts } = useSWR<GameCount[]>('/games/counts');

    if (!data || !counts) return <p>Loading...</p>;

    const gamesIndex = new Fuse(data, {
        keys: ['name'],
        threshold: 0.3,
    });

    const [search, setSearch] = useState('');

    const games = useMemo(() => {
        if (!search) return data;
        return gamesIndex.search(search).map(x => x.item);
    }, [data, search]);

    return <>
        <input type="text" value={search} onInput={e => setSearch((e.target as HTMLInputElement).value.trim().toLowerCase())} placeholder="Search..." className="mb-4 w-full px-4 py-2 bg-gray-100 dark:bg-slate-800 dark:text-slate-400" />
        <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-8 auto-rows-[1fr]">
            {games.map(x => <GameCard key={x.id} game={x} counts={(counts ?? []).filter(({ id, rarity }) => id === x.id && rarity !== 'platinum')} />)}
        </div>
    </>;
};

export default Home;
