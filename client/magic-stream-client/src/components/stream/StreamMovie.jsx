import { useParams } from "react-router-dom";
import ReactPlayer from "react-player";
import { useQuery } from "@tanstack/react-query";
import { fetchMovie } from "../home/GetMovieServices";
import PageTitle from "../../utils/PageTitle";

const styles = {
  playerContainer: {
    height: "90vh",
  },
};

const StreamMovie = () => {
  let params = useParams();
  let imdbId = params.imdb_id;

  const {
    data: movie,
    isLoading,
    isError,
  } = useQuery({
    queryKey: ["movie", imdbId],
    queryFn: () => fetchMovie(imdbId),
    staleTime: 10 * 60 * 1000,
  });

  if (isLoading) {
    return (
      <h2 className="m-auto text-white text-center">
        Loading...
      </h2>
    );
  }

  if (isError || !movie?.youtube_id) {
    return (
      <h2 className="m-auto text-white text-center">
        Trailer not available.
      </h2>
    );
  }

  const Player = ReactPlayer.default ?? ReactPlayer;

  return (
    <div style={styles.playerContainer}>
      <PageTitle title={`Streaming - ${movie.title}`} />

      <Player
        controls
        playing
        url={`https://www.youtube.com/watch?v=${movie.youtube_id}`}
        width="100%"
        height="100%"
      />
    </div>
  );
};

export default StreamMovie;
