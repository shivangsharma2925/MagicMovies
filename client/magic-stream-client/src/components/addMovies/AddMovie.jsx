import React, { useEffect, useState } from "react";
import UseAuth from "../../hook/UseAuth";
import useAxiosPrivate from "../../hook/UseAxiosPrivate";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import UseWebSocket from "../../hook/UseWebSocket";
import AskAiBtn from "../../utils/AskAiModal";
const JOBS_PER_PAGE = 8;

const JobRow = React.memo(({ job, retryMutation }) => {
  const badgeClass =
    {
      done: "queue-badge-done",
      failed: "queue-badge-failed",
      processing: "queue-badge-processing",
      queued: "queue-badge-queued",
    }[job.status] ?? "queue-badge-queued";

  return (
    <tr className="queue-row">
      <td>{job.title ? job.title : job.imdb_id}</td>

      <td>
        <span className={`queue-badge ${badgeClass}`}>{job.status}</span>
      </td>

      <td>{job.attempts}</td>

      <td className={job.error ? "queue-error-text" : ""}>
        {job.error || "—"}
      </td>

      <td>
        {job.status === "failed" ? (
          <button
            className="btn-retry-job"
            onClick={() => retryMutation.mutate(job.imdb_id)}
            disabled={retryMutation.isPending}
          >
            Retry
          </button>
        ) : (
          <span className={`queue-badge ${badgeClass}`}>{job.status}</span>
        )}
      </td>
    </tr>
  );
});

const AddMovie = () => {
  const [imdbId, setImdbId] = useState("");
  const [message, setMessage] = useState("");
  const [currentPage, setCurrentPage] = useState(1);

  // const socket = useRef(null);

  const axiosPrivate = useAxiosPrivate();

  const queryClient = useQueryClient();

  const { auth } = UseAuth(); //here provider has {{ auth }}

  const socket = UseWebSocket(); //here provider has { socket }

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
      setMessage("Movies added to queue!");
      return axiosPrivate.post("/addmovie", {
        imdb_ids: imdbIds,
      });
    },

    onSuccess: () => {
      setTimeout(() => {
        setMessage("");
      }, 2000);

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
        error?.response?.data?.error || error.message || "Error adding movies",
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
    const Socket = socket.current;

    const handler = (event) => {
      const updatedJob = JSON.parse(event.data);

      if (updatedJob.type === "job_update") {
        // Update React Query cache directly
        queryClient.setQueryData(["jobs"], (oldJobs = []) => {
          return oldJobs.map((job) =>
            job.imdb_id === updatedJob.imdb_id
              ? { ...job, ...updatedJob }
              : job,
          );
        });
      }
    };

    Socket.addEventListener("message", handler);

    return () => {
      Socket.removeEventListener("message", handler);
    };
  }, [queryClient, socket]);

  const handleSubmit = (e) => {
    e.preventDefault();

    setMessage("");

    if (!imdbId.trim()) {
      setMessage("Please enter IMDB ID");
      return;
    }

    // if(imdbId.split(',').length > 5){
    //   setMessage("Please enter 5 ids at max")
    //   return;
    // }

    addMovieMutation.mutate(imdbId);
  };

  if (isLoading) {
    return (
      <div className="text-white container text-center mt-4">
        <h3>Loading jobs...</h3>
      </div>
    );
  }

  if (error) {
    return (
      <div className="text-white text-center container mt-4">
        <h3>Error loading jobs</h3>
      </div>
    );
  }

  if (!auth || auth.role !== "ADMIN") {
    return (
      <div className="container text-center text-white mt-4">
        <h4>You don't have access to this functionality!</h4>
      </div>
    );
  }

  const totalPages = Math.ceil(jobs.length / JOBS_PER_PAGE);
  const paginatedJobs = jobs.slice(
    (currentPage - 1) * JOBS_PER_PAGE,
    currentPage * JOBS_PER_PAGE,
  );

  return (
    <div className="container mt-4">
      <h2 className="queue-page-title">Add Movie</h2>

      <form onSubmit={handleSubmit} className="add-movie-form">
        <input
          type="text"
          placeholder="Enter comma separated IMDB ID(s) e.g. tt0816692, tt1375666"
          value={imdbId}
          onChange={(e) => setImdbId(e.target.value)}
          className="add-movie-input"
        />

        <div className="d-flex justify-content-space-between gap-2">
          <AskAiBtn onUseIds={(ids) => setImdbId(ids)} />

          <button
            className="btn-add-queue"
            disabled={addMovieMutation.isPending}
          >
            {addMovieMutation.isPending ? "Adding..." : "Add to Queue"}
          </button>
        </div>
      </form>
      {message && <p className="queue-message">{message}</p>}

      <div className="mt-4">
        <h2 className="queue-page-title">Queue Dashboard</h2>

        {jobs.length === 0 ? (
          <div className="queue-empty">No movies processed yet!</div>
        ) : (
          <>
            <div className="queue-table-wrap">
              <table className="queue-table">
                <thead>
                  <tr>
                    <th>Title</th>
                    <th>Status</th>
                    <th>Attempts</th>
                    <th>Error</th>
                    <th>Action</th>
                  </tr>
                </thead>

                <tbody>
                  {paginatedJobs.map((job) => (
                    <JobRow
                      key={job._id}
                      job={job}
                      retryMutation={retryMutation}
                    />
                  ))}
                </tbody>
              </table>
            </div>

            {/* client side Pagination */}
            <div className="queue-pagination">
              <span className="queue-page-info">
                Showing {(currentPage - 1) * JOBS_PER_PAGE + 1}–
                {Math.min(currentPage * JOBS_PER_PAGE, jobs.length)} of{" "}
                {jobs.length} jobs
              </span>
              <div className="queue-page-btns">
                <button
                  className="queue-page-btn"
                  onClick={() => setCurrentPage((p) => p - 1)}
                  disabled={currentPage === 1}
                >
                  ‹
                </button>
                {Array.from({ length: totalPages }, (_, i) => i + 1).map(
                  (p) => (
                    <button
                      key={p}
                      className={`queue-page-btn ${p === currentPage ? "active" : ""}`}
                      onClick={() => setCurrentPage(p)}
                    >
                      {p}
                    </button>
                  ),
                )}
                <button
                  className="queue-page-btn"
                  onClick={() => setCurrentPage((p) => p + 1)}
                  disabled={currentPage === totalPages}
                >
                  ›
                </button>
              </div>
            </div>
          </>
        )}
      </div>
    </div>
  );
};

export default AddMovie;
