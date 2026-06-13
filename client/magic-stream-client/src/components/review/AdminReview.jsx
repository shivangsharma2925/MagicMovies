import { Form, Button } from "react-bootstrap";
import { useRef, useEffect, useState } from "react";
import { useParams } from "react-router-dom";
//import axiosPrivate from '../../api/axiosPrivateConfig';
import useAxiosPrivate from "../../hook/UseAxiosPrivate";
import UseAuth from "../../hook/UseAuth";
import Movie from "../movie/Movie";
import Spinner from "../../utils/Spinner";

const Review = () => {
  const [movie, setMovie] = useState({});
  const [loading, setLoading] = useState(false);
  const revText = useRef();
  const { imdb_id } = useParams();
  const { auth, setAuth } = UseAuth();
  const axiosPrivate = useAxiosPrivate();

  useEffect(() => {
    const fetchMovie = async () => {
      setLoading(true);
      try {
        const response = await axiosPrivate.get(`/movie/${imdb_id}`);
        setMovie(response.data);
      } catch (error) {
        console.error("Error fetching movie:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchMovie();
  }, [axiosPrivate, imdb_id]);

  const handleSubmit = async (e) => {
    e.preventDefault();

    setLoading(true);
    try {
      const response = await axiosPrivate.post(`/updatereview/${imdb_id}`, {
        admin_review: revText.current.value,
      });

      setMovie(() => ({
        ...movie,
        admin_review: response.data?.admin_review || movie.admin_review,
        ranking: {
          ranking_name:
            response.data?.ranking_name || movie.ranking?.ranking_name,
        },
      }));
    } catch (err) {
      console.error(err);
      if (err.response && err.response.status === 401) {
        console.error("Unauthorized access - redirecting to login");
        setAuth(null);
      } else {
        console.error("Error updating review:", err);
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
      <div className="container py-5">
        <h2 className="review-page-title">Admin Review</h2>
        <div className="row justify-content-center g-4">
          {/* left - movie card */}
          <div className="col-12 col-md-6">
            {/* <div className="review-panel h-100">
                <Movie movie={movie} />
              </div> */}
            <div className="review-movie-poster">
              <img src={movie.poster_path} alt={movie.title} />
              {movie.ranking?.ranking_name && (
                <span className="movie-rank-badge">
                  {movie.ranking.ranking_name}
                </span>
              )}
            </div>

            <div className="review-movie-meta">
              <h5 className="review-movie-title">{movie.title}</h5>
              <span className="review-movie-imdb">{movie.imdb_id}</span>
            </div>
          </div>

          {/* right - review card */}
          <div className="col-12 col-md-6">
            <div className="review-panel">
              {auth?.role === "ADMIN" ? (
                <Form onSubmit={handleSubmit}>
                  <p className="review-panel-title">Write Review</p>
                  <Form.Group className="mb-3" controlId="adminReviewTextarea">
                    <Form.Label className="review-form-label">
                      Admin Review
                    </Form.Label>
                    <Form.Control
                      ref={revText}
                      required
                      as="textarea"
                      rows={7}
                      disabled={loading}
                      defaultValue={movie?.admin_review}
                      placeholder="Write your review here..."
                      className="review-textarea"
                      style={{ resize: "vertical" }}
                    />
                  </Form.Group>
                  <div className="d-flex justify-content-end">
                    <Button className="btn-submit-review" type="submit" disabled={loading}>
                      {!loading ? "Submit Review" : "Submitting..."}
                    </Button>
                  </div>
                </Form>
              ) : (
                <div>
                  <p className="review-panel-title">Review</p>
                  <div className="review-display-box">
                    <p className="review-display-label">Admin Review</p>
                    <p className="review-display-text">{movie.admin_review}</p>
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </>
  );
};

export default Review;
