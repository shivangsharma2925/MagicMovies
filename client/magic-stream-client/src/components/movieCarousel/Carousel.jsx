import React from "react";

export const Carousel = ({ heroMovies = [] }) => {
  /* based on Ranking */
  const mustWatch = heroMovies?.excellent ?? [];
  const goodWatch = heroMovies?.good ?? [];
  const okayWatch = heroMovies?.okay ?? [];

  /* based on release date */
  const newRelease = heroMovies.newRelease;

  const nonGenreKeys = ["excellent", "good", "okay", "newRelease"];

  /* based on Genres */
  const genres = Object.entries(heroMovies).filter(([genreName, movies]) => {
    return !nonGenreKeys.includes(genreName) && movies.length > 3;
  });

  const carouselSections = [
    { title: "🏆 Excellent", section: "Excellent", movies: mustWatch },
    { title: "👍 Good Picks", section: "Good", movies: goodWatch },
    { title: "🎭 Okay Watches", section: "Okay", movies: okayWatch },
    { title: "✨ Recently Added", section: "newRelease", movies: newRelease },
    ...genres.map(([genreName, movies]) => ({
      title: `🎬 ${genreName}`,
      section: genreName,
      movies,
    })),
  ];

  return carouselSections;
};
