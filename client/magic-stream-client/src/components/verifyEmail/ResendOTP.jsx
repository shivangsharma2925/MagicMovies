import { useState } from "react";
import axiosConfig from "../../api/axiosConfig";
// import { useNavigate } from "react-router-dom";
import Button from "react-bootstrap/Button";

const ResendOTP = ({userId}) => {
    const [loading, setLoading] = useState(false);
    const [message, setMessage] = useState("");
    const [error, setError] = useState("");
    // const naviagte = useNavigate();

    const handleResend = async(e) => {
        e.preventDefault();
        setLoading(true);
        try {
            const res = await axiosConfig.post("/resend-verification", {
                "user_id": userId
            })

            setMessage(res.data.message);

        } catch (err) {
            setError(err?.response?.data?.error || "Failed to resend OTP");
            console.log(err.response.data.error);
        }finally{
            setLoading(false);
            setTimeout(() => {
                setMessage("");
            }, 2000);
        }
    }

    return (
        <div style={{ marginTop: "20px" }}>
            <Button
                variant="secondary"
                onClick={(e)=>handleResend(e)}
                disabled={loading}
            >
                {loading ? "Sending..." : "Resend OTP"}
            </Button>
            {message && <p>{message}</p>}
            {error && <p>{error}</p>}
        </div>
    )
}

export default ResendOTP;