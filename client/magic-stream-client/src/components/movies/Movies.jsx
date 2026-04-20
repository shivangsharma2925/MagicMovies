import Movie from "../movie/Movie";

const Movies = ({ movies, message }) => {
  if (!movies.length) {
    return <h2 className="text-center mt-4">{message}</h2>;
  }

  return (
    <div className="container mt-4">
      <div className="row">
        {movies.map((movie) => (
          <Movie key={movie._id} movie={movie} />
        ))}
      </div>
    </div>
  );
};

export default Movies;
