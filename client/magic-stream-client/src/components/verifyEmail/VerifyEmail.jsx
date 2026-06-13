import { useRef, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import axiosConfig from "../../api/axiosConfig";
import ResendOTP from "./ResendOTP";
import Form from "react-bootstrap/Form";
import Button from "react-bootstrap/Button";
import Container from "react-bootstrap/Container";

export default function VerifyEmail() {
  const { userId } = useParams();
  const navigate = useNavigate();

  const [otp, setOtp] = useState(new Array(6).fill(""));
  const inputRefs = useRef([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [message, setMessage] = useState("");

  const handleVerify = async (otpArray, event) => {
    if (event) event.preventDefault();

    const otpValue = otpArray.join("");

    if (otpValue.length !== 6) {
      setError("Please enter full 6-digit OTP");
      return;
    }

    setLoading(true);
    setError("");
    setMessage("");

    try {
      const res = await axiosConfig.post("/verify-email", {
        user_id: userId,
        otp: otpValue,
      });

      setMessage(res.data.message);

      setTimeout(() => navigate("/login"), 1000);
    } catch (err) {
      setError(err.response?.data?.error || "Verification failed");
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (value, index) => {
    if (!/^[0-9]?$/.test(value)) return;

    const newOtp = [...otp];
    newOtp[index] = value;
    setOtp(newOtp);

    // move to next box
    if (value && index < 5) {
      inputRefs.current[index + 1].focus();
    }

    // auto submit if complete
    if (newOtp.every((digit) => digit !== "")) {
      handleVerify(newOtp);
    }
  };

  const handleKeyDown = (e, index) => {
    if (e.key === "Backspace" && !otp[index] && index > 0) {
      inputRefs.current[index - 1].focus();
    }
  };

  const handlePaste = (e) => {
    const pasted = e.clipboardData.getData("text").slice(0, 6);

    if (!/^[0-9]+$/.test(pasted)) return;

    const newOtp = pasted.split("").concat(Array(6).fill("")).slice(0, 6);
    setOtp(newOtp);

    const lastIndex = Math.min(pasted.length, 5);
    inputRefs.current[lastIndex]?.focus();
  };

  return (
    <Container className="login-container d-flex align-items-center justify-content-center">
      <div className="login-card">
        <div className="text-center mb-4">
          <div className="login-logo fs-1" aria-label="Email icon">📧</div>
          <h2 className="login-title">Verify Email</h2>
          <p className="login-subtitle">Enter the OTP sent to your email</p>
        </div>

        {error && <div className="login-alert login-alert-error">{error}</div>}
        {message && (
          <div className="login-alert login-alert-success">{message}</div>
        )}

        <Form onSubmit={handleVerify}>
          <Form.Group className="mb-4">
            <div className="otp-inputs">
              {otp.map((digit, index) => (
                <input
                  key={index}
                  ref={(el) => (inputRefs.current[index] = el)}
                  type="text"
                  inputMode="numeric"
                  maxLength={1}
                  value={digit}
                  onChange={(e) => handleChange(e.target.value, index)}
                  onKeyDown={(e) => handleKeyDown(e, index)}
                  onPaste={handlePaste}
                  className={`otp-box ${digit ? "otp-box-filled" : ""}`}
                />
              ))}
            </div>
          </Form.Group>

          <Button
            type="submit"
            className="login-btn w-100"
            onClick={(e) => handleVerify(otp, e)}
            disabled={loading}
          >
            {loading ? (
              <>
                <span
                  className="spinner-border spinner-border-sm me-2"
                  role="status"
                  aria-hidden="true"
                />
                Verifying...
              </>
            ) : (
              "Verify Email"
            )}
          </Button>
        </Form>

        <div className="login-divider">
          <hr />
          <span>or</span>
          <hr />
        </div>

        <ResendOTP userId={userId} />
      </div>
    </Container>
  );
}
