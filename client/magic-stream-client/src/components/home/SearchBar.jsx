import Form from "react-bootstrap/Form";
import SuggestionsDropdown from "./SuggestionsDropdown";
import { useEffect, useRef } from "react";

const SearchBar = ({
  search,
  setSearch,
  suggestions,
  onSuggestionClick,
  showSuggestions,
  setShowSuggestions,
}) => {
  const searchRef = useRef(null);

  // close dropdown on outside click
  useEffect(() => {
    const handleOutsideClick = (event) => {
      if (searchRef.current && !searchRef.current.contains(event.target)) {
        setShowSuggestions(false);
      }
    };

    document.addEventListener("mousedown", handleOutsideClick);

    return () => {
      document.removeEventListener("mousedown", handleOutsideClick);
    };
  }, [setShowSuggestions]);

  return (
    <div ref={searchRef} className="search-input-wrapper">
      <span className="search-field-icon" aria-hidden="true">
        &#128269;
      </span>
      <Form.Control
        className="search-field"
        type="text"
        placeholder="Search movies..."
        value={search}
        onFocus={() => {
          if (suggestions.length > 0) {
            setShowSuggestions(true);
          }
        }}
        onChange={(e) => {
          setSearch(e.target.value);

          if (e.target.value.trim().length >= 2) {
            setShowSuggestions(true);
          } else {
            setShowSuggestions(false);
          }
        }}
      />
      {search && (
        <button
          className="search-clear-btn"
          onClick={() => {
            setSearch("");
            setShowSuggestions(false);
          }}
          aria-label="Clear search"
        >
          ×
        </button>
      )}
      {showSuggestions && (
        <SuggestionsDropdown
          suggestions={suggestions}
          onSuggestionClick={onSuggestionClick}
        />
      )}
    </div>
  );
};

export default SearchBar;
