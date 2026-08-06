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
import { Helmet } from "react-helmet-async";
import logo from "../../assets/logo.png";
import "../../styling/Home.css";

const Home = () => {
  const [hasNewMovies, setHasNewMovies] = useState(false);

  const queryClient = useQueryClient();

  const socket = UseWebSocket();

  let message = "";

  const title = "Magic Stream";
  const description = `Watch movies trailer and explore movie details, add movies to watchlists and see personalized recommended movies.`;
  const pageUrl = `https://www.magicmovies.site/`;
  // const pageUrl = `http://localhost:3000/`;
  const imageUrl = logo;

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
  const { data: heroMovies, isLoading: heroLoading } = useQuery({
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

  //for JSON-LD (SEO)
  const HomeStructuredData = {
    "@context": "https://schema.org",
    "@graph": [
      {
        "@type": "WebSite",
        "@id": "https://www.magicmovies.site/#website",
        url: "https://www.magicmovies.site/",
        name: "Magic Stream",
        alternateName: "Magic Movies",
        description:
          "Discover movie trailers, explore new releases, browse by genre, and find your next movie.",
        inLanguage: "en",
        publisher: {
          "@id": "https://www.magicmovies.site/#organization",
        },
      },
      {
        "@type": "Organization",
        "@id": "https://www.magicmovies.site/#organization",
        name: "Magic Stream",
        url: "https://www.magicmovies.site/",
        logo: {
          "@type": "ImageObject",
          "@id": "https://www.magicmovies.site/#logo",
          url: "https://www.magicmovies.site/logo.png",
          contentUrl: "https://www.magicmovies.site/logo.png",
          width: 512,
          height: 512,
        },
      },
      {
        "@type": "ItemList",
        "@id": "https://www.magicmovies.site/#featured-movies",
        name: "Featured Movies",
        itemListOrder: "https://schema.org/ItemListOrderAscending",
        numberOfItems: heroMovies?.excellent.length ?? 0,
        itemListElement: (heroMovies?.excellent?.slice(0, 10) ?? []).map(
          (movie, index) => ({
            "@type": "ListItem",
            position: index + 1,
            name: movie.title,
            url: `https://www.magicmovies.site/movie/detail/${movie.imdb_id}`,
          }),
        ),
      },
    ],
  };

  return (
    <>
      {!heroLoading && (
        <Helmet>
          <title>{title}</title>

          {/* Search engine description */}
          <meta name="description" content={description} />

          {/* Search engine crawling */}
          <meta name="robots" content="index, follow" />

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
          <script
            type="application/ld+json"
            dangerouslySetInnerHTML={{
              __html: JSON.stringify(HomeStructuredData),
            }}
          />
        </Helmet>
      )}

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
