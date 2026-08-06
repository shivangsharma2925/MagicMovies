import { useParams, useNavigate } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import { useState } from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faPlay,
  faThumbsUp,
  faThumbsDown,
  faPlus,
  faCheck,
} from "@fortawesome/free-solid-svg-icons";
import { Helmet } from "react-helmet-async";
import { fetchMovie } from "../home/GetMovieServices";
import Spinner from "../../utils/Spinner";
import { useWatchlists } from "../../hook/UseWatchlists";
import AddToWatchlistPicker from "../Watchlist/AddToWatchlistPicker";
import FormatDuration from "../../hook/UseFormatDurationSEO";
import PageTitle from "../../utils/PageTitle";
import { slugify } from "../../utils/SlugifyTitle";
import "../../styling/MovieDetail.css";
// import WatchlistPanel from "../Watchlist/AddToWatchlistModal";

const MovieDetails = () => {
  const { imdbId } = useParams();
  const navigate = useNavigate();
  const [liked, setLiked] = useState(false);
  const [disliked, setDisliked] = useState(false);
  const [showFullOverview, setShowFullOverview] = useState(false);
  const [showPicker, setShowPicker] = useState(false);

  const {
    data: movie,
    isLoading,
    error,
  } = useQuery({
    queryKey: ["movie", imdbId],
    queryFn: () => fetchMovie(imdbId),
    staleTime: 1000 * 60 * 10,
  });

  const { data: watchLists = [] } = useWatchlists();
  const watchlisted = watchLists.some((l) =>
    l.movies?.some((m) => m.imdb_id === movie?.imdb_id),
  );

  if (isLoading) {
    return (
      <div
        style={{
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          height: "100vh",
        }}
      >
        <Spinner />
      </div>
    );
  }

  if (error || !movie) {
    return (
      <div
        style={{
          textAlign: "center",
          padding: "80px",
          color: "rgba(255,255,255,0.4)",
        }}
      >
        Movie not found.
      </div>
    );
  }

  const year = movie.release_date?.split("-")[0];
  const hours = movie.runtime ? Math.floor(movie.runtime / 60) : null;
  const mins = movie.runtime ? movie.runtime % 60 : null;
  const runtime = hours ? `${hours}h ${mins}m` : null;
  const rating = movie.vote_average?.toFixed(1);

  const handleLike = () => {
    setLiked((p) => !p);
    setDisliked(false);
  };

  const handleDislike = () => {
    setDisliked((p) => !p);
    setLiked(false);
  };

  const handleWatchlist = () => {
    setShowPicker((p) => !p);
  };

  const title = movie?.title ? `${movie.title} — Magic Stream` : "Magic Stream";
  const description = movie?.overview
    ? movie.overview.length > 155
      ? movie.overview.slice(0, 152) + "..."
      : movie.overview
    : `Watch ${movie?.title ?? "movie"} trailer and explore ratings, genres and recommendations.`;
  const pageUrl = `https://www.magicmovies.site/movie/detail/${movie.imdb_id}`;
  const imageUrl = movie.backdrop_path;

  const genres = movie?.genre?.map((genre) => genre.genre_name) ?? [];

  const movieStructuredData = movie
    ? {
        "@context": "https://schema.org",
        "@type": "Movie",

        "@id": `https://www.magicmovies.site/movie/detail/${movie.imdb_id}#movie`,

        url: `https://www.magicmovies.site/movie/detail/${movie.imdb_id}`,

        name: movie.title,

        description: movie.overview,

        image: [movie.poster_path, movie.backdrop_path].filter(Boolean),

        dateCreated: movie.release_date,

        genre: genres,

        duration: FormatDuration(movie.runtime),

        contentRating: movie.adult ? "Adults Only" : undefined,

        sameAs: `https://www.imdb.com/title/${movie.imdb_id}/`,

        trailer: movie.youtube_id
          ? {
              "@type": "VideoObject",
              name: `${movie.title} Official Trailer`,
              description: `Watch the trailer for ${movie.title}.`,
              thumbnailUrl: `https://img.youtube.com/vi/${movie.youtube_id}/maxresdefault.jpg`,
              embedUrl: `https://www.youtube.com/embed/${movie.youtube_id}`,
              contentUrl: `https://www.youtube.com/watch?v=${movie.youtube_id}`,
            }
          : undefined,
      }
    : null;

  return (
    <div className="md-page">
      <Helmet>
        <title>{title}</title>

        {/* Search engine description */}
        <meta name="description" content={description} />

        {/* Search engine crawling */}
        <meta name="robots" content="index, follow" />

        <meta content={`${movie.title} Poster`} />

        {/* Main page URL */}
        <link rel="canonical" href={pageUrl} />

        {/* Open Graph - control how page appears when shared on platforms such as WhatsApp, Facebook, LinkedIn, etc */}
        <meta property="og:title" content={title} />
        <meta property="og:description" content={description} />
        <meta property="og:image" content={imageUrl} />
        <meta property="og:url" content={pageUrl} />
        <meta property="og:type" content="video.movie" />
        <meta property="og:site_name" content="Magic Stream" />

        {/* JSON-LD(JavaScript Object Notation for Linked Data), json like format useful for web crawlers and SEO */}
        <script type="application/ld+json">
          {JSON.stringify(movieStructuredData)}
        </script>
      </Helmet>

      <div className="md-hero">
        <img
          className="md-bg-img"
          src={movie.backdrop_path}
          alt={movie.title}
        />
        <div className="md-bg-fade" />

        <button className="md-back" onClick={() => navigate(-1)}>
          ← Back
        </button>

        <div className="md-title-block">
          <h1 className="md-title">{movie.title}</h1>
          {movie.tagline && <p className="md-tagline">"{movie.tagline}"</p>}
        </div>

        <div className="md-content">
          {movie.ranking?.ranking_name && (
            <div className="md-award">✦ {movie.ranking.ranking_name}</div>
          )}

          {movie.overview && (
            <p
              className={`md-overview ${!showFullOverview ? "md-overview-clamp" : ""}`}
              onClick={() => setShowFullOverview((prev) => !prev)}
            >
              ’’ {movie.overview}
              <span className="md-read-more">
                {showFullOverview ? " Show less" : ""}
              </span>
            </p>
          )}
          <div className="md-meta">
            {movie.genre?.map((g) => (
              <button
                key={g.genre_id}
                className="md-pill pill-genre"
                onClick={() => navigate(`/category/${g.genre_name}`)}
              >
                {g.genre_name}
              </button>
            ))}
            {rating && (
              <span className="md-pill pill-rating">⭐ {rating} / 10</span>
            )}
            {year && <span className="md-pill pill-info">{year}</span>}
            {runtime && <span className="md-pill pill-info">{runtime}</span>}
            <span
              className={`md-pill ${movie.is_adult ? "pill-adult" : "pill-pg"}`}
            >
              {movie.is_adult ? "18+" : "PG-13"}
            </span>
          </div>

          <div className="md-actions-row">
            <button
              className="md-primary-btn"
              onClick={() => navigate(`/stream/${slugify(movie.title)}/${movie.imdb_id}`)}
            >
              <FontAwesomeIcon icon={faPlay} />
              Watch Trailer
            </button>

            {[
              {
                icon: watchlisted ? faCheck : faPlus,
                label: watchlisted ? "In Watchlist" : "Add to Watchlist",
                cls: watchlisted ? "active-watch" : "",
                action: handleWatchlist,
              },
              {
                icon: faThumbsUp,
                label: "Like",
                cls: liked ? "active-like" : "",
                action: handleLike,
              },
              {
                icon: faThumbsDown,
                label: "Dislike",
                cls: disliked ? "active-dislike" : "",
                action: handleDislike,
              },
            ].map(({ icon, label, cls, action }) => (
              <div key={label} className="md-tooltip-wrap">
                <button
                  key={label}
                  className={`md-round ${cls}`}
                  onClick={action}
                  aria-label={label}
                >
                  <FontAwesomeIcon icon={icon} />
                </button>
                <span className="md-tooltip">{label}</span>
              </div>
            ))}
          </div>

          {showPicker && (
            <AddToWatchlistPicker
              movie={movie}
              onClose={() => setShowPicker(false)}
            />
          )}

          {movie.admin_review && (
            <div className="md-review-box">
              <div className="md-review-lbl">Admin Review</div>
              <div className="md-review-txt">"{movie.admin_review}"</div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default MovieDetails;
