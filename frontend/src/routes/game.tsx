import type {FunctionalComponent} from "preact";
import useSWR from "swr";

type Metadata = {
    name: string;
    backgroundImage?: string;
};

const Page: FunctionalComponent<{ game: string }> = ({ game }) => {
    const { data: metadata } = useSWR<Metadata>(`/metadata?id=${encodeURIComponent(game)}`)

    if (!metadata) return <p>Loading...</p>

    return <>
        {metadata.backgroundImage ? <div className="h-96 flex items-center bg-cover bg-center" style={{backgroundImage: `url(${metadata.backgroundImage})`}}>
            <div className="backdrop-blur-sm w-full bg-black/10">
                <div className="container px-4 py-16 flex items-center">
                    <h2 className="text-4xl font-extrabold">{metadata.name}</h2>
                </div>
            </div>
        </div> : <div className="bg-slate-800">
            <div className="container px-4 py-16 flex items-center">
                <h2 className="text-4xl font-extrabold">{metadata.name}</h2>
            </div>
        </div>}
    </>
};

export default Page;
