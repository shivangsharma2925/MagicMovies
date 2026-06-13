import Movie from "../movie/Movie";

const Movies = ({ movies, message }) => {
  if (movies == null || !movies.length) {
    return <h2 className="text-center text-white mt-4">{message}</h2>;
  }

  return (
    <div className="container-fluid px-4 mt-4">
      <div className="row g-4">
        {movies.map((movie) => (
          <Movie key={movie._id} movie={movie} />
        ))}
      </div>
    </div>
  );
};

export default Movies;
