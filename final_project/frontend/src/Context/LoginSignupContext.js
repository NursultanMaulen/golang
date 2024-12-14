import { createContext, useContext, useReducer, useEffect } from "react";
import { toast } from "react-toastify";

const loginSignupContext = createContext();
export const useLoginSignupContext = () => useContext(loginSignupContext);

function LoginSignupContext({ children }) {
  const [state, dispatch] = useReducer(reducerFn, {
    user: null,
    token: null,
    isAuthenticated: false,
  });

  function reducerFn(state, action) {
    switch (action.type) {
      case "LOGIN":
        return {
          ...state,
          user: action.payload.user,
          token: action.payload.token,
          isAuthenticated: true,
        };
      case "LOGOUT":
        return {
          ...state,
          user: null,
          token: null,
          isAuthenticated: false,
        };
      case "UPDATE_USER":
        return {
          ...state,
          user: { ...state.user, ...action.payload },
        };
      default:
        return state;
    }
  }

  useEffect(() => {
    const storedToken = localStorage.getItem("token");
    const storedUser = localStorage.getItem("userData");

    if (storedToken && storedUser) {
      try {
        const parsedUser = JSON.parse(storedUser);
        dispatch({
          type: "LOGIN",
          payload: {
            user: parsedUser,
            token: storedToken,
          },
        });
        console.log("Restored user from localStorage:", parsedUser);
      } catch (error) {
        console.error("Error parsing user data from localStorage:", error);
      }
    } else {
      console.warn("No token or userData in localStorage");
    }
  }, []);

  const loginUser = async (email, password) => {
    try {
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

      dispatch({
        type: "LOGIN",
        payload: {
          user: data.user,
          token: data.token,
        },
      });

      return { success: true }; // Возвращаем успех
    } catch (error) {
      console.error("Login error:", error);
      return { success: false, message: error.message || "Failed to log in" }; // Возвращаем ошибку
    }
  };

  const logoutUser = () => {
    localStorage.removeItem("token");
    localStorage.removeItem("userData");
    dispatch({ type: "LOGOUT" });
    toast.info("Logged out successfully.");
  };

  const updateUser = async (updatedData) => {
    try {
      const response = await fetch(
        `http://localhost:8080/api/users/${state.user.id}`,
        {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${state.token}`,
          },
          body: JSON.stringify(updatedData),
        }
      );

      if (!response.ok) {
        throw new Error("Failed to update user data");
      }

      const updatedUser = await response.json();

      // Обновляем данные пользователя как в localStorage, так и в состоянии
      localStorage.setItem("userData", JSON.stringify(updatedUser));
      dispatch({ type: "UPDATE_USER", payload: updatedUser });

      toast.success("User data updated successfully!");
    } catch (error) {
      console.error("Error updating user data:", error);
      toast.error("Failed to update user data.");
    }
  };

  const { user, token, isAuthenticated } = state;

  return (
    <loginSignupContext.Provider
      value={{
        state,
        dispatch,
        user,
        token,
        isAuthenticated,
        loginUser,
        logoutUser,
        updateUser,
      }}
    >
      {children}
    </loginSignupContext.Provider>
  );
}

export default LoginSignupContext;
