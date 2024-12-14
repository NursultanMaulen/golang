import { React, toast } from "../Utils/CustomUtils";
import { useLoginSignupContext } from "../Context/LoginSignupContext"; 

export const logoutHandler = (dispatch) => {
  localStorage.clear();
  dispatch({ type: "LOGOUT" });
  toast.success("Logout success!");
};

export const signUpHandler = async (userData) => {
  try {
    const userWithRole = { ...userData, role: "user" };

    const response = await fetch("http://localhost:8080/api/users", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(userWithRole),
    });

    if (!response.ok) {
      throw new Error("Failed to sign up");
    }

    toast.success("Signup success!");
  } catch (error) {
    console.error("Error during signup:", error);
    toast.error("Signup failed. Please try again.");
  }
};

export const loginHandler = async (email, password, dispatch) => {
  try {
    dispatch({ type: "LOADING", payload: true });

    const response = await fetch(
      "http://localhost:8080/api/users/authenticate",
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ email, password }),
      }
    );

    if (!response.ok) {
      throw new Error("Invalid email or password");
    }

    const data = await response.json();

    localStorage.setItem("token", data.token);
    localStorage.setItem("userData", JSON.stringify(data.user));

    dispatch({ type: "LOGINDATA", payload: data.user });

    toast.success(`Welcome ${data.user.username}!`);
  } catch (error) {
    console.error("Error during login:", error);
    toast.error("Login failed. Please check your credentials.");
  } finally {
    dispatch({ type: "LOADING", payload: false });
  }
};
