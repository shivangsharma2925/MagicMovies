import { createContext, useEffect, useState } from "react";

const AuthContext = createContext({});

export const AuthProvider = ({children}) => {

    const [auth, setAuth] = useState();
    const [loading, setLoading] = useState(true);

    useEffect(()=>{
        const checklogin = () => {
            const userData = localStorage.getItem("user");
            if (userData){
                setAuth(JSON.parse(userData)); //since localstorage stores string
            }else{
                setAuth(null);
            }
            setLoading(false);
        }
        checklogin();
    }, [])
    return (
        <AuthContext.Provider value={{auth, setAuth, loading}}>
            {children}
        </AuthContext.Provider>
    )
}

export default AuthContext;