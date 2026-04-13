import { useParams } from "react-router-dom";
import ReactPlayer from "react-player";

const styles = {
  playerContainer: {
    height: "90vh",
  },
};

const StreamMovie = () => {
  let params = useParams();
  let key = params.yt_id;
  const Player = ReactPlayer.default;

  return (
    <div style={styles.playerContainer}>
      {key != null ? (
        <Player
          controls={true}
          playing={true}
          url={`https://www.youtube.com/watch?v=${key}`}
          width="100%"
          height="100%"
        />
      ) : null}
    </div>
  );
};

export default StreamMovie;
