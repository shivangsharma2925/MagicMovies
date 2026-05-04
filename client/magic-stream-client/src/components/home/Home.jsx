import { useEffect, useRef, useState } from "react";
import axiosConfig from "../../api/axiosConfig";
import Movies from "../movies/Movies";
import { Button, Form } from "react-bootstrap";
import Spinner from "../../utils/Spinner";
import { useNavigate } from "react-router-dom";
import { useInfiniteQuery } from "@tanstack/react-query";

const Home = () => {
  const [search, setSearch] = useState("");
  const [debouncedSearch, setDebouncedSearch] = useState("");
  const navigate = useNavigate();
  const observerRef = useRef();
  let message = "";

  useEffect(() => {
    const timer = setTimeout(() => {
      setDebouncedSearch(search);
    }, 500); // 500ms delay

    return () => clearTimeout(timer);
  }, [search]);

  const trimmedSearch = debouncedSearch.trim();

  const fetchMovies = async ({ pageParam = 1 }) => {
    const url = trimmedSearch
      ? `/movies?search=${trimmedSearch}&page=${pageParam}&limit=10`
      : `/movies?page=${pageParam}&limit=10`;

    const response = await axiosConfig.get(url);

    return response.data;
  };

  /*
    The data in the react query movies array is stored as an set of all the pages we have ordered so far from server,
    example,
    movies = {
        "pages": [
            {
                "movies": [first 10 movies], "page": 1, "hasmore": true
            },
            {
                "movies": [second 10 movies], "page": 2, "hasmore": true
            },
        ],
        pageParams: [1, 2]
    } 

    Complete Process,
    0. Initially we have pageparam as 1 to show data.
    1. Intersecton oberser keep the track of the marked element (in this case the div at the bottom), as soon as that element intersects or visible on the screen, the event fires up.
    2. inside the event we have "hasNextPage", which uses the "fetchNextPage", a messesnger that tells react query to fetch next page. But react query doesn't know the value of next page so it calls the function "getNextPageParam" by giving the last page object, the function returns the next page value by looking at the hasMore flag we send from backend.
    3. The next page value now gets stored inside the pageParams array.
    4. Now the queryFn gets executed and the latest value of the pageParams array is used for getting the next set of movies.
    5. New data gets stored similarily in the pages array as another object, then we use the flatMap to get all the movies in single array.
    6. React query will store all the pages data in its cache giving almost 0 second latency to show data.
  */

  const {
    data: movies,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
    isLoading,
    error,
  } = useInfiniteQuery({
    queryKey: ["movies", trimmedSearch],

    queryFn: fetchMovies,

    initialPageParam: 1,

    getNextPageParam: (lastPage) => {
      return lastPage.hasMore ? lastPage.page + 1 : undefined;
    },

    placeholderData: (prev) => prev,

    staleTime: 1000 * 60 * 5,
  });

  useEffect(() => {
    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting && hasNextPage) {
          fetchNextPage();
        }
      },
      {
        rootMargin: "100px", //buffer till when the target element intersects or visible on the screen
        threshold: 0.1, // tells how much the target element should intersects/falls with/under the margin where 1 = complete
      },
    );

    if (observerRef.current) {
      observer.observe(observerRef.current);
    }

    return () => {
      if (observerRef.current) {
        observer.unobserve(observerRef.current);
      }
    };
  }, [fetchNextPage, hasNextPage]);

  const allMovies = movies?.pages?.flatMap((page) => page.movies ?? []) ?? [];

  if (allMovies == null || allMovies.length === 0) {
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
        <>
          <Movies movies={allMovies} message={message} />
          <h2>
            {isFetchingNextPage && <Spinner />}
          </h2>
          <div ref={observerRef} style={{ height: "20px" }} />
        </>
      )}
    </>
  );
};

export default Home;
