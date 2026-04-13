import { Route, Routes, useNavigate } from "react-router-dom";
import "./App.css";
import Header from "./components/header/Header";
import Home from "./components/home/Home";
import Register from "./components/regsiter/Regsiter";
import Login from "./components/login/Login";
import Layout from "./components/Layout";
import RequiredAuth from "./components/RequiredAuth";
import Recommended from "./components/recommended/Recommended";
import Review from "./components/review/AdminReview";
import axiosConfig from "./api/axiosConfig";
import UseAuth from "./hook/UseAuth";
import StreamMovie from "./components/stream/StreamMovie";

function App() {

  const navigate = useNavigate();
  const {auth, setAuth} = UseAuth();

  const updateMovieReview = (imdb_id) => {
      navigate(`/review/${imdb_id}`);
  };

  const handleLogout = async () => {
      try {
          await axiosConfig.post("/logout",{user_id: auth.user_id});
          setAuth(null);
          localStorage.removeItem('user');
          // console.log('User logged out');

      } catch (error) {
          console.error('Error logging out:', error);
      } 
  };

  return (
    <>
      <Header handleLogout = {handleLogout}/>
      <Routes>
        {/* Layout wrapper */}
        <Route path="/" element={<Layout />}>
          {/* Public routes */}
          <Route index element={<Home updateMovieReview={updateMovieReview}/>} />
          <Route path="register" element={<Register />} />
          <Route path="login" element={<Login />} />

          {/* Protected routes */}
          <Route element={<RequiredAuth />}>
            <Route path="recommended" element={<Recommended />} />
            <Route path="review/:imdb_id" element={<Review />} />
            <Route path="stream/:yt_id" element={<StreamMovie />} />
          </Route>
        </Route>
      </Routes>
    </>
  );
}

export default App;
