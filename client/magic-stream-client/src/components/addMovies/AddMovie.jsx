import { useState } from "react";
import useAxiosPrivate from "../../hook/UseAxiosPrivate";
import UseAuth from "../../hook/UseAuth";

const AddMovie = () => {
  const [imdbId, setImdbId] = useState("");
  const [message, setMessage] = useState("");
  const axiosPrivate = useAxiosPrivate();
  const { auth } = UseAuth();

  const handleSubmit = async (e) => {
    e.preventDefault();

    try {
      const resp = await axiosPrivate.post("/addmovie", { imdb_ids: imdbId });
      setMessage("Added to queue!");
      if(resp?.data?.error){
        setMessage(resp.data.error);
      }
      setImdbId("");
    } catch (error) {
      setMessage(error?.response?.data?.error || error.message || "Error adding movies");
      console.log(error);
    } finally{
        setMessage("");
    }
  };

  return (
    <div className="container mt-4">
      <h2>Add Movie</h2>
      {
        auth && auth.role === "ADMIN" ? (
          <>
            <form onSubmit={handleSubmit}>
              <input
                type="text"
                placeholder="Enter comma separated IMDB ID(s)"
                value={imdbId}
                onChange={(e) => setImdbId(e.target.value)}
                className="form-control mb-3"
              />

              <button className="btn btn-primary">Add to Queue</button>
            </form>
            {message && <p className="mt-3">{message}</p>}
          </>
        ) : (
          <h4>You don't have access to this functionality!!</h4>
        )
      }
    </div>
  );
};

export default AddMovie;
