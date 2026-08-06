import { Route, Routes } from "react-router-dom";
import { lazy, Suspense } from "react";

import Header from "./components/navbar/NavBar";
import Home from "./components/home/Home";
import Layout from "./components/Layout";
import UseAuth from "./hook/UseAuth";
import RequiredAuth from "./components/RequiredAuth";
import Spinner from "./utils/Spinner";

// Public pages
const MovieDetails = lazy(() => import("./components/movie/MovieDetails"));
const Register = lazy(() => import("./components/regsiter/Regsiter"));
const Login = lazy(() => import("./components/login/Login"));
const VerifyEmail = lazy(() => import("./components/verifyEmail/VerifyEmail"));
const ResetPassword = lazy(
  () => import("./components/resetPassword/ResetPassword"),
);
const CategoryMovies = lazy(
  () => import("./components/categoryMovies/CategoryMovies"),
);

// Protected pages
const Recommended = lazy(() => import("./components/recommended/Recommended"));
const WatchlistPage = lazy(
  () => import("./components/Watchlist/WatchlistPage"),
);
const Review = lazy(() => import("./components/review/AdminReview"));
const StreamMovie = lazy(() => import("./components/stream/StreamMovie"));
const AddMovie = lazy(() => import("./components/addMovies/AddMovie"));
const DeleteAccount = lazy(
  () => import("./components/deleteAccount/DeleteAccount"),
);

import "./App.css";

function App() {
  const { loading } = UseAuth();

  if (loading)
    return (
      <div>
        <Spinner />
      </div>
    );
  return (
    <>
      <Header />
      <Suspense fallback={<Spinner />}>
        <Routes>
          <Route path="/" element={<Layout />}>
            {/* Public routes */}
            <Route index element={<Home />} />

            <Route path="movie/detail/:imdbId" element={<MovieDetails />} />
            <Route path="register" element={<Register />} />
            <Route path="login" element={<Login />} />
            <Route path="verify-email/:userId" element={<VerifyEmail />} />
            <Route path="reset-password" element={<ResetPassword />} />
            <Route path="category/:title" element={<CategoryMovies />} />

            {/* Protected routes */}
            <Route element={<RequiredAuth />}>
              <Route path="recommended" element={<Recommended />} />
              <Route path="watchlist" element={<WatchlistPage />} />
              <Route path="review/:imdb_id" element={<Review />} />
              <Route path="stream/:title/:imdb_id" element={<StreamMovie />} />
              <Route path="admin/add-movie" element={<AddMovie />} />
              <Route path="account/delete" element={<DeleteAccount />} />
            </Route>
          </Route>
        </Routes>
      </Suspense>
    </>
  );
}

export default App;
