import { useState } from "react";
import Button from "react-bootstrap/Button";
import { Link, useNavigate } from "react-router-dom";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faCirclePlay,
  faHeart,
  faHeartCrack,
  faCircleInfo,
} from "@fortawesome/free-solid-svg-icons";
import { useWatchlists } from "../../hook/UseWatchlists";
import AddToWatchlistPicker from "../Watchlist/AddToWatchlistPicker";
import { slugify } from "../../utils/SlugifyTitle";
import "../../styling/Movie.css";

const Movie = ({ movie }) => {
  const [showPicker, setShowPicker] = useState(false);

  const { data: watchLists = [] } = useWatchlists();
  const navigate = useNavigate();

  const isWatchlisted = watchLists.some((l) =>
    l.movies?.some((m) => m.imdb_id === movie.imdb_id),
  );

  const year = movie.release_date?.split("-")[0];
  const rating = movie.vote_average?.toFixed(1);

  const updateMovieReview = (imdb_id) => {
    navigate(`/review/${imdb_id}`);
  };

  return (
    <div className="col-6 col-md-4 col-lg-3 mb-4" id={movie._id}>
      <div className="movie-card-dark h-100">
        <Link
          to={`/stream/${slugify(movie.title)}/${movie.imdb_id}`}
          style={{ textDecoration: "none", color: "inherit" }}
        >
          <div className="movie-poster">
            <img
              src={movie.poster_path}
              alt={movie.title}
              className="movie-poster-img"
              loading="lazy"
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

            <button
              className={`movie-like-badge ${isWatchlisted ? "active" : ""}`}
              onClick={(e) => {
                e.preventDefault();
                e.stopPropagation();
                setShowPicker(true);
              }}
              aria-label={isWatchlisted ? "In watchlist" : "Add to watchlist"}
              title={isWatchlisted ? "In watchlist" : "Add to watchlist"}
            >
              <FontAwesomeIcon icon={isWatchlisted ? faHeart : faHeartCrack} />
            </button>
          </div>
        </Link>

        {/* Watchlist toggle — top right */}

        <div className="movie-info">
          <h6 className="movie-title">{movie.title}</h6>
          <div className="movie-meta-wrapper">
            <div className="movie-meta-row">
              {year && <span className="movie-year">{year}</span>}
              {year && <span className="movie-meta-dot">·</span>}
              {rating && <span className="movie-rating-badge">⭐ {rating}</span>}
            </div>
            <button
              className="mc-info-icon"
              onClick={(e) => {
                e.stopPropagation();
                navigate(`/movie/detail/${movie.imdb_id}`);
              }}
              aria-label="More info"
              title="More Info"
            >
              <FontAwesomeIcon icon={faCircleInfo} />
            </button>
          </div>
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

      {showPicker && (
        <AddToWatchlistPicker
          movie={movie}
          onClose={() => setShowPicker(false)}
        />
      )}
    </div>
  );
};

export default Movie;
