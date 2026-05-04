import React, { useEffect, useRef, useState } from "react";
import UseAuth from "../../hook/UseAuth";
import useAxiosPrivate from "../../hook/UseAxiosPrivate";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

const JobRow = React.memo(({ job, retryMutation }) => {
  return (
    <tr>
      <td>{job.imdb_id}</td>

      <td>
        <span
          className={`badge ${
            job.status === "done"
              ? "bg-success"
              : job.status === "failed"
                ? "bg-danger"
                : job.status === "processing"
                  ? "bg-warning text-dark"
                  : "bg-secondary"
          }`}
        >
          {job.status}
        </span>
      </td>

      <td>{job.attempts}</td>

      <td>{job.error || "-"}</td>

      <td>
        {job.status === "failed" ? (
          <button
            className="btn btn-sm btn-warning"
            onClick={() => retryMutation.mutate(job.imdb_id)}
            disabled={retryMutation.isPending}
          >
            Retry
          </button>
        ) : (
          job.status
        )}
      </td>
    </tr>
  );
});

const AddMovie = () => {
  const [imdbId, setImdbId] = useState("");
  const [message, setMessage] = useState("");

  const socket = useRef(null);

  const axiosPrivate = useAxiosPrivate();

  const queryClient = useQueryClient();

  const { auth } = UseAuth();

  /*
    =========================================
    FETCH JOBS (React Query)
    =========================================
  */

  const {
    data: jobs = [],
    isLoading,
    error,
  } = useQuery({
    queryKey: ["jobs"],

    queryFn: async () => {
      const res = await axiosPrivate.get("/jobs");
      return res.data;
    },

    staleTime: 1000 * 60,

    refetchOnWindowFocus: false,
  });

  /*
    =========================================
    ADD MOVIE MUTATION
    =========================================
  */

  const addMovieMutation = useMutation({
    mutationFn: async (imdbIds) => {
      return axiosPrivate.post("/addmovie", {
        imdb_ids: imdbIds,
      });
    },

    onSuccess: () => {
      setMessage("Movies added to queue!");

      setImdbId("");

      /*
        Marks query stale
        React Query auto refetches
      */
      queryClient.invalidateQueries({
        queryKey: ["jobs"],
      });
    },

    onError: (error) => {
      setMessage(
        error?.response?.data?.error ||
          error.message ||
          "Error adding movies",
      );
    },
  });

  /*
    =========================================
    RETRY JOB MUTATION
    =========================================
  */

  const retryMutation = useMutation({
    mutationFn: async (imdbId) => {
      return axiosPrivate.post(`/jobs/retry/${imdbId}`);
    },

    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: ["jobs"],
      });
    },
  });

  /*
    =========================================
    WEBSOCKET CONNECTION
    =========================================
  */

  useEffect(() => {
    socket.current = new WebSocket(
      import.meta.env.VITE_WS_URL || "ws://localhost:8080/api/v1/ws",
    );

    socket.current.onopen = () => {
      console.log("WebSocket Connected");
    };

    socket.current.onmessage = (event) => {
      const updatedJob = JSON.parse(event.data);

      /*
        Update React Query cache directly
      */
      queryClient.setQueryData(["jobs"], (oldJobs = []) => {
        return oldJobs.map((job) =>
          job.imdb_id === updatedJob.imdb_id
            ? { ...job, ...updatedJob }
            : job,
        );
      });
    };

    socket.current.onerror = (err) => {
      console.log("WebSocket Error:", err);
    };

    socket.current.onclose = () => {
      console.log("WebSocket Closed");
    };

    return () => {
      if (socket.current) {
        socket.current.close(1000, "Normal close");
        socket.current = null;
      }
    };
  }, [queryClient]);

  const handleSubmit = (e) => {
    e.preventDefault();

    setMessage("");

    if (!imdbId.trim()) {
      setMessage("Please enter IMDB ID");
      return;
    }

    // if(length(imdbId.split(',')) > 5){
    //   setMessage("Please enter 5 ids at max")
    //   return;
    // }

    addMovieMutation.mutate(imdbId);
  };

  if (isLoading) {
    return (
      <div className="container mt-4">
        <h3>Loading jobs...</h3>
      </div>
    );
  }

  if (error) {
    return (
      <div className="container mt-4">
        <h3>Error loading jobs</h3>
      </div>
    );
  }

  if (!auth || auth.role !== "ADMIN") {
    return (
      <div className="container mt-4">
        <h4>You don't have access to this functionality!</h4>
      </div>
    );
  }

  return (
    <div className="container mt-4">
      <h2>Add Movie</h2>

      <form onSubmit={handleSubmit}>
        <input
          type="text"
          placeholder="Enter comma separated IMDB ID(s)"
          value={imdbId}
          onChange={(e) => setImdbId(e.target.value)}
          className="form-control mb-3"
        />

        <button
          className="btn btn-primary"
          disabled={addMovieMutation.isPending}
        >
          {addMovieMutation.isPending ? "Adding..." : "Add to Queue"}
        </button>
      </form>

      {message && <p className="mt-3">{message}</p>}

      <div className="container mt-4">
        <h2>Queue Dashboard</h2>

        {jobs.length === 0 ? (
          <h4>No Movies processed yet!</h4>
        ) : (
          <table className="table table-bordered table-hover">
            <thead>
              <tr>
                <th>IMDB ID</th>
                <th>Status</th>
                <th>Attempts</th>
                <th>Error</th>
                <th>Action</th>
              </tr>
            </thead>

            <tbody>
              {jobs.map((job) => (
                <JobRow
                  key={job._id}
                  job={job}
                  retryMutation={retryMutation}
                />
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
};

export default AddMovie;