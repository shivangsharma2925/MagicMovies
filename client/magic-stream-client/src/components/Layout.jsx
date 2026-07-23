import { Outlet } from "react-router-dom";
import ScrollToTopBtn from "../utils/ScrollToTopBtn";

const Layout = () => {
    return (
        <main>
            <Outlet />
            <ScrollToTopBtn />
        </main>
    )
}

export default Layout;