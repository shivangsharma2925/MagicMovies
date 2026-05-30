import axiosConfig from "../../api/axiosConfig"

export const fetchMovies = async ({ search = "", pageParam = "" }) => {
  const trimmedSearch = search.trim();

  const url = trimmedSearch
    ? `/movies?search=${encodeURIComponent(trimmedSearch)}&cursor=${pageParam}&limit=10`
    : `/movies?cursor=${pageParam}&limit=10`;

  const response = await axiosConfig.get(url);

  return response.data;
};

export const fetchMovieSuggestions = async (search) => {
  if (!search.trim()) return [];

  const response = await axiosConfig.get(
    `/movies/suggestions?q=${encodeURIComponent(search)}`
  );
  return response.data;
};