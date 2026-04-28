import { Route, Routes } from "react-router-dom";
import "./App.css";
import Header from "./components/header/Header";
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
          <Route path="register" element={<Register />} />
          <Route path="login" element={<Login />} />

          {/* Protected routes */}
          <Route element={<RequiredAuth />}>
            <Route path="recommended" element={<Recommended />} />
            <Route path="review/:imdb_id" element={<Review />} />
            <Route path="stream/:yt_id" element={<StreamMovie />} />
            <Route path="admin/add-movie" element={<AddMovie />} />
          </Route>
        </Route>
      </Routes>
    </>
  );
}

export default App;
