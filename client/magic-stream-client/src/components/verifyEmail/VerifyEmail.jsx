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

            setTimeout(() => navigate("/login"), 2000);

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
        <Container className="d-flex align-items-center justify-content-center">
            <div className="">
                <h2 className="fw-bold">Verify Email</h2>

                <Form onSubmit={handleVerify}>
                    <Form.Group className="mb-3">
                        <Form.Label>Enter the OTP send to your Email Id</Form.Label>
                        <div style={{ display: "flex", gap: "10px", justifyContent: "center" }}>
                            {otp.map((digit, index) => (
                                <input
                                    key={index}
                                    ref={(el) => (inputRefs.current[index] = el)}
                                    type="text"
                                    maxLength={1}
                                    value={digit}
                                    onChange={(e) => handleChange(e.target.value, index)}
                                    onKeyDown={(e) => handleKeyDown(e, index)}
                                    onPaste={handlePaste}
                                    style={{
                                        width: "45px",
                                        height: "45px",
                                        textAlign: "center",
                                        fontSize: "20px",
                                        border: "1px solid #ccc",
                                        borderRadius: "8px",
                                        outline: "none",
                                    }}
                                />
                            ))}
                        </div>
                    </Form.Group>
                    <Form.Group className="mb-3">
                        <Button
                            variant="primary"
                            onClick={(e) => handleVerify(otp, e)}
                            disabled={loading}
                            className="mt-3"
                        >
                            Verify
                        </Button>
                    </Form.Group>
                </Form>

                {message && (
                    <p>{message}</p>
                )}

                {error && (
                    <p>{error}</p>
                )}

                <ResendOTP userId={userId} />
            </div>
       </Container>
    );
}