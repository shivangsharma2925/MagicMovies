import { useEffect, useState } from "react";
import axiosConfig from "../../api/axiosConfig";
import Movies from "../movies/Movies";
import { Button, Form } from "react-bootstrap";
import Spinner from "../../utils/Spinner";
import { useNavigate } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";

const Home = () => {
  const [search, setSearch] = useState("");
  const [debouncedSearch, setDebouncedSearch] = useState("");
  const navigate = useNavigate();
  let message = "";

  useEffect(() => {
    const timer = setTimeout(() => {
      setDebouncedSearch(search);
    }, 500); // 500ms delay

    return () => clearTimeout(timer);
  }, [search]);

  const trimmedSearch = debouncedSearch.trim();

  const {
    data: movies = [],
    isLoading,
    error,
  } = useQuery({
    queryKey: ["movies", trimmedSearch],

    queryFn: async () => {
      const url = trimmedSearch ? `/movies?search=${trimmedSearch}` : "/movies";

      const response = await axiosConfig.get(url);

      return response.data;
    },

    staleTime: 1000 * 60 * 5,

    keepPreviousData: true,
  });

  if (movies == null || movies.length === 0) {
    message = `There are currently no movies available ${trimmedSearch ? "for the provided search" : ""}`;
  }

  if (error) {
    console.log(error);
    message = "Error fetching movies";
  }

  return (
    <>
      <div
        className="sticky-top py-3"
        style={{
          zIndex: 1020,
          top: "60px",
          backgroundColor: "rgba(255, 255, 255, 0.7)", // Semi-transparent white
          backdropFilter: "blur(10px)", // Blurs content behind
          WebkitBackdropFilter: "blur(10px)", // Safari support
        }}
      >
        <div className="container-fluid d-flex align-items-center">
          <div style={{ flex: 1 }}></div>

          <div style={{ flex: 2 }}>
            <Form.Control
              className="shadow-sm border-0 bg-light"
              type="text"
              placeholder="Search movies..."
              style={{ borderRadius: "20px" }}
              value={search}
              onChange={(e) => setSearch(e.target.value)}
            />
          </div>

          <div style={{ flex: 1, textAlign: "right" }}>
            <Button
              variant="primary"
              size="sm"
              className="rounded-pill px-2 shadow-sm"
              onClick={() => navigate("admin/add-movie")}
            >
              Add Movies
            </Button>
          </div>
        </div>
      </div>
      {isLoading ? (
        <h2>
          <Spinner />
        </h2>
      ) : (
        <Movies movies={movies} message={message} />
      )}
    </>
  );
};

export default Home;
