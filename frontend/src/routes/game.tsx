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
  backgroundImageURL?: string;
  metacriticScore?: number;
  releaseDate?: string;
};

const GameHeader: FunctionalComponent<Game> = ({ name, metacriticScore, releaseDate, platform }) => {
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

const Page: FunctionalComponent<{ game: string }> = ({ game }) => {
  const { data: metadata } = useSWR<Game>(`/games/${encodeURIComponent(game)}`);

  if (!metadata) return <p>Loading...</p>;

  return (
    <>
      {metadata.backgroundImageURL
        ? (
          <div
            className="h-96 flex items-center bg-cover bg-center"
            style={{ backgroundImage: `url(${metadata.backgroundImageURL})` }}
          >
            <div className="backdrop-blur-sm w-full bg-black/10">
              <GameHeader {...metadata} />
            </div>
          </div>
        )
        : (
          <div className="bg-slate-800">
            <GameHeader {...metadata} />
          </div>
        )}
    </>
  );
};

export default Page;
