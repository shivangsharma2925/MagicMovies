import { Outlet } from "react-router-dom";
import { Navigate, useLocation } from "react-router-dom";
import UseAuth from "../hook/UseAuth";
import Spinner from "../utils/Spinner";

const RequiredAuth = () => {
    const {auth, loading} = UseAuth();
    const location = useLocation();

    return loading ? ( <Spinner /> ) :
    auth ? (
        <Outlet />
    ) : (
        <Navigate to="/login" replace state={{ from: location }} />
    )
}

export default RequiredAuth;