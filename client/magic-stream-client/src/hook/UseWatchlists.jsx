import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import useAxiosPrivate from "../hook/UseAxiosPrivate";
import UseAuth from "./UseAuth";

export const useWatchlists = () => {
  const axiosPrivate = useAxiosPrivate();
  const { auth } = UseAuth();
  return useQuery({
    queryKey: ["watchlists"],
    queryFn: () => axiosPrivate.get("/watchlists").then((r) => r.data),
    enabled: auth !== null,
    staleTime: 5 * 60 * 1000,
  });
};

export const useCreateWatchlist = () => {
  const qc = useQueryClient();
  const axiosPrivate = useAxiosPrivate();
  return useMutation({
    mutationFn: (name) => axiosPrivate.post("/watchlists", { name }),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["watchlists"] }),
  });
};

export const useDeleteWatchlist = () => {
  const axiosPrivate = useAxiosPrivate();
  const qc = useQueryClient();
  return useMutation({
    mutationFn: (id) => axiosPrivate.delete(`/watchlists/${id}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["watchlists"] }),
  });
};

export const useRenameWatchlist = () => {
  const qc = useQueryClient();
  const axiosPrivate = useAxiosPrivate();
  return useMutation({
    mutationFn: ({ id, name }) =>
      axiosPrivate.post(`/watchlists/${id}`, { name }),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["watchlists"] }),
  });
};

export const useAddToWatchlist = () => {
  const axiosPrivate = useAxiosPrivate();
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ listId, movie }) =>
      axiosPrivate.post(`/watchlists/${listId}/movies`, movie),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["watchlists"] }),
  });
};

export const useRemoveFromWatchlist = () => {
  const axiosPrivate = useAxiosPrivate();
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ listId, imdbId }) =>
      axiosPrivate.delete(`/watchlists/${listId}/movies/${imdbId}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["watchlists"] }),
  });
};

export const useToggleWatched = () => {
  const axiosPrivate = useAxiosPrivate();
  const qc = useQueryClient();
  return useMutation({
    mutationFn: ({ listId, imdbId }) =>
      axiosPrivate.post(`/watchlists/${listId}/movies/${imdbId}`),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["watchlists"] }),
  });
};
