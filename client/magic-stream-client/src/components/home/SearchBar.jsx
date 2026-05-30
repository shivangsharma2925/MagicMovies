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
    <div
      ref={searchRef}
      className="position-relative"
      style={{
        width: "100%",
      }}
    >
      <Form.Control
        className="shadow-sm border-0 bg-light py-2 px-4"
        type="text"
        placeholder="Search movies..."
        style={{
          borderRadius: "999px",
          fontSize: "1.05rem",
        }}
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
