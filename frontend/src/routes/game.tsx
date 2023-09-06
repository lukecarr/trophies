import type { FunctionalComponent } from "preact";
import useSWR from "swr";

type Game = {
  id: number;
  psnID: string;
  psnServiceName: string;
  name: string;
  description: string;
  iconURL: string;
  platform: string;
};

type Metadata = {
  backgroundImageURL?: string;
  metacriticScore?: number;
  releaseDate?: string;
};

const GameHeader: FunctionalComponent<Game & Metadata> = ({ name, metacriticScore, releaseDate, platform }) => {
  return (
    <div className="container px-4 py-16 flex flex-col justify-center items-start space-y-4">
      <h2 className="text-4xl text-white font-extrabold">{name}</h2>
      <div className="space-x-4">
        {platform.split(",").map((platform, i) => (
          <span key={i} className="bg-gray-100 text-gray-700 px-2 py-1 font-semibold">{platform}</span>
        ))}
        {metacriticScore && metacriticScore > 0 && (
          <span className="bg-gray-100 text-gray-700 px-2 py-1 font-semibold">{metacriticScore}% Metacritic</span>
        )}
        <span className="bg-gray-100 text-gray-700 px-2 py-1 font-semibold">Released on {releaseDate}</span>
      </div>
    </div>
  );
};

const Page: FunctionalComponent<{ id: string }> = ({ id }) => {
  const { data: game } = useSWR<Game>(`/games/${encodeURIComponent(id)}`);
  const { data: metadata } = useSWR<Metadata>(`/games/${encodeURIComponent(id)}/metadata`);

  if (!game || !metadata) return <p>Loading...</p>;

  return (
    <>
      {metadata.backgroundImageURL
        ? (
          <div
            className="h-96 flex items-center bg-cover bg-center"
            style={{ backgroundImage: `url(${metadata.backgroundImageURL})` }}
          >
            <div className="backdrop-blur-sm w-full bg-black/10">
              <GameHeader {...metadata} {...game} />
            </div>
          </div>
        )
        : (
          <div className="bg-slate-800">
            <GameHeader {...metadata} {...game} />
          </div>
        )}
    </>
  );
};

export default Page;
