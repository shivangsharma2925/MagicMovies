import Button from "react-bootstrap/Button";
import { Link, useNavigate } from "react-router-dom";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCirclePlay } from "@fortawesome/free-solid-svg-icons";
import "./Movie.css";

const Movie = ({ movie }) => {
  const navigate = useNavigate();

  const updateMovieReview = (imdb_id) => {
    navigate(`/review/${imdb_id}`);
  };

  return (
    <div className="col-6 col-md-4 col-lg-3 mb-4" id={movie._id}>
      <div className="movie-card-dark h-100">
        <Link
          to={`/stream/${movie.youtube_id?movie.youtube_id:null}`}
          style={{ textDecoration: "none", color: "inherit" }}
        >
          <div className="movie-poster">
            <img
              src={movie.poster_path}
              alt={movie.title}
              className="movie-poster-img"
            />
            <div className="movie-poster-overlay">
              <div className="movie-play-btn">
                <FontAwesomeIcon icon={faCirclePlay} />
              </div>
            </div>
            {movie.ranking?.ranking_name && (
              <span className="movie-rank-badge">
                {movie.ranking.ranking_name}
              </span>
            )}
          </div>
        </Link>
        <div className="movie-info">
          <h6 className="movie-title">{movie.title}</h6>
          <span className="movie-imdb">{movie.imdb_id}</span>
        </div>
        <div className="movie-actions">
          <button
            className="btn-movie-review"
            onClick={() => updateMovieReview(movie.imdb_id)}
          >
            Review
          </button>
        </div>
      </div>
    </div>
  );
};

export default Movie;
