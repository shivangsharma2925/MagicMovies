import { useState } from "react";
import axiosConfig from "../../api/axiosConfig";
// import { useNavigate } from "react-router-dom";
import Button from "react-bootstrap/Button";

const ResendOTP = ({ userId }) => {
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");
  // const naviagte = useNavigate();

  const handleResend = async (e) => {
    e.preventDefault();
    setLoading(true);
    try {
      const res = await axiosConfig.post("/resend-verification", {
        user_id: userId,
      });

      setMessage(res.data.message);
    } catch (err) {
      setError(err?.response?.data?.error || "Failed to resend OTP");
      console.log(err.response.data.error);
    } finally {
      setLoading(false);
      setTimeout(() => {
        setMessage("");
      }, 2000);
    }
  };

  return (
    <div className="text-center">
      <Button
        variant="secondary"
        onClick={(e) => handleResend(e)}
        disabled={loading}
        className="mb-3"
      >
        {loading ? "Sending..." : "Resend OTP"}
      </Button>
      {error && <div className="login-alert login-alert-error">{error}</div>}
      {message && (
        <div className="login-alert login-alert-success">{message}</div>
      )}
    </div>
  );
};

export default ResendOTP;
