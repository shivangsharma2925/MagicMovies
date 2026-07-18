import apiwocredential from "../../api/axiosConfigwocredential";

export const fetchMovies = async ({ search = "", pageParam = "" }) => {
  const trimmedSearch = search.trim();

  const url = trimmedSearch
    ? `/movies?search=${encodeURIComponent(trimmedSearch)}`
    : `/movies?cursor=${pageParam}&limit=10`;

  const response = await apiwocredential.get(url);

  return response.data;
};

export const fetchMovieSuggestions = async (search) => {
  if (!search.trim()) return [];

  const response = await apiwocredential.get(
    `/movies/suggestions?q=${encodeURIComponent(search)}`,
  );
  return response.data;
};

export const getHeroMovies = async () => {
  const res = await apiwocredential.get("/movies/heroMovies");
  return res.data;
};

export const fetchTitles = async ({ title = "" }) => {
  if (!title.trim()) return [];
  
  const res = await apiwocredential.get(`/category/${title}`);
  return res.data;
};

export const fetchMovie = async (imdbId) => {
  if (!imdbId.trim()) return ;

  const { data } = await apiwocredential.get(`/movie/${imdbId}`);
  return data;
};
