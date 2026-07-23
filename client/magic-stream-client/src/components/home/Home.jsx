import React, { useEffect, useState } from "react";
import Movies from "../movies/Movies";
import { Button } from "react-bootstrap";
import Spinner from "../../utils/Spinner";
import {
  useInfiniteQuery,
  useQuery,
  useQueryClient,
} from "@tanstack/react-query";
import UseWebSocket from "../../hook/UseWebSocket";
import { fetchMovies, getHeroMovies } from "./GetMovieServices";
import useInfiniteScroll from "../../hook/UseInfiniteScroll";
import SearchBar from "./SearchBar";
import { Carousel } from "../movieCarousel/Carousel";
import MovieCarousel from "../movieCarousel/MovieCarousel";
import HeroSection from "../heroSection/HeroSection";
import "../../styling/Home.css";

const Home = () => {
  const [hasNewMovies, setHasNewMovies] = useState(false);

  const queryClient = useQueryClient();

  const socket = UseWebSocket();

  let message = "";

  /*
    The data in the react query movies array is stored as an set of all the pages we have ordered so far from server,
    example,
    movies = {
        "pages": [
            {
                "movies": [first 10 movies], "nextCursor": "last movie id from this batch", "hasmore": true
            },
            {
                "movies": [second 10 movies], "nextCursor": "last movie id from this batch", "hasmore": true
            },
        ],
        pageParams: [cursor1, cursor2]
    } 

    Complete Process for page and cursor based,
    0. Initially we have pageparam as 1 or "" to show data.
    1. Intersecton oberser keep the track of the marked element (in this case the div at the bottom), as soon as that element intersects or visible on the screen, the event fires up.
    2. inside the event we have "hasNextPage", which uses the "fetchNextPage", a messesnger that tells react query to fetch next page or cursor. But react query doesn't know the value of next page or cursor so it calls the function "getNextPageParam" by giving the last page object, the function returns the next page value by looking at the hasMore flag we have send from backend.
    3. The next page or cursor value now gets stored inside the pageParams array.
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
    queryKey: ["movies"],

    queryFn: ({ pageParam = "" }) => fetchMovies({ pageParam }),

    initialPageParam: "",

    getNextPageParam: (lastPage) => {
      return lastPage.hasMore ? lastPage.nextCursor : undefined;
    },

    placeholderData: (prev) => prev,

    staleTime: 1000 * 60 * 5,
  });

  /* fetching hero movies */
  const {
    data: heroMovies,
    isLoading: heroLoading,
  } = useQuery({
    queryKey: ["heroMovies"],

    queryFn: () => getHeroMovies(),

    staleTime: 1000 * 60 * 5,
  });

  const observerRef = useInfiniteScroll(fetchNextPage, hasNextPage);

  useEffect(() => {
    const Socket = socket.current;

    const handler = (event) => {
      const data = JSON.parse(event.data);

      if (data.type === "new_movie") {
        setHasNewMovies(true);
      }
    };

    Socket.addEventListener("message", handler);

    return () => {
      Socket.removeEventListener("message", handler);
    };
  }, [socket]);

  // flatten movies data
  const allMovies = movies?.pages?.flatMap((page) => page.movies ?? []) ?? [];

  if (allMovies == null || allMovies.length === 0) {
    message = `There are currently no movies available`;
  }

  if (error) {
    console.log(error);
    message = "Error fetching movies";
  }

  // get carousels movies
  const carousels = React.useMemo(() => {
    return Carousel({ heroMovies });
  }, [heroMovies]);

  return (
    <>
      {hasNewMovies && (
        <div className="new-movies-banner">
          <button
            className="new-movies-btn"
            onClick={() => {
              queryClient.invalidateQueries({
                queryKey: ["movies"],
              });

              setHasNewMovies(false);
            }}
          >
            <span className="pulse-dot" />
            New movies available
          </button>
        </div>
      )}

      <div>
        {heroLoading ? (
          <Spinner />
        ) : (
          <>
            <HeroSection movies={heroMovies} />

            <div className="carousels-container">
              {carousels.map((carousel) => (
                <MovieCarousel
                  key={carousel.title}
                  section={carousel.section}
                  title={carousel.title}
                  movies={carousel.movies}
                ></MovieCarousel>
              ))}
            </div>
          </>
        )}

        {isLoading ? (
          <Spinner />
        ) : (
          <Movies movies={allMovies} message={message} />
        )}
        {isFetchingNextPage && <Spinner />}

        <div ref={observerRef} style={{ height: "20px" }} />
      </div>
    </>
  );
};

export default Home;
