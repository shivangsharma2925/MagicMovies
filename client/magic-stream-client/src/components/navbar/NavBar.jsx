import { useState, useRef, useEffect } from "react";
import { NavLink, useNavigate } from "react-router-dom";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faMagnifyingGlass,
  faCamera,
  faCircleUser,
  faDownLong,
  faBars,
  faArrowRightToBracket,
  faPersonCirclePlus,
  faTrashCan,
  faPlus,
  faArrowRightFromBracket,
  faUser,
} from "@fortawesome/free-solid-svg-icons";
import UseAuth from "../../hook/UseAuth";
import axiosConfig from "../../api/axiosConfigwocredential";
import logo from "../../../logo.png";
import SearchBar from "../home/SearchBar";
import { useQuery } from "@tanstack/react-query";
import useDebounce from "../../hook/UseDebounce";
import { fetchMovieSuggestions, fetchMovies } from "../home/GetMovieServices";
import SuggestionsDropdown from "../home/SuggestionsDropdown";
import Spinner from "../../utils/Spinner";
import SearchedMovies from "../movie/SearchedMovies";
import "../../styling/Header.css";
import "../../styling/SearchSuggestion.css";

export default function Header() {
  const { auth, logout } = UseAuth();
  const [genres, setGenres] = useState([]);
  const [genresLoaded, setGenresLoaded] = useState(false); // cached, no extra fetch
  const [searchOpen, setSearchOpen] = useState(false);
  const [search, setSearch] = useState("");
  const [acctOpen, setAcctOpen] = useState(false);
  const [genreOpen, setGenreOpen] = useState(false);
  const [menuOpen, setMenuOpen] = useState(false);
  const [selectedGenre, setSelectedGenre] = useState(null);
  const [showSuggestions, setShowSuggestions] = useState(false);

  const inputRef = useRef(null);
  const acctRef = useRef(null);
  const genreRef = useRef(null);
  const menuRef = useRef(null);
  const navigate = useNavigate();

  const initials = auth
    ? `${auth.first_name?.[0] ?? ""}${auth.last_name?.[0] ?? ""}`.toUpperCase()
    : "";

  useEffect(() => {
    const FetchGenres = async () => {
      try {
        if (genresLoaded) return;
        const genres = await axiosConfig.get("/genres");
        const genredata = genres.data;
        setGenres(genredata);
        setGenresLoaded(true);
      } catch (err) {
        console.log(err);
      }
    };
    FetchGenres();
  }, [genresLoaded]);

  const debouncedSearch = useDebounce(search, 500);

  const trimmedSearch = debouncedSearch.trim();

  const { data: searchedMovies = [], isLoading } = useQuery({
    queryKey: ["searchedMovies", trimmedSearch],

    queryFn: () => fetchMovies({ search: trimmedSearch }),

    enabled: trimmedSearch.trim().length >= 2,

    staleTime: 1000 * 60 * 5,

    refetchOnWindowFocus: false,
  });

  /* fetching suggestions */
  const { data: suggestions = [] } = useQuery({
    queryKey: ["movie-suggestions", debouncedSearch],

    queryFn: () => fetchMovieSuggestions(debouncedSearch),

    enabled: debouncedSearch.length >= 2,

    staleTime: 1000 * 60 * 10,
  });

  // Close dropdowns on outside click
  useEffect(() => {
    const handler = (e) => {
      if (!acctRef.current?.contains(e.target)) setAcctOpen(false);
      if (!genreRef.current?.contains(e.target)) {
        setGenreOpen(false);
        setSelectedGenre(null);
      }
      if (!menuRef.current?.contains(e.target)) setMenuOpen(false);
    };
    document.addEventListener("mousedown", handler);
    return () => document.removeEventListener("mousedown", handler);
  }, []);

  const openSearch = () => {
    setSearchOpen(true);
    setAcctOpen(false);
    setGenreOpen(false);
    setTimeout(() => inputRef.current?.focus(), 150);
  };

  const closeSearch = () => {
    setSearchOpen(false);
    setSearch("");
    setShowSuggestions(false);
  };

  const handleGenreClick = (genre) => {
    setSelectedGenre((prev) =>
      prev === genre.genre_name ? null : genre.genre_name,
    );
    if (selectedGenre === genre.genre_name) {
      navigate("/");
    } else navigate(`/category/${genre.genre_name}`);
    setGenreOpen(false);
  };

  return (
    <>
      <nav className="app-navbar justify-content-between">
        {/* Brand */}
        <NavLink to="/" className="nb-brand-left">
          <img src={logo} alt="Magic Stream" className="nb-logo" />
          <span className="nb-brand-name">Magic Stream</span>
        </NavLink>

        <div ref={menuRef}>
          <button
            className={`menu-btn ${menuOpen ? "active" : ""}`}
            onClick={() => setMenuOpen(!menuOpen)}
          >
            Menu{" "}
            <FontAwesomeIcon
              icon={faDownLong}
              className={`${menuOpen ? "rotate-x-180" : ""}`}
            />
          </button>

          <button className="menu-icon" onClick={() => setMenuOpen(!menuOpen)}>
            <FontAwesomeIcon icon={faBars} />
          </button>
          {menuOpen && (
            <div className={"menu-nb-links"}>
              <NavLink
                onClick={() => setMenuOpen(false)}
                to="/"
                end
                className={({ isActive }) =>
                  `nb-link${isActive ? " active" : ""}`
                }
              >
                Home
              </NavLink>
              <NavLink
                to="/recommended"
                onClick={() => setMenuOpen(false)}
                className={({ isActive }) =>
                  `nb-link${isActive ? " active" : ""}`
                }
              >
                Recommended
              </NavLink>
              <NavLink
                to="/watchlist"
                onClick={() => setMenuOpen(false)}
                className={({ isActive }) =>
                  `nb-link${isActive ? " active" : ""}`
                }
              >
                Watch List
              </NavLink>
            </div>
          )}
        </div>

        <NavLink to="/" className="nb-brand-mid">
          <img src={logo} alt="Magic Stream" className="nb-logo" />
          <span className="nb-brand-name">Magic Stream</span>
        </NavLink>

        {/* Nav links */}
        <div className={"nb-links"}>
          <NavLink
            to="/"
            end
            className={({ isActive }) => `nb-link${isActive ? " active" : ""}`}
          >
            Home
          </NavLink>
          <NavLink
            to="/recommended"
            className={({ isActive }) => `nb-link${isActive ? " active" : ""}`}
          >
            Recommended
          </NavLink>
          <NavLink
            to="/watchlist"
            className={({ isActive }) => `nb-link${isActive ? " active" : ""}`}
          >
            Watch List
          </NavLink>
        </div>

        {/* Right icons */}
        <div className="nb-right">
          {/* Search */}
          <button
            className="nb-icon-btn"
            onClick={openSearch}
            aria-label="Search"
            title="Search movies"
          >
            <FontAwesomeIcon icon={faMagnifyingGlass} />
          </button>

          {/* Genre / Camera */}
          <div ref={genreRef} className="nb-dd-wrap">
            <button
              className={`nb-icon-btn ${genreOpen ? "active" : ""}`}
              onClick={() => {
                setGenreOpen((p) => !p);
                setAcctOpen(false);
              }}
              aria-label="Browse genres"
              title="Browse genres"
            >
              <FontAwesomeIcon icon={faCamera} />
            </button>

            {genreOpen && (
              <div className="nb-genre-panel">
                <p className="nb-panel-label">Browse by genre</p>
                <div className="nb-genre-grid">
                  {genres.map((genre) => (
                    <button
                      key={genre.genre_id}
                      className={`nb-genre-chip ${selectedGenre === genre.genre_name ? "selected" : ""}`}
                      onClick={() => handleGenreClick(genre)}
                    >
                      {genre.genre_name}
                    </button>
                  ))}
                </div>
              </div>
            )}
          </div>

          {/* Account */}
          <div ref={acctRef} className="nb-dd-wrap">
            {auth ? (
              <button
                className={`nb-avatar-btn ${acctOpen ? "open" : ""}`}
                onClick={() => {
                  setAcctOpen((p) => !p);
                  setGenreOpen(false);
                }}
                aria-label="Account menu"
                aria-expanded={acctOpen}
                title="View Profile"
              >
                {initials}
              </button>
            ) : (
              <button
                className={`nb-icon-btn ${acctOpen ? "active" : ""}`}
                onClick={() => {
                  setAcctOpen((p) => !p);
                  // setGenreOpen(false);
                }}
                aria-label="Account"
                title="Account"
              >
                <FontAwesomeIcon icon={faCircleUser} />
              </button>
            )}

            {acctOpen && (
              <div className="nb-dropdown">
                {auth ? (
                  <>
                    <div className="nb-dd-header">
                      <div className="nb-dd-avatar">{initials}</div>
                      <div>
                        <div className="nb-dd-name">
                          {auth.first_name} {auth.last_name}
                        </div>
                        <div className="nb-dd-sub">Magic Stream member</div>
                      </div>
                    </div>
                    <div className="nb-dd-body">
                      <button
                        className="nb-dd-item"
                        onClick={() => {
                          navigate("/profile");
                          setAcctOpen(false);
                        }}
                      >
                        <FontAwesomeIcon icon={faUser} /> Profile
                      </button>
                      <button
                        className="nb-dd-item"
                        onClick={() => {
                          navigate("/admin/add-movie");
                          setAcctOpen(false);
                        }}
                      >
                        <FontAwesomeIcon icon={faPlus} /> Add Movie
                      </button>
                      <div className="nb-dd-divider" />
                      <button
                        className="nb-dd-item danger"
                        onClick={() => {
                          logout();
                          setAcctOpen(false);
                        }}
                      >
                        <FontAwesomeIcon icon={faArrowRightFromBracket} />{" "}
                        Logout
                      </button>
                      <button
                        className="nb-dd-item danger"
                        onClick={() => {
                          navigate("/account/delete");
                          setAcctOpen(false);
                        }}
                      >
                        <FontAwesomeIcon icon={faTrashCan} /> Delete Account
                      </button>
                    </div>
                  </>
                ) : (
                  <div className="nb-dd-body">
                    <button
                      className="nb-dd-item"
                      onClick={() => {
                        navigate("/login");
                        setAcctOpen(false);
                      }}
                    >
                      <FontAwesomeIcon icon={faArrowRightToBracket} /> Login
                    </button>
                    <button
                      className="nb-dd-item"
                      onClick={() => {
                        navigate("/register");
                        setAcctOpen(false);
                      }}
                    >
                      <FontAwesomeIcon icon={faPersonCirclePlus} /> Register
                    </button>
                  </div>
                )}
              </div>
            )}
          </div>
        </div>
      </nav>

      {/* Search overlay — rendered outside navbar so it covers full page */}
      {searchOpen && (
        <div
          className="nb-search-overlay"
          onClick={(e) => e.target === e.currentTarget && closeSearch()}
        >
          <div className="nb-search-content">
            <SearchBar
              search={search}
              setSearch={setSearch}
              suggestions={suggestions}
              showSuggestions={showSuggestions}
              setShowSuggestions={setShowSuggestions}
              setSearchOpen={setSearchOpen}
            />
            {showSuggestions && (
              <SuggestionsDropdown
                suggestions={suggestions}
                onSuggestionClick={(title) => {
                  setSearch(title);
                  setShowSuggestions(false);
                }}
              />
            )}
            {isLoading ? (
              <Spinner />
            ) : (
              <SearchedMovies
                movies={searchedMovies?.movies}
                search={trimmedSearch}
                setSearchOpen={setSearchOpen}
                setSearch={setSearch}
                setShowSuggestions={setShowSuggestions}
              />
            )}
          </div>
        </div>
      )}
    </>
  );
}
