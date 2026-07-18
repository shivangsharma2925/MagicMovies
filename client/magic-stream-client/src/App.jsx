import { Route, Routes } from "react-router-dom";
import "./App.css";
import Header from "./components/navbar/NavBar";
import Home from "./components/home/Home";
import Register from "./components/regsiter/Regsiter";
import Login from "./components/login/Login";
import Layout from "./components/Layout";
import RequiredAuth from "./components/RequiredAuth";
import Recommended from "./components/recommended/Recommended";
import Review from "./components/review/AdminReview";
import UseAuth from "./hook/UseAuth";
import StreamMovie from "./components/stream/StreamMovie";
import Spinner from "./utils/Spinner";
import AddMovie from "./components/addMovies/AddMovie";
import VerifyEmail from "./components/verifyEmail/VerifyEmail";
import ResetPassword from "./components/resetPassword/ResetPassword";
import MovieDetails from "./components/movie/MovieDetails";
import CategoryMovies from "./components/categoryMovies/CategoryMovies";
import WatchlistPage from "./components/Watchlist/WatchlistPage";
import DeleteAccount from "./components/deleteAccount/DeleteAccount";

function App() {
  const { loading } = UseAuth();
  
  if (loading) return <div><Spinner /></div>;
  return (
    <>
      <Header />
      <Routes>
        {/* Layout wrapper */}
        <Route path="/" element={<Layout />}>
          {/* Public routes */}
          <Route index element={<Home />} />
          <Route path="movie/detail/:imdbId" element={<MovieDetails />}></Route>
          <Route path="register" element={<Register />} />
          <Route path="login" element={<Login />} />
          <Route path="verify-email/:userId" element={<VerifyEmail />}/>
          <Route path="reset-password" element={<ResetPassword />}/>
          <Route path="category/:title" element={<CategoryMovies />}/>

          {/* Protected routes */}
          <Route element={<RequiredAuth />}>
            <Route path="recommended" element={<Recommended />} />
            <Route path="watchlist" element={<WatchlistPage />} />
            <Route path="review/:imdb_id" element={<Review />} />
            <Route path="stream/:yt_id" element={<StreamMovie />} />
            <Route path="admin/add-movie" element={<AddMovie />} />
            <Route path="account/delete" element={<DeleteAccount />} />
          </Route>
        </Route>
      </Routes>
    </>
  );
}

export default App;
