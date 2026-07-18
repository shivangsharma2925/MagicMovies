import { useState, useEffect, useRef, useMemo, useCallback } from "react";
import { useNavigate } from "react-router-dom";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faPlay, faCircleInfo } from "@fortawesome/free-solid-svg-icons";
import "../../styling/Hero.css";

const HERO_INTERVAL = 6000;

export default function HeroSection({ movies }) {
  const [current, setCurrent] = useState(0);
  const [fading, setFading] = useState(false);
  const timerRef = useRef(null);
  const navigate = useNavigate();

  const shuffleMovies = useCallback((movies) => {
    const shuffled = [...movies];

    for (let i = shuffled.length - 1; i > 0; i--) {
      const randomIndex = Math.floor(Math.random() * (i + 1));

      [shuffled[i], shuffled[randomIndex]] = [
        shuffled[randomIndex],
        shuffled[i],
      ];
    }

    return shuffled;
  }, []);

  const featured = useMemo(() => {
    return shuffleMovies(movies?.excellent ?? []).slice(0, 5);
  }, [movies?.excellent, shuffleMovies]);
  // const featured = movies?.excellent.slice(0, 5) ?? [];

  const goTo = (i) => {
    setFading(true);
    setTimeout(() => {
      setCurrent(i);
      setFading(false);
    }, 300);
  };

  const startTimer = () => {
    clearInterval(timerRef.current);
    timerRef.current = setInterval(() => {
      setCurrent((prev) => (prev + 1) % featured.length);
    }, HERO_INTERVAL);
  };

  useEffect(() => {
    if (!featured.length) return;
    startTimer();
    return () => clearInterval(timerRef.current);
  }, [featured.length]);

  if (!featured.length) return null;

  const movie = featured[current];

  return (
    <div className="hero-section">
      {/* <div
    className={`hero-backdrop ${fading ? "hero-fading" : ""}`}
    style={{ backgroundImage: `url(${movie.backdrop_path || movie.poster_path})` }}
  /> */}
      <div
        className={`hero-backdrop ${fading ? "hero-fading" : ""}`}
        style={{
          "--backdrop": `url(${movie.backdrop_path})`,
          "--poster": `url(${movie.poster_path})`,
        }}
      />
      <div className="hero-gradient" />

      <div className="hero-content">
        {movie.ranking?.ranking_name && (
          <span className="hero-badge">{movie.ranking.ranking_name}</span>
        )}
        <h1 className="hero-title">{movie.title}</h1>
        {movie.overview && <p className="hero-desc">{movie.overview}</p>}

        <div className="hero-actions">
          <button
            className="btn-hero-play"
            onClick={() => navigate(`/stream/${movie.youtube_id}`)}
          >
            <FontAwesomeIcon icon={faPlay} /> Play Trailer
          </button>

          <button
            className="btn-hero-icon"
            onClick={() => navigate(`/movie/detail/${movie.imdb_id}`)}
            aria-label="More info"
            title="More Info"
          >
            <FontAwesomeIcon icon={faCircleInfo} />
          </button>

          {/* <button
        className="btn-hero-icon"
        // onClick={() => handleAddToWatchlist(movie)}
        aria-label="Add to watchlist"
        title="Add to Watchlist"
      >
        <FontAwesomeIcon icon={faPlus} />
      </button> */}
        </div>
      </div>

      <div className="hero-dots">
        {featured.map((_, i) => (
          <button
            key={i}
            className={`hero-dot ${i === current ? "active" : ""}`}
            onClick={() => {
              goTo(i);
              startTimer();
            }}
          />
        ))}
      </div>

      <div className="hero-progress-bar">
        <div key={current} className="hero-progress-fill" />
      </div>
    </div>
  );
}
