import { useParams, useNavigate } from "react-router-dom";
import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faArrowLeft,
  faPlay,
  faCheck,
  faPlus,
} from "@fortawesome/free-solid-svg-icons";
import { fetchTitles } from "../home/GetMovieServices";
import Spinner from "../../utils/Spinner";
import { useWatchlists } from "../../hook/UseWatchlists";
import AddToWatchlistPicker from "../Watchlist/AddToWatchlistPicker";
import PageTitle from "../../utils/PageTitle";
import { Helmet } from "react-helmet-async";
import { slugify } from "../../utils/SlugifyTitle";
import "../../styling/Category.css";

// Emoji/label map for known category names
const CATEGORY_META = {
  Excellent: { label: "Top Rated", color: "#fbbf24" },
  Good: { label: "Good Picks", color: "rgb(74, 222, 128)" },
  Okay: { label: "Casual Watch", color: "rgb(96, 165, 250)" },
  newRelease: { label: "Recently Added", color: "#ec5050df" },
  Bad: { label: "Category", color: "#ff6b8a" },
  Terrible: { label: "Category", color: "#ff6b8a" },
};

export default function CategoryMovies() {
  const { title } = useParams();
  const navigate = useNavigate();
  const [selectedMovie, setSelectedMovie] = useState(null);

  const meta = CATEGORY_META[title] ?? {
    label: "Category",
    color: "#f32256da",
  };

  const { data: watchLists = [] } = useWatchlists();

  const {
    data: movies = [],
    isLoading,
    error,
  } = useQuery({
    queryKey: ["category", title],
    queryFn: () => fetchTitles({ title }),
    staleTime: 5 * 60 * 60 * 1000,
  });

  const handleWatchlist = (e, movie) => {
    e.preventDefault();
    e.stopPropagation();
    setSelectedMovie(movie);
  };

  const description = `watch ${title} movies trailer and explore movie details, genres, release date, and other info.`;
  const pageUrl = `https://www.magicmovies.site/category/${title}`;

  if (isLoading) return <Spinner />;

  return (
    <div className="cat-page">
      <Helmet>
        <title>{`${title} Movies — Magic Stream`}</title>

        {/* Search engine description */}
        <meta name="description" content={description} />

        {/* Search engine crawling */}
        <meta name="robots" content="index, follow" />

        {/* Main page URL */}
        <link rel="canonical" href={pageUrl} />

        {/* Open Graph - control how page appears when shared on platforms such as WhatsApp, Facebook, LinkedIn, etc */}
        <meta property="og:title" content={`${title}`} />
        <meta property="og:description" content={description} />
        <meta property="og:url" content={pageUrl} />
        <meta property="og:type" content="video.movie" />
        <meta property="og:site_name" content="Magic Stream" />
      </Helmet>
      {/* Hero header */}
      <div className="cat-hero">
        <div className="cat-hero-bg" style={{ "--bg-color": meta.color }} />
        <div className="cat-hero-glow" />
        <button className="cat-back" onClick={() => navigate(-1)}>
          <FontAwesomeIcon icon={faArrowLeft} /> Back
        </button>
        <div className="cat-hero-content">
          <div className="cat-label" style={{ color: meta.color }}>
            {meta.label}
          </div>
          <h1 className="cat-title">{title}</h1>
          <div className="cat-count">
            {error
              ? "Error loading movies"
              : `${movies.length} movie${movies.length !== 1 ? "s" : ""}`}
          </div>
        </div>
      </div>

      {/* Grid */}
      {error ? (
        <div className="cat-empty">
          <div className="cat-empty-icon">⚠️</div>
          <div className="cat-empty-text">
            Failed to load movies for "{title}"
          </div>
        </div>
      ) : movies.length === 0 ? (
        <div className="cat-empty">
          <div className="cat-empty-icon">🎬</div>
          <div className="cat-empty-text">No movies found in this category</div>
        </div>
      ) : (
        <>
          <div className="cat-grid">
            {movies.map((movie) => {
              const year = movie.release_date?.split("-")[0];
              const rating = movie.vote_average?.toFixed(1);

              const watchlisted = watchLists.some((l) =>
                l.movies?.some((m) => m.imdb_id === movie?.imdb_id),
              );

              return (
                <div
                  key={movie._id}
                  className="cat-card"
                  onClick={() => navigate(`/stream/${slugify(movie.title)}/${movie.imdb_id}`)}
                >
                  <div className="cat-poster">
                    <img
                      src={movie.poster_path}
                      alt={movie.title}
                      loading="lazy"
                    />
                    <div className="cat-overlay">
                      <div className="cat-play">
                        <FontAwesomeIcon icon={faPlay} />
                      </div>
                    </div>
                    {movie.ranking?.ranking_name && (
                      <span className="cat-rank-badge">
                        {movie.ranking.ranking_name}
                      </span>
                    )}
                  </div>
                  <div className="cat-info">
                    <div className="cat-movie-title">{movie.title}</div>
                    <div className="cat-movie-meta-wrapper">
                      <div className="cat-movie-meta">
                        {year && <span>{year}</span>}
                        {year && rating && <span>·</span>}
                        {rating && (
                          <span className="cat-movie-rating">⭐ {rating}</span>
                        )}
                      </div>
                      <button
                        type="button"
                        className={`cat-movie-watchlist ${watchlisted ? "active" : ""}`}
                        onClick={(e) => handleWatchlist(e, movie)}
                        aria-label={
                          watchlisted ? "In Watchlist" : "Add to Watchlist"
                        }
                        title={
                          watchlisted ? "In Watchlist" : "Add to Watchlist"
                        }
                      >
                        <FontAwesomeIcon
                          icon={watchlisted ? faCheck : faPlus}
                        />
                      </button>
                    </div>
                  </div>
                </div>
              );
            })}
          </div>
          {selectedMovie && (
            <AddToWatchlistPicker
              movie={selectedMovie}
              onClose={() => setSelectedMovie(null)}
            />
          )}
        </>
      )}
    </div>
  );
}
