import type {FunctionalComponent} from "preact";
import useSWR from "swr";

type Metadata = {
    name: string;
    backgroundImage?: string;
    metacritic: number;
    released: string;
    platforms: string[];
};

const GameHeader: FunctionalComponent<Metadata> = ({ name, metacritic, released, platforms }) => {
    return <div className="container px-4 py-16 flex flex-col justify-center items-start space-y-4">
        <h2 className="text-4xl font-extrabold">{name}</h2>
        <div className="space-x-4">
            {platforms.map((platform, i) => <span key={i} className="bg-gray-100 text-gray-700 px-2 py-1 font-semibold">{platform}</span>)}
            <span className="bg-gray-100 text-gray-700 px-2 py-1 font-semibold">{metacritic}% Metacritic</span>
            <span className="bg-gray-100 text-gray-700 px-2 py-1 font-semibold">Released on {released}</span>
        </div>
    </div>;
};

const Page: FunctionalComponent<{ game: string }> = ({ game }) => {
    const { data: metadata } = useSWR<Metadata>(`/metadata?id=${encodeURIComponent(game)}`)

    if (!metadata) return <p>Loading...</p>

    return <>
        {metadata.backgroundImage ? <div className="h-96 flex items-center bg-cover bg-center" style={{backgroundImage: `url(${metadata.backgroundImage})`}}>
            <div className="backdrop-blur-sm w-full bg-black/10">
                <GameHeader {...metadata} />
            </div>
        </div> : <div className="bg-slate-800">
            <GameHeader {...metadata} />
        </div>}
    </>
};

export default Page;
