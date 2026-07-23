import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import useAxiosPrivate from "../hook/UseAxiosPrivate";
import UseAuth from "./UseAuth";

const WATCHLISTS_KEY = ["watchlists"];

function rollbackWatchlists(queryClient, context) {
  if (context?.previousWatchlists !== undefined) {
    queryClient.setQueryData(WATCHLISTS_KEY, context.previousWatchlists);
  }
}

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

// export const useAddToWatchlist = () => {
//   const axiosPrivate = useAxiosPrivate();
//   const qc = useQueryClient();
//   return useMutation({
//     mutationFn: ({ listId, movie }) =>
//       axiosPrivate.post(`/watchlists/${listId}/movies`, movie),
//     onSuccess: () => qc.invalidateQueries({ queryKey: ["watchlists"] }),
//   });
// };

export const useAddToWatchlist = () => {
  const axiosPrivate = useAxiosPrivate();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ listId, movie }) => {
      const response = await axiosPrivate.post(
        `/watchlists/${listId}/movies`,
        movie,
      );
      return response.data;
    },

    onMutate: async ({ listId, movie }) => {
      //cancel all the ongoing apis
      await queryClient.cancelQueries({
        queryKey: WATCHLISTS_KEY,
      });

      //store the current watchlist data for rollback if needed
      const previousWatchlists = queryClient.getQueryData(WATCHLISTS_KEY);

      const now = new Date().toISOString();

      queryClient.setQueryData(WATCHLISTS_KEY, (oldWatchlists = []) =>
        oldWatchlists.map((watchlist) => {
          //do nothing if watchlist id are different, meaning operation is performed on some other watchlist
          if (String(watchlist.id) !== String(listId)) {
            return watchlist;
          }

          const currentMovies = watchlist.movies ?? [];

          //if movie already exists, update nothing in query cache
          const alreadyExists = currentMovies.some(
            (existingMovie) => existingMovie.imdb_id === movie.imdb_id,
          );

          if (alreadyExists) {
            return watchlist;
          }

          return {
            ...watchlist,
            movies: [
              ...currentMovies,
              {
                ...movie,
                watched: movie.watched ?? false,
                added_at: movie.added_at ?? now,
                updated_at: now,
              },
            ],
          };
        }),
      );

      return {
        previousWatchlists,
      };
    },

    onError: (_error, _variables, context) => {
      rollbackWatchlists(queryClient, context);
    },

    onSettled: () => {
      // Do not return this promise.
      // The UI has already updated optimistically.
      void queryClient.invalidateQueries({
        queryKey: WATCHLISTS_KEY,
      });
    },
  });
};

// export const useRemoveFromWatchlist = () => {
//   const axiosPrivate = useAxiosPrivate();
//   const qc = useQueryClient();
//   return useMutation({
//     mutationFn: ({ listId, imdbId }) =>
//       axiosPrivate.delete(`/watchlists/${listId}/movies/${imdbId}`),
//     onSuccess: () => qc.invalidateQueries({ queryKey: ["watchlists"] }),
//   });
// };

export const useRemoveFromWatchlist = () => {
  const axiosPrivate = useAxiosPrivate();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ listId, imdbId }) => {
      const response = await axiosPrivate.delete(
        `/watchlists/${listId}/movies/${imdbId}`,
      );

      return response.data;
    },

    onMutate: async ({ listId, imdbId }) => {
      await queryClient.cancelQueries({
        queryKey: WATCHLISTS_KEY,
      });

      const previousWatchlists = queryClient.getQueryData(WATCHLISTS_KEY);

      queryClient.setQueryData(WATCHLISTS_KEY, (oldWatchlists = []) =>
        oldWatchlists.map((watchlist) =>
          String(watchlist.id) === String(listId)
            ? {
                ...watchlist,
                movies: (watchlist.movies ?? []).filter(
                  (movie) => movie.imdb_id !== imdbId,
                ),
              }
            : watchlist,
        ),
      );

      return {
        previousWatchlists,
      };
    },

    onError: (_error, _variables, context) => {
      rollbackWatchlists(queryClient, context);
    },

    onSettled: () => {
      void queryClient.invalidateQueries({
        queryKey: WATCHLISTS_KEY,
      });
    },
  });
};

// export const useToggleWatched = () => {
//   const axiosPrivate = useAxiosPrivate();
//   const qc = useQueryClient();
//   return useMutation({
//     mutationFn: ({ listId, imdbId }) =>
//       axiosPrivate.post(`/watchlists/${listId}/movies/${imdbId}`),
//     onSuccess: () => qc.invalidateQueries({ queryKey: ["watchlists"] }),
//   });
// };

export const useToggleWatched = () => {
  const axiosPrivate = useAxiosPrivate();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async ({ listId, imdbId, watched }) => {
      const response = await axiosPrivate.post(
        `/watchlists/${listId}/movies/${imdbId}/watched`,
        { watched },
      );

      return response.data;
    },

    onMutate: async ({ listId, imdbId }) => {
      await queryClient.cancelQueries({
        queryKey: WATCHLISTS_KEY,
      });

      const previousWatchlists = queryClient.getQueryData(WATCHLISTS_KEY);

      const now = new Date().toISOString();

      queryClient.setQueryData(WATCHLISTS_KEY, (oldWatchlists = []) =>
        oldWatchlists.map((watchlist) =>
          String(watchlist.id) === String(listId)
            ? {
                ...watchlist,
                movies: (watchlist.movies ?? []).map((movie) =>
                  movie.imdb_id === imdbId
                    ? {
                        ...movie,
                        watched: !movie.watched,
                        updated_at: now,
                      }
                    : movie,
                ),
              }
            : watchlist,
        ),
      );

      return {
        previousWatchlists,
      };
    },

    onError: (_error, _variables, context) => {
      rollbackWatchlists(queryClient, context);
    },

    onSettled: () => {
      void queryClient.invalidateQueries({
        queryKey: WATCHLISTS_KEY,
      });
    },
  });
};
